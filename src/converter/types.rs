use core::fmt;
use std::fmt::Display;

use bytes::Bytes;
use can_socket::CanData;

pub type MQTTPayload = Bytes;
pub type CANFrame = CanData;

pub trait Converter: Display + fmt::Debug {
    fn towards_mqtt(self: &Self, cf: CANFrame) -> Result<MQTTPayload, String>;
    fn towards_can(self: &Self, msg: MQTTPayload) -> Result<CANFrame, String>;
}
