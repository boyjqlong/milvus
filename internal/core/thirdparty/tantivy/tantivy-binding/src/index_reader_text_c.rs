use std::ffi::CStr;

use libc::{c_char, c_void};

use crate::{array::RustArray, index_reader::IndexReaderWrapper};

#[no_mangle]
pub extern "C" fn tantivy_match_query(ptr: *mut c_void, query: *const c_char) -> RustArray {
    let real = ptr as *mut IndexReaderWrapper;
    unsafe {
        let c_str = CStr::from_ptr(query);
        let hits = (*real).match_query(c_str.to_str().unwrap());
        RustArray::from_vec(hits)
    }
}
