use std::collections::HashMap;
use std::ffi::c_char;
use std::ffi::c_void;
use std::ffi::CStr;

use crate::index_writer::IndexWriterWrapper;
use crate::tokenizer::create_tokenizer;
use crate::util::create_binding;

#[no_mangle]
pub extern "C" fn tantivy_create_default_text_writer(
    field_name: *const c_char,
    path: *const c_char,
) -> *mut c_void {
    let field_name_str = unsafe { CStr::from_ptr(field_name) };
    let path_str = unsafe { CStr::from_ptr(path) };
    let wrapper = IndexWriterWrapper::from_text_default(
        String::from(field_name_str.to_str().unwrap()),
        String::from(path_str.to_str().unwrap()),
    );
    create_binding(wrapper)
}

#[no_mangle]
pub extern "C" fn tantivy_create_text_writer(
    field_name: *const c_char,
    path: *const c_char,
    tokenizer_name: *const c_char,
    tokenizer_params: *mut c_void,
) -> *mut c_void {
    let field_name_str = unsafe { CStr::from_ptr(field_name) };
    let path_str = unsafe { CStr::from_ptr(path) };
    let tokenizer_name_str = unsafe { CStr::from_ptr(tokenizer_name) };
    let analyzer = unsafe {
        let m = tokenizer_params as *const HashMap<String, String>;
        create_tokenizer(&(*m))
    };
    match analyzer {
        Some(text_analyzer) => {
            let wrapper = IndexWriterWrapper::from_text_with_tokenizer(
                String::from(field_name_str.to_str().unwrap()),
                String::from(path_str.to_str().unwrap()),
                String::from(tokenizer_name_str.to_str().unwrap()),
                text_analyzer,
            );
            return create_binding(wrapper);
        },
        None => {
            return std::ptr::null_mut();
        },
    }
}
