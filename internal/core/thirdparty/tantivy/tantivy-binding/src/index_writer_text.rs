





use tantivy::schema::{IndexRecordOption, Schema, TextFieldIndexing, TextOptions};
use tantivy::tokenizer::TextAnalyzer;
use tantivy::{Index, SingleSegmentIndexWriter};

use crate::data_type::TantivyDataType;
use crate::tokenizer::default_tokenizer;
use crate::{index_writer::IndexWriterWrapper, log::init_log};

impl IndexWriterWrapper {
    pub(crate) fn from_text_with_tokenizer(
        field_name: String,
        path: String,
        tokenizer_name: String,
        tokenizer: TextAnalyzer,
    ) -> IndexWriterWrapper {
        init_log();

        let mut schema_builder = Schema::builder();
        // positions is required for matching phase.
        let indexing = TextFieldIndexing::default()
            .set_tokenizer(&tokenizer_name)
            .set_index_option(IndexRecordOption::WithFreqsAndPositions);
        let option = TextOptions::default().set_indexing_options(indexing);
        let field = schema_builder.add_text_field(&field_name, option);
        let schema = schema_builder.build();
        let index = Index::create_in_dir(path.clone(), schema).unwrap();
        index.tokenizers().register(&tokenizer_name, tokenizer);
        let index_writer = SingleSegmentIndexWriter::new(index, 15 * 1024 * 1024).unwrap();
        let data_type = TantivyDataType::Text;

        IndexWriterWrapper {
            field_name,
            field,
            data_type,
            path,
            index_writer,
        }
    }

    pub(crate) fn from_text_default(field_name: String, path: String) -> IndexWriterWrapper {
        IndexWriterWrapper::from_text_with_tokenizer(
            field_name,
            path,
            "default".to_string(),
            default_tokenizer(),
        )
    }
}
