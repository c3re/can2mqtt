use core::fmt;
use std::fmt::Display;

use bytes::Bytes;
use can_socket::CanData;

pub type MQTTPayload = Bytes;
pub type CANFrame = CanData;

pub trait Converter: Display + fmt::Debug {
    fn towards_mqtt(&self, cf: CANFrame) -> Result<MQTTPayload, String>;
    fn towards_can(&self, msg: MQTTPayload) -> Result<CANFrame, String>;
}
