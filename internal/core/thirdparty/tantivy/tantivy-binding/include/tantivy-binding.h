#pragma once

#include <cstdarg>
#include <cstdint>
#include <cstdlib>
#include <ostream>
#include <new>

enum class TantivyDataType : uint8_t {
  Text,
  Keyword,
  I64,
  F64,
  Bool,
};

struct RustArray {
  uint32_t *array;
  size_t len;
  size_t cap;
};

extern "C" {

void free_rust_array(RustArray array);

void print_vector_of_strings(const char *const *ptr, uintptr_t len);

void *tantivy_load_index(const char *path);

void tantivy_free_index_reader(void *ptr);

uint32_t tantivy_index_count(void *ptr);

RustArray tantivy_term_query_i64(void *ptr, int64_t term);

RustArray tantivy_lower_bound_range_query_i64(void *ptr, int64_t lower_bound, bool inclusive);

RustArray tantivy_upper_bound_range_query_i64(void *ptr, int64_t upper_bound, bool inclusive);

RustArray tantivy_range_query_i64(void *ptr,
                                  int64_t lower_bound,
                                  int64_t upper_bound,
                                  bool lb_inclusive,
                                  bool ub_inclusive);

RustArray tantivy_term_query_f64(void *ptr, double term);

RustArray tantivy_lower_bound_range_query_f64(void *ptr, double lower_bound, bool inclusive);

RustArray tantivy_upper_bound_range_query_f64(void *ptr, double upper_bound, bool inclusive);

RustArray tantivy_range_query_f64(void *ptr,
                                  double lower_bound,
                                  double upper_bound,
                                  bool lb_inclusive,
                                  bool ub_inclusive);

RustArray tantivy_term_query_bool(void *ptr, bool term);

RustArray tantivy_term_query_keyword(void *ptr, const char *term);

RustArray tantivy_lower_bound_range_query_keyword(void *ptr,
                                                  const char *lower_bound,
                                                  bool inclusive);

RustArray tantivy_upper_bound_range_query_keyword(void *ptr,
                                                  const char *upper_bound,
                                                  bool inclusive);

RustArray tantivy_range_query_keyword(void *ptr,
                                      const char *lower_bound,
                                      const char *upper_bound,
                                      bool lb_inclusive,
                                      bool ub_inclusive);

RustArray tantivy_prefix_query_keyword(void *ptr, const char *prefix);

RustArray tantivy_regex_query(void *ptr, const char *pattern);

RustArray tantivy_match_query(void *ptr, const char *query);

void tantivy_register_tokenizer(void *ptr, const char *tokenizer_name, void *tokenizer_params);

void *tantivy_create_index(const char *field_name, TantivyDataType data_type, const char *path);

void tantivy_free_index_writer(void *ptr);

void tantivy_finish_index(void *ptr);

void tantivy_index_add_int8s(void *ptr, const int8_t *array, uintptr_t len);

void tantivy_index_add_int16s(void *ptr, const int16_t *array, uintptr_t len);

void tantivy_index_add_int32s(void *ptr, const int32_t *array, uintptr_t len);

void tantivy_index_add_int64s(void *ptr, const int64_t *array, uintptr_t len);

void tantivy_index_add_f32s(void *ptr, const float *array, uintptr_t len);

void tantivy_index_add_f64s(void *ptr, const double *array, uintptr_t len);

void tantivy_index_add_bools(void *ptr, const bool *array, uintptr_t len);

void tantivy_index_add_string(void *ptr, const char *s);

void tantivy_index_add_multi_int8s(void *ptr, const int8_t *array, uintptr_t len);

void tantivy_index_add_multi_int16s(void *ptr, const int16_t *array, uintptr_t len);

void tantivy_index_add_multi_int32s(void *ptr, const int32_t *array, uintptr_t len);

void tantivy_index_add_multi_int64s(void *ptr, const int64_t *array, uintptr_t len);

void tantivy_index_add_multi_f32s(void *ptr, const float *array, uintptr_t len);

void tantivy_index_add_multi_f64s(void *ptr, const double *array, uintptr_t len);

void tantivy_index_add_multi_bools(void *ptr, const bool *array, uintptr_t len);

void tantivy_index_add_multi_keywords(void *ptr, const char *const *array, uintptr_t len);

void *tantivy_create_default_text_writer(const char *field_name, const char *path);

void *tantivy_create_text_writer(const char *field_name,
                                 const char *path,
                                 const char *tokenizer_name,
                                 void *tokenizer_params);

bool tantivy_index_exist(const char *path);

void *create_hashmap();

void hashmap_set_value(void *map, const char *key, const char *value);

void free_hashmap(void *map);

void *tantivy_create_token_stream(void *tokenizer, const char *text);

void tantivy_free_token_stream(void *token_stream);

bool tantivy_token_stream_advance(void *token_stream);

const char *tantivy_token_stream_get_token(void *token_stream);

void *tantivy_create_tokenizer(void *tokenizer_params);

void tantivy_free_tokenizer(void *tokenizer);

void free_rust_string(const char *ptr);

} // extern "C"
