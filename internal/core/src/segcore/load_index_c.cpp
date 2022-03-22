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

#include "common/LoadInfo.h"
#include "exceptions/EasyAssert.h"
#include "knowhere/common/BinarySet.h"
#include "knowhere/index/vector_index/VecIndexFactory.h"
#include "segcore/load_index_c.h"

CStatus
NewLoadIndexInfo(CLoadIndexInfo* c_load_index_info) {
    try {
        auto load_index_info = std::make_unique<LoadIndexInfo>();
        *c_load_index_info = load_index_info.release();
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

void
DeleteLoadIndexInfo(CLoadIndexInfo c_load_index_info) {
    auto info = (LoadIndexInfo*)c_load_index_info;
    delete info;
}

CStatus
AppendIndexParam(CLoadIndexInfo c_load_index_info, const char* c_index_key, const char* c_index_value) {
    try {
        auto load_index_info = (LoadIndexInfo*)c_load_index_info;
        std::string index_key(c_index_key);
        std::string index_value(c_index_value);
        load_index_info->index_params[index_key] = index_value;

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

CStatus
AppendFieldInfo(CLoadIndexInfo c_load_index_info, int64_t field_id) {
    try {
        auto load_index_info = (LoadIndexInfo*)c_load_index_info;
        load_index_info->field_id = field_id;

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

void
DebugBinarySet(const milvus::knowhere::BinarySet* binary_set) {
    std::string ss;
    ss += std::string("len: ") + std::to_string(binary_set->binary_map_.size()) + "\n";
    for (const auto&[k, v] : binary_set->binary_map_) {
        // std::cout << k << ":\n" << v->data.get() << "\n" << v->size << std::endl;
        ss +=  std::string(k) + ":\n" + std::to_string(v->size) + "\n";
    }
    std::cout << ss << std::endl;
}

CStatus
AppendIndex(CLoadIndexInfo c_load_index_info, CBinarySet c_binary_set) {
    try {
        auto load_index_info = (LoadIndexInfo*)c_load_index_info;
        auto binary_set = (milvus::knowhere::BinarySet*)c_binary_set;
        // DebugBinarySet(binary_set);
        auto& index_params = load_index_info->index_params;
        bool find_index_type = index_params.count("index_type") > 0 ? true : false;
        bool find_index_mode = index_params.count("index_mode") > 0 ? true : false;
        AssertInfo(find_index_type == true, "Can't find index type in index_params");
        milvus::knowhere::IndexMode mode;
        if (find_index_mode) {
            mode = index_params["index_mode"] == "CPU" ? milvus::knowhere::IndexMode::MODE_CPU
                                                       : milvus::knowhere::IndexMode::MODE_GPU;
        } else {
            mode = milvus::knowhere::IndexMode::MODE_CPU;
        }
        load_index_info->index =
            milvus::knowhere::VecIndexFactory::GetInstance().CreateVecIndex(index_params["index_type"], mode);
        load_index_info->index->Load(*binary_set);
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
