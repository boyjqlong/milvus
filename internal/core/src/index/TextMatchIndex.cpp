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

#include <boost/uuid/random_generator.hpp>
#include <boost/uuid/uuid_io.hpp>

#include "index/TextMatchIndex.h"

namespace milvus::index {
TextMatchIndex::TextMatchIndex() {
    auto uuid = boost::uuids::random_generator()();
    auto uuid_string = boost::uuids::to_string(uuid);
    auto file_name = std::string("/tmp/") + uuid_string;
    path_ = file_name;
    boost::filesystem::create_directories(path_);
    d_type_ = TantivyDataType::Text;
    std::string field_name = "tmp_text_index";
    wrapper_ = std::make_shared<TantivyIndexWrapper>(field_name.c_str(),
                                                     path_.c_str());
}

TextMatchIndex::TextMatchIndex(const storage::FileManagerContext& ctx) {
    space_ = nullptr;
    schema_ = ctx.fieldDataMeta.field_schema;
    mem_file_manager_ = std::make_shared<MemFileManager>(ctx, ctx.space_);
    disk_file_manager_ = std::make_shared<DiskFileManager>(ctx, ctx.space_);
    std::string field_name =
        std::to_string(disk_file_manager_->GetFieldDataMeta().field_id);
    auto prefix = disk_file_manager_->GetLocalIndexObjectPrefix();
    path_ = prefix;
    boost::filesystem::create_directories(path_);
    d_type_ = TantivyDataType::Text;
    if (tantivy_index_exist(path_.c_str())) {
        LOG_INFO(
            "index {} already exists, which should happen in loading progress",
            path_);
        reader_ = std::make_shared<TantivyIndexWrapper>(path_.c_str());
    } else {
        wrapper_ = std::make_shared<TantivyIndexWrapper>(field_name.c_str(),
                                                         path_.c_str());
    }
}

void
TextMatchIndex::AddText(const std::string& text) {
    wrapper_->add_data(&text, 1);
}

void
TextMatchIndex::AddTexts(size_t n, const std::string* texts) {
    wrapper_->add_data(texts, n);
}

void
TextMatchIndex::Finish() {
    finish();
}

void
TextMatchIndex::CreateReader() {
    AssertInfo(tantivy_index_exist(path_.c_str()),
               "inverted index not exist: {}",
               path_);
    reader_ = std::make_shared<TantivyIndexWrapper>(path_.c_str());
}
}  // namespace milvus::index
