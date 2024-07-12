use tantivy::schema::{IndexRecordOption, Schema, TextFieldIndexing, TextOptions, FAST};
use tantivy::tokenizer::TextAnalyzer;
use tantivy::Index;

use crate::data_type::TantivyDataType;
use crate::tokenizer::default_tokenizer;
use crate::{index_writer::IndexWriterWrapper, log::init_log};

impl IndexWriterWrapper {
    pub(crate) fn from_text_with_tokenizer(
        field_name: String,
        path: String,
        tokenizer_name: String,
        tokenizer: TextAnalyzer,
        num_threads: usize,
        overall_memory_budget_in_bytes: usize,
    ) -> IndexWriterWrapper {
        init_log();

        let mut schema_builder = Schema::builder();
        // positions is required for matching phase.
        let indexing = TextFieldIndexing::default()
            .set_tokenizer(&tokenizer_name)
            .set_index_option(IndexRecordOption::WithFreqsAndPositions);
        let option = TextOptions::default().set_indexing_options(indexing);
        let field = schema_builder.add_text_field(&field_name, option);
        let id_field = schema_builder.add_i64_field("doc_id", FAST);
        let schema = schema_builder.build();
        let index = Index::create_in_dir(path.clone(), schema).unwrap();
        index.tokenizers().register(&tokenizer_name, tokenizer);
        let index_writer = index
            .writer_with_num_threads(num_threads, overall_memory_budget_in_bytes)
            .unwrap();
        let data_type = TantivyDataType::Text;

        IndexWriterWrapper {
            field_name,
            field,
            data_type,
            path,
            index_writer,
            id_field,
        }
    }

    pub(crate) fn from_text_default(
        field_name: String,
        path: String,
        num_threads: usize,
        overall_memory_budget_in_bytes: usize,
    ) -> IndexWriterWrapper {
        IndexWriterWrapper::from_text_with_tokenizer(
            field_name,
            path,
            "default".to_string(),
            default_tokenizer(),
            num_threads,
            overall_memory_budget_in_bytes,
        )
    }
}
