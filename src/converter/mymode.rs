use std::fmt;
use crate::converter::types::{CANFrame, Converter, MQTTPayload};

#[derive(Debug)]
pub struct MyModeConverter {}

impl MyModeConverter {
    pub fn new() -> MyModeConverter{
        MyModeConverter{}
    }
}

impl Converter for MyModeConverter{
    fn towards_mqtt(self: &Self, cf: CANFrame) -> Result<MQTTPayload, String> {
        // here you can implement your own convertmode towards mqtt
        Ok(MQTTPayload::copy_from_slice(cf.as_slice()))
    }

    fn towards_can(self: &Self, mut msg: MQTTPayload) -> Result<CANFrame, String> {
        // here you can implement your own convertmode towards can
        msg.truncate(8);
        match CANFrame::try_new(msg.as_ref()) {
            Ok(cf) => Ok(cf),
            Err(e) => Err(e.to_string()),
        }
    }
}

impl fmt::Display for MyModeConverter {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "mymode")
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn back_and_forth() {
        let cv = MyModeConverter{};
        let cf = CANFrame::new([5; 3]);
        match cv.towards_mqtt(cf) {
            Ok(msg) => {
                println!("Successful conversion {:?}", msg);
                // Aaaaand backwards
                match cv.towards_can(msg) {
                    Ok(cf2) => println!("Successful reverse conversion {:?}", cf2),
                    Err(e) => println!("Error in reverse conversion {}", e),
                }
            }
            Err(e) => println!("Error in conversion {}", e),
        }
    }

    #[test]
    fn zero_bytes_to_mqtt() {
        let cv = MyModeConverter{};
        let cf = CANFrame::new([5; 0]);
        match cv.towards_mqtt(cf) {
            Ok(msg) => {
                println!("Successful conversion {:?}", msg);
                // Aaaaand backwards
                match cv.towards_can(msg) {
                    Ok(cf2) => {
                        assert_eq!(cf2.len(), 0);
                    }
                    Err(e) => panic!("Error in reverse conversion {}", e),
                }
            }
            Err(e) => println!("Error in conversion {}", e),
        }
    }

    #[test]
    fn test_name() {
        let cv = MyModeConverter{};
        assert_eq!("mymode", cv.to_string());
    }

    #[test]
    fn too_long_mqtt() {
        let cv = MyModeConverter{};
        let msg = MQTTPayload::copy_from_slice("overlongstring".as_bytes());
        match cv.towards_can(msg) {
            Ok(cf) => {
                assert_eq!(cf.len(), 8);
                assert_eq!(cf[0], 111);
                assert_eq!(cf[1], 118);
                assert_eq!(cf[2], 101);
                assert_eq!(cf[3], 114);
                assert_eq!(cf[4], 108);
                assert_eq!(cf[5], 111);
                assert_eq!(cf[6], 110);
                assert_eq!(cf[7], 103);
            }
            Err(e) => panic!("Error in reverse conversion {}", e),
        }
    }
}
