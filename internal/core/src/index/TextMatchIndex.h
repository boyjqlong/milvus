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

#pragma once

#include <string>
#include <boost/filesystem.hpp>

#include "index/InvertedIndexTantivy.h"

namespace milvus::index {
class TextMatchIndex : public InvertedIndexTantivy<std::string> {
 public:
    explicit TextMatchIndex();
    explicit TextMatchIndex(const storage::FileManagerContext& ctx);

 public:
    void
    AddText(const std::string& text);

    void
    AddTexts(size_t n, const std::string* texts);

    void
    Finish();

 public:
    void
    CreateReader();

    TargetBitmap
    MatchQuery(const std::string& query);

 private:
    std::shared_ptr<TantivyIndexWrapper> reader_;
};
}  // namespace milvus::index
