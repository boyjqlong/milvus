use std::collections::HashMap;
use std::ffi::{CStr};
use std::os::raw::{c_char};

use libc::c_void;

#[no_mangle]
pub extern "C" fn create_hashmap() -> *mut c_void {
    let map: HashMap<String, String> = HashMap::new();
    Box::into_raw(Box::new(map)) as *mut c_void
}

#[no_mangle]
pub extern "C" fn set_value(map: *mut c_void, key: *const c_char, value: *const c_char) {
    let m = map as *mut HashMap<String, String>;
    let k = unsafe { CStr::from_ptr(key).to_str().unwrap() };
    let v = unsafe { CStr::from_ptr(value).to_str().unwrap() };
    unsafe {
        (*m).insert(String::from(k), String::from(v));
    }
}

#[no_mangle]
pub extern "C" fn free_hashmap(map: *mut c_void) {
    let _ = unsafe { Box::from_raw(map as *mut HashMap<String, String>) };
}
