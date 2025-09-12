#![no_main]

use libfuzzer_sys::fuzz_target;
use can2mqtt::converter::bytecolor2colorcode::ByteColor2ColorCodeConverter;
use can2mqtt::converter::types::*;
use regex::Regex;

fuzz_target!(|data: &[u8]| {
    let cv = ByteColor2ColorCodeConverter{};
    let msg = MQTTPayload::copy_from_slice(data);
    let res = cv.towards_can(msg.clone());
    let re = Regex::new(r"^#?(?i)[0-9a-f]{6}$").unwrap();
    match res {
        Err(_) => {
            if msg.is_ascii() {
                // leading # is optional
                // (?i) turns on case insensitivity for the hexchars.
                if re.is_match(&String::from(msg.escape_ascii().to_string())) {
                    panic!("should not panic");
                }
            }
        },
        Ok(_) => {
                if !re.is_match(&String::from(msg.escape_ascii().to_string())) {
                    panic!("should not be ok");
                }
        },
    };
    } 
);
