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

    auto text_writer = TantivyIndexWrapper("text_demo", path);
    auto write_single_text = [&text_writer](const std::string& s) {
        text_writer.add_data(&s, 1);
    };

    {
        write_single_text("football, basketball, pingpang");
        write_single_text("swimming, football");
        text_writer.finish();
    }

    auto text_reader = TantivyIndexWrapper(path);

    {
        auto result = to_set(text_reader.match_query("football"));
        assert(result.size() == 2);
        assert(result.find(0) != result.end());
        assert(result.find(1) != result.end());
    }

    {
        auto result = to_set(text_reader.match_query("basketball"));
        assert(result.size() == 1);
        assert(result.find(0) != result.end());
    }

    {
        auto result = to_set(text_reader.match_query("swimming"));
        assert(result.size() == 1);
        assert(result.find(1) != result.end());
    }

    {
        auto result = to_set(text_reader.match_query("basketball, swimming"));
        assert(result.size() == 2);
        assert(result.find(0) != result.end());
        assert(result.find(1) != result.end());
    }

    return 0;
}
