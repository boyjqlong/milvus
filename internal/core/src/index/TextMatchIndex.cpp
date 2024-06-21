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

#include "index/TextMatchIndex.h"

namespace milvus::index {
TextMatchIndex::TextMatchIndex(const storage::FileManagerContext& ctx) {
    space_ = nullptr;
    schema_ = ctx.fieldDataMeta.field_schema;
    mem_file_manager_ = std::make_shared<MemFileManager>(ctx, ctx.space_);
    disk_file_manager_ = std::make_shared<DiskFileManager>(ctx, ctx.space_);
    auto field =
        std::to_string(disk_file_manager_->GetFieldDataMeta().field_id);
    auto prefix = disk_file_manager_->GetLocalIndexObjectPrefix();
    path_ = prefix;
    boost::filesystem::create_directories(path_);
    d_type_ = TantivyDataType::Text;
    if (tantivy_index_exist(path_.c_str())) {
        LOG_INFO(
            "index {} already exists, which should happen in loading progress",
            path_);
    } else {
        wrapper_ =
            std::make_shared<TantivyIndexWrapper>(field.c_str(), path_.c_str());
    }
}

}  // namespace milvus::index
