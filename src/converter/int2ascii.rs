use std::fmt;
use crate::converter::types::{CANFrame, Converter, MQTTPayload};
use bytes::BufMut;


#[derive(Debug)]
pub struct IntConverter {
    instances: usize, // these fantastically private fields enforce the usage of the constructor :)
    bits: usize,
}

impl IntConverter {
    pub fn new(instances: usize, bits: usize) -> Result<IntConverter, String> {
        match instances {
            1 | 2 | 4 | 8 => {}
            _ => {
                return Err(format!(
                    "Invalid instance size, allowed values: 1, 2, 4, 8. got {}",
                    instances
                ));
            }
        }
        match bits {
            8 | 16 | 32 | 64 => {}
            _ => {
                return Err(format!(
                    "Invalid bit size, allowed values: 8, 16, 32, 64. got {}",
                    bits
                ));
            }
        }

        if (bits / 8) * instances > 8 {
            return Err(format!(
                "{} instances of {} exceed a CAN-Frame (8 byte)",
                instances, bits
            ));
        }

        Ok(IntConverter { instances, bits })
    }

    fn expected_len(self: &Self) -> usize {
        (self.bits / 8) * self.instances
    }
}

impl Converter for IntConverter {
    fn towards_mqtt(self: &Self, cf: CANFrame) -> Result<MQTTPayload, String> {
        if cf.len() != self.expected_len() {
            return Err(format!(
                "Length mismatch, expected {} bytes, got {} bytes",
                self.expected_len(),
                cf.len()
            ));
        }

        let bytesize = self.bits / 8;

        // I chose i64 here because i8, i16, i32, and i64 will fit into it,
        // values still should not exceed the size because I feed it 1, 2, 4, and 8 bytes respectively
        let mut r: Vec<i64> = vec![];
        match self.bits {
            8 => {
                for i in 0..self.instances {
                    r.push(i64::from(i8::from_le_bytes(
                        cf[bytesize * i..bytesize * (i + 1)].try_into().unwrap(),
                    )));
                }
            }
            16 => {
                for i in 0..self.instances {
                    r.push(i64::from(i16::from_le_bytes(
                        cf[bytesize * i..bytesize * (i + 1)].try_into().unwrap(),
                    )));
                }
            }
            32 => {
                for i in 0..self.instances {
                    r.push(i64::from(i32::from_le_bytes(
                        cf[bytesize * i..bytesize * (i + 1)].try_into().unwrap(),
                    )));
                }
            }
            _ => {
                // 64 case is the only possibility here, we have restricted (immutable) bits to 8,16,32,64 in the constructor
                for i in 0..self.instances {
                    r.push(i64::from_le_bytes(
                        cf[bytesize * i..bytesize * (i + 1)].try_into().unwrap(),
                    ));
                }
            }
        }

        let res = r
            .iter()
            .map(|i| format!("{}", *i))
            .collect::<Vec<String>>()
            .join(" ");

        Ok(MQTTPayload::from(res))
    }

    fn towards_can(self: &Self, msg: MQTTPayload) -> Result<CANFrame, String> {
        if !msg.is_ascii() {
            return Err(("input contains non-ASCII characters").to_string());
        }

        let str = String::from(msg.escape_ascii().to_string());

        let number_strs = str.split(" ").collect::<Vec<_>>();

        if number_strs.len() != self.instances {
            return Err(format!(
                "Wrong amount of instances, expected {}, got {}",
                self.instances,
                number_strs.len()
            ));
        }

        let mut numbers: Vec<u8> = vec![];
        match self.bits {
            8 => {
                for str in number_strs {
                    match str.parse::<i8>() {
                        Ok(i) => numbers.put_i8(i),
                        Err(e) => return Err(format!("Error parsing number: {}", e.to_string())),
                    }
                }
            }
            16 => {
                for str in number_strs {
                    match str.parse::<i16>() {
                        Ok(i) => numbers.put_i16_le(i),
                        Err(e) => return Err(format!("Error parsing number: {}", e.to_string())),
                    }
                }
            }
            32 => {
                for str in number_strs {
                    match str.parse::<i32>() {
                        Ok(i) => numbers.put_i32_le(i),
                        Err(e) => return Err(format!("Error parsing number: {}", e.to_string())),
                    }
                }
            }
            _ => {
                for str in number_strs {
                    match str.parse::<i64>() {
                        Ok(i) => numbers.put_i64_le(i),
                        Err(e) => return Err(format!("Error parsing number: {}", e.to_string())),
                    }
                }
            }
        }
        match CANFrame::try_new(&numbers[0..self.expected_len()]) {
            Ok(cf) => Ok(cf),
            Err(e) => Err(e.to_string()),
        }
    }

}

impl fmt::Display for IntConverter {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        let instance_string = match self.instances {
            1 => "".to_string(),
            i => format!("{}", i),
        };
        write!(f, "{}int{}2ascii", instance_string, self.bits )
    }
}


#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_new() {
        // good cases
        IntConverter::new(1, 8).expect("should not err, 1*8 fits");
        IntConverter::new(2, 8).expect("should not err, 2*8 fits");
        IntConverter::new(4, 8).expect("should not err, 4*8 fits");
        IntConverter::new(8, 8).expect("should not err, 8*8 fits");

        IntConverter::new(1, 16).expect("should not err, 1*16 fits");
        IntConverter::new(2, 16).expect("should not err, 2*16 fits");
        IntConverter::new(4, 16).expect("should not err, 4*16 fits");

        IntConverter::new(1, 32).expect("should not err, 1*32 fits");
        IntConverter::new(2, 32).expect("should not err, 2*32 fits");

        IntConverter::new(1, 64).expect("should not err, 1*64 fits");

        // invalid instance
        IntConverter::new(3, 8).expect_err("should err, 3 is no valid instance count");

        // invalid bitcount
        IntConverter::new(1, 7).expect_err("should err, 7 is no valid bit count");

        // invalid total size
        IntConverter::new(2, 64).expect_err("should err, 2*64 exceeds CAN Frame size of 8 byte");
    }

    #[test]
    fn test_name() {
        // implicit "1"
        assert_eq!(IntConverter::new(1, 8).unwrap().to_string(), "int82ascii");
        assert_eq!(IntConverter::new(2, 8).unwrap().to_string(), "2int82ascii");
        assert_eq!(IntConverter::new(2, 32).unwrap().to_string(), "2int322ascii");
    }

    #[test]
    fn test_towards_mqtt18() {
        let cv = IntConverter::new(1, 8).unwrap();
        let cf = CANFrame::new([1]);
        assert_eq!(cv.towards_mqtt(cf).unwrap(), "1");
    }

    #[test]
    fn test_towards_mqtt28() {
        let cv = IntConverter::new(2, 8).unwrap();
        let cf = CANFrame::new([1, 2]);
        assert_eq!(cv.towards_mqtt(cf).unwrap(), "1 2");
    }

    #[test]
    fn test_towards_mqtt28e() {
        let cv = IntConverter::new(2, 8).unwrap();
        let cf = CANFrame::new([1, 2, 3]);
        cv.towards_mqtt(cf).expect_err("should err, wrong length");
    }

    #[test]
    fn test_towards_mqtt416() {
        let cv = IntConverter::new(4, 16).unwrap();
        let cf = CANFrame::new([1, 0, 2, 0, 3, 0, 5, 0]);
        assert_eq!(cv.towards_mqtt(cf).unwrap(), "1 2 3 5");
    }

    #[test]
    fn test_towards_mqtt232() {
        let cv = IntConverter::new(2, 32).unwrap();
        let cf = CANFrame::new([1, 0, 0, 0, 2, 0, 0, 0]);
        assert_eq!(cv.towards_mqtt(cf).unwrap(), "1 2");
    }

    #[test]
    fn test_towards_can88() {
        let cv = IntConverter::new(8, 8).unwrap();
        let msg = MQTTPayload::copy_from_slice("4 2 3 4 5 6 7 8".as_bytes());
        let cf = cv.towards_can(msg).unwrap();
        assert_eq!(cf, [4, 2, 3, 4, 5, 6, 7, 8]);
    }

    #[test]
    fn test_towards_can88e() {
        let cv = IntConverter::new(8, 8).unwrap();
        let msg = MQTTPayload::copy_from_slice("128 2 3 4 5 6 7 8".as_bytes());
        cv.towards_can(msg)
            .expect_err("should err, 128 too large for a i8");
    }

    #[test]
    fn test_towards_can88e2() {
        let cv = IntConverter::new(8, 8).unwrap();
        let msg = MQTTPayload::copy_from_slice("-1 2 3 4 5 6 7 8".as_bytes());
        let mut expected_res: Vec<u8> = vec![];
        expected_res.put_i8(-1);
        expected_res.put_i8(2);
        expected_res.put_i8(3);
        expected_res.put_i8(4);
        expected_res.put_i8(5);
        expected_res.put_i8(6);
        expected_res.put_i8(7);
        expected_res.put_i8(8);

        assert_eq!(expected_res[0..8], cv.towards_can(msg).unwrap());
    }

    #[test]
    fn test_towards_can88e22() {
        let cv = IntConverter::new(8, 8).unwrap();
        let msg = MQTTPayload::copy_from_slice("-128 -129 3 4 5 6 7 8".as_bytes());
        cv.towards_can(msg)
            .expect_err("should err, -129 too large for a i8");
    }


    #[test]
    fn test_towards_can88e3() {
        let cv = IntConverter::new(8, 8).unwrap();
        let msg = MQTTPayload::copy_from_slice("2 3 4 5 6 7 8".as_bytes());
        cv.towards_can(msg)
            .expect_err("should err, too few instances");
    }

    #[test]
    fn test_towards_can88e4() {
        let cv = IntConverter::new(8, 8).unwrap();
        let msg = MQTTPayload::copy_from_slice("1 1 2 3 4 5 6 7 8".as_bytes());
        cv.towards_can(msg)
            .expect_err("should err, too many instances");
    }

    #[test]
    fn test_towards_can88e5() {
        let cv = IntConverter::new(8, 8).unwrap();
        let msg = MQTTPayload::copy_from_slice("a 2 3 4 5 6 7 8".as_bytes());
        cv.towards_can(msg)
            .expect_err("should err, non-numbers involved");
    }


    #[test]
    fn test_towards_can416() {
        let cv = IntConverter::new(4,16).unwrap();
        let msg = MQTTPayload::copy_from_slice("1 2 3 4".as_bytes());
        let cf = cv.towards_can(msg).unwrap();
        assert_eq!(cf, [1, 0, 2, 0, 3, 0, 4, 0]);
    }

    #[test]
    fn test_towards_can216() {
        let cv = IntConverter::new(2,16).unwrap();
        let msg = MQTTPayload::copy_from_slice("1 2".as_bytes());
        let cf = cv.towards_can(msg).unwrap();
        assert_eq!(cf, [1, 0, 2, 0]);
    }

}
