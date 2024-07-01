use std::collections::HashMap;

use log::{info, warn};
use tantivy::tokenizer::{TextAnalyzer, TokenizerManager};

pub(crate) fn default_tokenizer() -> TextAnalyzer {
    TokenizerManager::default().get("default").unwrap()
}

pub(crate) fn create_tokenizer(params: &HashMap<String, String>) -> Option<TextAnalyzer> {
    match params.get("tokenizer") {
        Some(tokenizer_name) => match tokenizer_name.as_str() {
            "default" => {
                return Some(default_tokenizer());
            }
            _ => {
                return None;
            }
        },
        None => {
            info!("no tokenizer is specific, use default tokenizer");
            return Some(default_tokenizer());
        }
    }
}
