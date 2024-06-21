use tantivy::tokenizer::{TextAnalyzer, TokenizerManager};

pub(crate) fn default_tokenizer() -> TextAnalyzer {
	TokenizerManager::default().get("default").unwrap()
}
