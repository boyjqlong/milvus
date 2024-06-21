use tantivy::{query::BooleanQuery, tokenizer::TokenStream, Term};

use crate::{index_reader::IndexReaderWrapper, tokenizer::default_tokenizer};

impl IndexReaderWrapper {
    // split the query string into multiple tokens using index's default tokenizer,
    // and then execute the conjunction of term query.
    pub fn match_query(&self, q: &str) -> Vec<u32> {
        // clone the tokenizer to make `match_query` thread-safe.
        let mut tokenizer = self
            .index
            .tokenizer_for_field(self.field)
			.unwrap_or(default_tokenizer()) // TODO: register the runtime tokenizer.
            .clone();
        let mut token_stream = tokenizer.token_stream(q);
        let mut terms: Vec<Term> = Vec::new();
        while token_stream.advance() {
            let token = token_stream.token();
            terms.push(Term::from_field_text(self.field, &token.text));
        }
        let query = BooleanQuery::new_multiterms_query(terms);
        self.search(&query)
    }
}
