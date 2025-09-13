#![no_main]

use libfuzzer_sys::fuzz_target;
use can2mqtt_rs::converter::none::NoneConverter;
use can2mqtt_rs::converter::types::*;

// Start with cargo fuzz run convert_none2can --  -max_total_time=60
// length over 8 are filtered out below
fuzz_target!(|data: &[u8]| {
    if data.len() < 9 { // we a are interested in fuzzing the convertmode not the frame creation
        let cv = NoneConverter{};
        let cf = CANFrame::try_new(data).expect("something went wrong");
        // so far only checking if we create an error and then panic 
        let _ = cv.towards_mqtt(cf).expect("panic");
    }
});

