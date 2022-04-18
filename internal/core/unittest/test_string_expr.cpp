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

#include <gtest/gtest.h>
#include <memory>
#include <boost/format.hpp>
#include <regex>

#include "pb/plan.pb.h"
#include "query/Expr.h"
#include "query/Plan.h"
#include "query/PlanNode.h"
#include "query/generated/ExprVisitor.h"
#include "query/generated/PlanNodeVisitor.h"
#include "query/generated/ShowPlanNodeVisitor.h"
#include "query/generated/ExecExprVisitor.h"
#include "segcore/SegmentGrowingImpl.h"
#include "test_utils/DataGen.h"
#include "utils/Utils.h"
#include "segcore/SegmentSealed.h"
#include "query/PlanProto.h"

using namespace milvus;

namespace {
auto
GenColumnInfo(int64_t field_id, proto::schema::DataType field_type, bool auto_id, bool is_pk) {
    auto column_info = new proto::plan::ColumnInfo();
    column_info->set_field_id(field_id);
    column_info->set_data_type(field_type);
    column_info->set_is_autoid(auto_id);
    column_info->set_is_primary_key(is_pk);
    return column_info;
}

auto
GenQueryInfo(int64_t topk, std::string metric_type, std::string search_params, int64_t round_decimal = -1) {
    auto query_info = new proto::plan::QueryInfo();
    query_info->set_topk(topk);
    query_info->set_metric_type(metric_type);
    query_info->set_search_params(search_params);
    query_info->set_round_decimal(round_decimal);
    return query_info;
}

auto
GenAnns(bool is_binary, int64_t field_id, std::string placeholder_tag = "$0") {
    auto anns = new proto::plan::VectorANNS();
    anns->set_is_binary(is_binary);
    anns->set_field_id(field_id);
    anns->set_placeholder_tag(placeholder_tag);
    return anns;
}

template <typename T>
auto
GenTermExpr(const std::vector<T>& values) {
    auto term_expr = new proto::plan::TermExpr();
    for (int i = 0; i < values.size(); i++) {
        auto add_value = term_expr->add_values();
        if constexpr (std::is_same_v<T, bool>) {
            add_value->set_bool_val(static_cast<T>(values[i]));
        } else if constexpr (std::is_integral_v<T>) {
            add_value->set_int64_val(static_cast<int64_t>(values[i]));
        } else if constexpr (std::is_floating_point_v<T>) {
            add_value->set_float_val(static_cast<double>(values[i]));
        } else if constexpr (std::is_same_v<T, std::string>) {
            add_value->set_string_val(static_cast<T>(values[i]));
        } else {
            static_assert(always_false<T>);
        }
    }
    return term_expr;
}

auto
GenExpr() {
    return new proto::plan::Expr();
}

SchemaPtr
GenTestSchema() {
    auto schema = std::make_shared<Schema>();
    schema->AddDebugField("fvec", DataType::VECTOR_FLOAT, 16, MetricType::METRIC_L2);
    schema->AddDebugField("str", DataType::VarChar);
    auto pk = schema->AddDebugField("int64", DataType::INT64);
    schema->set_primary_field_id(pk);
    return schema;
}
}  // namespace

TEST(StringExpr, Dummy) {
    InsertData data;
    data.set_num_rows(3);

    auto str_f = new proto::schema::StringArray();
    str_f->add_data("aaaa");
    str_f->add_data("bbbb");
    str_f->add_data("cccc");

    auto scalar_f = new proto::schema::ScalarField();
    scalar_f->set_allocated_string_data(str_f);

    auto newly_f = data.add_fields_data();
    newly_f->set_allocated_scalars(scalar_f);

    std::vector<std::string> dst(3, "");
    auto src_data = reinterpret_cast<const std::string*>(data.fields_data(0).scalars().string_data().data().data());
    std::cout << *(src_data + 0) << std::endl;
    std::cout << *(src_data + 1) << std::endl;
    std::cout << *(src_data + 2) << std::endl;
    std::copy_n(src_data, 3, dst.data());
    for (const auto& str : dst) {
        std::cout << str << std::endl;
    }
}

TEST(StringExpr, Term) {
    using namespace milvus::query;
    using namespace milvus::segcore;

    auto schema = GenTestSchema();
    const auto& fvec_meta = schema->operator[](FieldName("fvec"));
    const auto& str_meta = schema->operator[](FieldName("str"));
    const auto& i64_meta = schema->operator[](FieldName("int64"));

    auto gen_term_expr = [&, fvec_meta,
                          str_meta](const std::vector<std::string>& strs) -> std::unique_ptr<proto::plan::PlanNode> {
        auto column_info = GenColumnInfo(str_meta.get_id().get(), proto::schema::DataType::VarChar, false, false);
        auto term_expr = GenTermExpr<std::string>(strs);
        term_expr->set_allocated_column_info(column_info);

        auto expr = GenExpr();
        expr->set_allocated_term_expr(term_expr);

        auto query_info = GenQueryInfo(10, "L2", "{\"nprobe\": 10}", -1);
        auto anns = GenAnns(fvec_meta.get_data_type() == DataType::VECTOR_BINARY, fvec_meta.get_id().get(), "$0");
        anns->set_allocated_query_info(query_info);
        anns->set_allocated_predicates(expr);

        auto plan_node = std::make_unique<proto::plan::PlanNode>();
        plan_node->set_allocated_vector_anns(anns);
        return std::move(plan_node);
    };

    auto vec_2k_3k = []() -> std::vector<std::string> {
        std::vector<std::string> ret;
        for (int i = 2000; i < 3000; i++) {
            ret.push_back(std::to_string(i));
        }
        return ret;
    }();

    std::map<int, std::vector<std::string>> terms = {
            {0, {"2000", "3000"}},
            {1, {"2000"}},
            {2, {"3000"}},
            {3, {}},
            {4, {vec_2k_3k}},
    };

    auto seg = CreateGrowingSegment(schema);
    int N = 1000;
    std::vector<std::string> str_col;
    int num_iters = 100;
    for (int iter = 0; iter < num_iters; ++iter) {
        auto raw_data = DataGen(schema, N, iter);
        auto new_str_col = raw_data.get_col<std::string>(str_meta.get_id());
        str_col.insert(str_col.end(), new_str_col.begin(), new_str_col.end());
        seg->PreInsert(N);
        seg->Insert(iter * N, N, raw_data.row_ids_.data(), raw_data.timestamps_.data(), raw_data.raw_);
    }

    auto seg_promote = dynamic_cast<SegmentGrowingImpl*>(seg.get());
    ExecExprVisitor visitor(*seg_promote, seg_promote->get_row_count(), MAX_TIMESTAMP);
    for (const auto& [i, term] : terms) {
        auto plan_proto = gen_term_expr(term);
        auto plan = ProtoParser(*schema).CreatePlan(*plan_proto);
        auto final = visitor.call_child(*plan->plan_node_->predicate_.value());
        EXPECT_EQ(final.size(), N * num_iters);

        for (int i = 0; i < N * num_iters; ++i) {
            auto ans = final[i];

            auto val = str_col[i];
            auto ref = std::find(term.begin(), term.end(), val) != term.end();
            ASSERT_EQ(ans, ref) << "@" << i << "!!" << val;
        }
    }
}
