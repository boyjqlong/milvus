// Copyright (C) 2019-2020 Zilliz. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License
// is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
// or implied. See the License for the specific language governing permissions and limitations under the License

#include <limits>
#include <unordered_set>
#include <vector>

#include "Reduce.h"
#include "common/CGoHelper.h"
#include "common/Consts.h"
#include "common/Types.h"
#include "common/QueryResult.h"
#include "exceptions/EasyAssert.h"
#include "log/Log.h"
#include "pb/milvus.pb.h"
#include "query/Plan.h"
#include "segcore/ReduceStructure.h"
#include "segcore/SegmentInterface.h"
#include "segcore/reduce_c.h"
#include "segcore/Utils.h"

using SearchResult = milvus::SearchResult;

// void
// PrintSearchResult(char* buf, const milvus::SearchResult* result, int64_t seg_idx, int64_t from, int64_t to) {
//    const int64_t MAXLEN = 32;
//    snprintf(buf + strlen(buf), MAXLEN, "{ seg No.%ld ", seg_idx);
//    for (int64_t i = from; i < to; i++) {
//        snprintf(buf + strlen(buf), MAXLEN, "(%ld, %ld, %f), ", i, result->primary_keys_[i], result->distances_[i]);
//    }
//    snprintf(buf + strlen(buf), MAXLEN, "} ");
//}

void
ReduceResultData(std::vector<SearchResult*>& search_results, int64_t nq, int64_t topk) {
    AssertInfo(topk > 0, "topk must greater than 0");
    auto num_segments = search_results.size();
    AssertInfo(num_segments > 0, "num segment must greater than 0");
    for (int i = 0; i < num_segments; i++) {
        auto search_result = search_results[i];
        AssertInfo(search_result != nullptr, "search result must not equal to nullptr");
        AssertInfo(search_result->primary_keys_.size() == nq * topk, "incorrect search result primary key size");
        AssertInfo(search_result->distances_.size() == nq * topk, "incorrect search result distance size");
    }

    std::vector<std::vector<int64_t>> search_records(num_segments);
    std::unordered_set<milvus::PkType> pk_set;
    int64_t skip_dup_cnt = 0;

    // reduce search results
    for (int64_t qi = 0; qi < nq; qi++) {
        std::vector<SearchResultPair> result_pairs;
        int64_t base_offset = qi * topk;
        for (int i = 0; i < num_segments; i++) {
            auto search_result = search_results[i];
            auto primary_key = search_result->primary_keys_[base_offset];
            auto distance = search_result->distances_[base_offset];
            result_pairs.push_back(
                SearchResultPair(primary_key, distance, search_result, i, base_offset, base_offset + topk));
        }
        int64_t curr_offset = base_offset;

#if 0
        for (int i = 0; i < topk; ++i) {
            result_pairs[0].reset_distance();
            std::sort(result_pairs.begin(), result_pairs.end(), std::greater<>());
            auto& result_pair = result_pairs[0];
            auto index = result_pair.index_;
            result_pair.search_result_->result_offsets_.push_back(loc_offset++);
            search_records[index].push_back(result_pair.offset_++);
        }
#else
        pk_set.clear();
        while (curr_offset - base_offset < topk) {
            std::sort(result_pairs.begin(), result_pairs.end(), std::greater<>());
            auto& pilot = result_pairs[0];
            auto index = pilot.index_;
            auto curr_pk = pilot.primary_key_;
            // remove duplicates
            auto valid_pk = milvus::IsValidPrimaryKey(curr_pk);
            if (!valid_pk || pk_set.count(curr_pk) == 0) {
                pilot.search_result_->result_offsets_.push_back(curr_offset++);
                // when inserted data are dirty, it's possible that primary keys are duplicated,
                // in this case, "offset_" may be greater than "offset_rb_" (#10530)
                search_records[index].push_back(pilot.offset_ < pilot.offset_rb_ ? pilot.offset_ : INVALID_OFFSET);
                if (valid_pk) {
                    pk_set.insert(curr_pk);
                }
            } else {
                // skip entity with same primary key
                skip_dup_cnt++;
            }
            pilot.reset();
        }
#endif
    }
    LOG_SEGCORE_DEBUG_ << "skip duplicated search result, count = " << skip_dup_cnt;

    // after reduce, remove redundant values in primary_keys, distances and ids
    for (int i = 0; i < num_segments; i++) {
        auto search_result = search_results[i];
        if (search_result->result_offsets_.size() == 0) {
            continue;
        }

        std::vector<milvus::PkType> primary_keys;
        std::vector<float> distances;
        std::vector<int64_t> ids;
        for (int j = 0; j < search_records[i].size(); j++) {
            auto& offset = search_records[i][j];
            primary_keys.push_back(offset != INVALID_OFFSET ? search_result->primary_keys_[offset] : INVALID_PK);
            distances.push_back(offset != INVALID_OFFSET ? search_result->distances_[offset]
                                                         : std::numeric_limits<float>::max());
            ids.push_back(offset != INVALID_OFFSET ? search_result->seg_offsets_[offset] : INVALID_SEG_OFFSET);
        }

        search_result->primary_keys_ = primary_keys;
        search_result->distances_ = distances;
        search_result->seg_offsets_ = ids;
    }
}

std::unique_ptr<SearchResult>
ReorganizeSearchResults(std::vector<SearchResult*>& search_results,
                        milvus::query::Plan* plan,
                        int32_t nq,
                        int32_t topK) {
    AssertInfo(search_results.size() > 0, "empty search result when reorganize");
    auto final_search_result = std::make_unique<SearchResult>();
    final_search_result->primary_keys_.resize(nq * topK);
    final_search_result->distances_.resize(nq * topK);
    std::vector<std::pair<SearchResult*, int64_t>> result_offsets(nq *topK);

    auto num_segments = search_results.size();
    auto results_count = 0;

    for (int i = 0; i < num_segments; i++) {
        auto search_result = search_results[i];
        AssertInfo(search_result != nullptr, "null search result when reorganize");
        auto num_results = search_result->result_offsets_.size();
        if (num_results == 0) {
            continue;
        }
#pragma omp parallel for
        for (int j = 0; j < num_results; j++) {
            auto loc = search_result->result_offsets_[j];
            // set result pks
            final_search_result->primary_keys_[loc] = search_result->primary_keys_[j];
            // set result distances
            final_search_result->distances_[loc] = search_result->distances_[j];
            // set result offset to fill output fields data
            result_offsets[results_count + j] = std::make_pair(search_result, j);
        }
        results_count += num_results;
    }

    AssertInfo(results_count == nq * topK,
               "size of reduce result is less than nq * topK"
               ", result_count = " +
                   std::to_string(results_count) + ", nq * topK = " + std::to_string(nq * topK));

    for (auto field_id : plan->target_entries_) {
        auto& field_meta = plan->schema_[field_id];
        auto field_data = milvus::segcore::MergeDataArray(result_offsets, field_meta);
        final_search_result->output_fields_data_[field_id] = std::move(field_data);
    }

    return final_search_result;
}

std::vector<char>
GetSearchResultDataSlice(
                         SearchResult* final_search_result,
                         milvus::query::Plan* plan,
                         int64_t offset_begin,
                         int64_t offset_end,
                         int64_t nq,
                         int64_t topK) {
    auto search_result_data = std::make_unique<milvus::proto::schema::SearchResultData>();
    // set topK and nq
    search_result_data->set_top_k(topK);
    search_result_data->set_num_queries(nq);

    AssertInfo(offset_begin <= offset_end,
               "illegal offsets when GetSearchResultDataSlice"
               ", offset_begin = " +
                   std::to_string(offset_begin) + ", offset_end = " + std::to_string(offset_end));
    AssertInfo(offset_end <= topK * nq,
               "illegal offset_end when GetSearchResultDataSlice"
               ", offset_end = " +
                   std::to_string(offset_end) + ", nq = " + std::to_string(nq) + ", topK = " + std::to_string(topK));

    // set ids
    auto primary_field_id = plan->schema_.get_primary_field_id().value_or(milvus::FieldId(-1));
    AssertInfo(primary_field_id.get() != -1, "Primary key is -1");
    auto pk_type = plan->schema_[primary_field_id].get_data_type();
    auto proto_ids = std::make_unique<milvus::proto::schema::IDs>();
    switch (pk_type) {
        case milvus::DataType::INT64: {
            auto ids = std::make_unique<milvus::proto::schema::LongArray>();
            for (auto iter = offset_begin; iter < offset_end; iter++) {
                ids->mutable_data()->Add(std::get<int64_t>(final_search_result->primary_keys_[iter]));
            }

            proto_ids->set_allocated_int_id(ids.release());
            break;
        }
        case milvus::DataType::VarChar: {
            auto ids = std::make_unique<milvus::proto::schema::StringArray>();
            for (auto iter = offset_begin; iter < offset_end; iter++) {
                *(ids->mutable_data()->Add())= std::get<std::string>(final_search_result->primary_keys_[iter]);
            }

            proto_ids->set_allocated_str_id(ids.release());
            break;
        }
        default: {
            PanicInfo("unsupported primary key type");
        }
    }
    search_result_data->set_allocated_ids(proto_ids.release());

    // set scores
    *search_result_data->mutable_scores() = {final_search_result->distances_.begin() + offset_begin,
                                             final_search_result->distances_.begin() + offset_end};
    AssertInfo(search_result_data->scores_size() == offset_end - offset_begin,
               "wrong scores size"
               ", size = " +
                   std::to_string(search_result_data->scores_size()) +
                   ", expected size = " + std::to_string(offset_end - offset_begin));

    // set output fields
    for (int i = 0; i < plan->target_entries_.size(); i++) {
        auto field_id = plan->target_entries_[i];
        auto& field_meta = plan->schema_[field_id];
        std::unique_ptr<milvus::DataArray> result_data = nullptr;
        if (field_meta.is_vector()) {
            auto dim = field_meta.get_dim();
            if (field_meta.get_data_type() == milvus::DataType::VECTOR_FLOAT) {
                auto src_data = final_search_result->output_fields_data_.at(
                        field_id)->vectors().float_vector().data().data();
                auto array = milvus::segcore::CreateVectorDataArrayFrom(src_data + dim * offset_begin, offset_end - offset_begin, field_meta);
                result_data = std::move(array);
            } else if (field_meta.get_data_type() == milvus::DataType::VECTOR_BINARY) {
                auto src_data = final_search_result->output_fields_data_.at(
                        field_id)->vectors().binary_vector().data();
                auto array = milvus::segcore::CreateVectorDataArrayFrom(src_data + (dim / 8) * offset_begin, offset_end - offset_begin, field_meta);
                result_data = std::move(array);
            } else {
                PanicInfo("unsupported vector data type");
            }
            continue;
        }
        switch (field_meta.get_data_type()) {
            case milvus::DataType::BOOL: {
                auto src_data = final_search_result->output_fields_data_.at(
                        field_id)->scalars().bool_data().data().data();
                auto array = milvus::segcore::CreateScalarDataArrayFrom(src_data + offset_begin, offset_end - offset_begin, field_meta);
                result_data = std::move(array);
            }
            case milvus::DataType::INT8:
            case milvus::DataType::INT16:
            case milvus::DataType::INT32: {
                auto src_data = final_search_result->output_fields_data_.at(
                        field_id)->scalars().int_data().data().data();
                auto array = milvus::segcore::CreateScalarDataArrayFrom(src_data + offset_begin, offset_end - offset_begin, field_meta);
                result_data = std::move(array);
            }
            case milvus::DataType::INT64: {
                auto src_data = final_search_result->output_fields_data_.at(
                        field_id)->scalars().long_data().data().data();
                auto array = milvus::segcore::CreateScalarDataArrayFrom(src_data + offset_begin, offset_end - offset_begin, field_meta);
                result_data = std::move(array);
            }
            case milvus::DataType::FLOAT: {
                auto src_data = final_search_result->output_fields_data_.at(
                        field_id)->scalars().float_data().data().data();
                auto array = milvus::segcore::CreateScalarDataArrayFrom(src_data + offset_begin, offset_end - offset_begin, field_meta);
                result_data = std::move(array);
            }
            case milvus::DataType::DOUBLE: {
                auto src_data = final_search_result->output_fields_data_.at(
                        field_id)->scalars().double_data().data().data();
                auto array = milvus::segcore::CreateScalarDataArrayFrom(src_data + offset_begin, offset_end - offset_begin, field_meta);
                result_data = std::move(array);
            }
            case milvus::DataType::VarChar: {
                auto src_data = final_search_result->output_fields_data_.at(
                        field_id)->scalars().string_data().data().data();
                auto array = milvus::segcore::CreateScalarDataArrayFrom(src_data + offset_begin, offset_end - offset_begin, field_meta);
                result_data = std::move(array);
            }
            default: {
                PanicInfo("unsupported vector data type");
            }
        }

        AssertInfo(result_data, "empty output field data");
        search_result_data->mutable_fields_data()->AddAllocated(result_data.release());
    }

    // SearchResultData to blob
    auto size = search_result_data->ByteSize();
    auto buffer = std::vector<char>(size);
    search_result_data->SerializePartialToArray(buffer.data(), size);

    return buffer;
}

CStatus
Marshal(CSearchResultDataBlobs* cSearchResultDataBlobs,
        CSearchResult* c_search_results,
        CSearchPlan c_plan,
        int32_t num_segments,
        int32_t* nq_slice_sizes,
        int32_t num_slices) {
    try {
        // parse search results and get topK, nq
        std::vector<SearchResult*> search_results(num_segments);
        for (int i = 0; i < num_segments; ++i) {
            search_results[i] = static_cast<SearchResult*>(c_search_results[i]);
        }
        AssertInfo(search_results.size() > 0, "empty search result when Marshal");
        auto plan = (milvus::query::Plan*)c_plan;
        auto topK = search_results[0]->topk_;
        auto nq = search_results[0]->num_queries_;

        // Reorganize search results, get result ids, distances and output fields data
        auto final_search_result = ReorganizeSearchResults(search_results, plan, nq, topK);

        // prefix sum, get slices offsets
        AssertInfo(num_slices > 0, "empty nq_slice_sizes is not allowed");
        auto slice_offsets_size = num_slices + 1;
        auto slice_offsets = std::vector<int32_t>(slice_offsets_size);
        slice_offsets[0] = 0;
        slice_offsets[1] = nq_slice_sizes[0] * topK;
        for (int i = 2; i < slice_offsets_size; i++) {
            slice_offsets[i] = slice_offsets[i - 1] + nq_slice_sizes[i - 1] * topK;
        }
        AssertInfo(slice_offsets[num_slices] == nq * topK,
                   "illegal req sizes"
                   ", slice_offsets[last] = " +
                       std::to_string(slice_offsets[num_slices]) + ", nq = " + std::to_string(nq));

        // get search result data blobs by slices
        auto search_result_data_blobs = std::make_unique<milvus::segcore::SearchResultDataBlobs>();
        search_result_data_blobs->blobs.resize(num_slices);
#pragma omp parallel for
        for (int i = 0; i < num_slices; i++) {
            auto proto = GetSearchResultDataSlice(final_search_result.get(), plan, slice_offsets[i], slice_offsets[i + 1], nq, topK);
            search_result_data_blobs->blobs[i] = proto;
        }

        // set final result ptr
        *cSearchResultDataBlobs = search_result_data_blobs.release();
        return milvus::SuccessCStatus();
    } catch (std::exception& e) {
        DeleteSearchResultDataBlobs(cSearchResultDataBlobs);
        return milvus::FailureCStatus(UnexpectedError, e.what());
    }
}

CStatus
GetSearchResultDataBlob(CProto* searchResultDataBlob,
                        CSearchResultDataBlobs cSearchResultDataBlobs,
                        int32_t blob_index) {
    try {
        auto search_result_data_blobs =
            reinterpret_cast<milvus::segcore::SearchResultDataBlobs*>(cSearchResultDataBlobs);
        AssertInfo(blob_index < search_result_data_blobs->blobs.size(), "blob_index out of range");
        searchResultDataBlob->proto_blob = search_result_data_blobs->blobs[blob_index].data();
        searchResultDataBlob->proto_size = search_result_data_blobs->blobs[blob_index].size();
        return milvus::SuccessCStatus();
    } catch (std::exception& e) {
        searchResultDataBlob->proto_blob = nullptr;
        searchResultDataBlob->proto_size = 0;
        return milvus::FailureCStatus(UnexpectedError, e.what());
    }
}

void
DeleteSearchResultDataBlobs(CSearchResultDataBlobs cSearchResultDataBlobs) {
    if (cSearchResultDataBlobs == nullptr) {
        return;
    }
    auto search_result_data_blobs = reinterpret_cast<milvus::segcore::SearchResultDataBlobs*>(cSearchResultDataBlobs);
    delete search_result_data_blobs;
}

CStatus
ReduceSearchResultsAndFillData(CSearchPlan c_plan, CSearchResult* c_search_results, int64_t num_segments) {
    try {
        auto plan = (milvus::query::Plan*)c_plan;
        std::vector<SearchResult*> search_results;
        for (int i = 0; i < num_segments; ++i) {
            search_results.push_back((SearchResult*)c_search_results[i]);
        }
        auto topk = search_results[0]->topk_;
        auto num_queries = search_results[0]->num_queries_;

        // get primary keys for duplicates removal
        for (auto& search_result : search_results) {
            auto segment = (milvus::segcore::SegmentInterface*)(search_result->segment_);
            segment->FillPrimaryKeys(plan, *search_result);
        }

        ReduceResultData(search_results, num_queries, topk);

        // fill in other entities
        for (auto& search_result : search_results) {
            auto segment = (milvus::segcore::SegmentInterface*)(search_result->segment_);
            segment->FillTargetEntry(plan, *search_result);
        }

        auto status = CStatus();
        status.error_code = Success;
        status.error_msg = "";
        return status;
    } catch (std::exception& e) {
        auto status = CStatus();
        status.error_code = UnexpectedError;
        status.error_msg = strdup(e.what());
        return status;
    }
}
