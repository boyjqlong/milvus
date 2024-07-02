#include <string>
#include <vector>
#include <boost/filesystem/operations.hpp>

#include "tantivy-binding.h"
#include "tantivy-wrapper.h"

using namespace milvus::tantivy;

std::set<uint32_t>
to_set(const RustArrayWrapper& w) {
    std::set<uint32_t> s(w.array_.array, w.array_.array + w.array_.len);
    return s;
}

int
main(int argc, char* argv[]) {
    auto path = "/tmp/inverted-index/text-demo/";
    boost::filesystem::remove_all(path);
    boost::filesystem::create_directories(path);

    std::string tokenizer_name = "jieba";
    std::unordered_map<std::string, std::string> tokenizer_params;
    tokenizer_params["tokenizer"] = tokenizer_name;

    auto text_writer = TantivyIndexWrapper(
        "text_demo", path, tokenizer_name.c_str(), tokenizer_params);
    auto write_single_text = [&text_writer](const std::string& s) {
        text_writer.add_data(&s, 1);
    };

    {
        write_single_text(
            "张华考上了北京大学；李萍进了中等技术学校；我在百货公司当售货员：我"
            "们都有光明的前途");
        write_single_text("测试中文分词器的效果");
        text_writer.finish();
    }

    auto text_reader = TantivyIndexWrapper(path);
    text_reader.register_tokenizer(tokenizer_name.c_str(), tokenizer_params);

    {
        auto result = to_set(text_reader.match_query("北京"));
        assert(result.size() == 1);
        assert(result.find(0) != result.end());
    }

    {
        auto result = to_set(text_reader.match_query("效果"));
        assert(result.size() == 1);
        assert(result.find(1) != result.end());
    }

    return 0;
}
