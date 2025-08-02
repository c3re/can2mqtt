use crate::config::ToCanMap;
use crate::config::ToMqttMap;
use can_socket::CanFrame;
use rumqttc::Publish;
use ctflag::Flags;
/// MQTTMngEvent is the data structure handled by the MQTT Manager                                                                                
///                                                                                                                                               
/// it can have three different states with a fourth one as an idea:    
/// 1. Pairlist (a list of pairs to subscribe and their convertmodes)   
/// 2. Message In (a Message that has been received that needs to be processed now)                                                               
/// 3. Message Out (a Message that shall be send (on behalf of the CAN Manager))                                                                  
#[derive(Debug)]
pub enum MQTTMngEvent {                                                                                                                           
    Config(Box<ToCanMap>), // new config to be used                                                                                       
    RX(Publish),                   // Message received (to be converted)                                                                          
    TX(Publish),                   // Frame converted (to be sent)                                                                                
}                                                                             

#[derive(Debug)]
pub enum CANMngEvent {                                                                                                                           
    Config(Box<ToMqttMap>), // new config to be used                                                                                       
    RX(CanFrame),                   // Message received (to be converted)                                                                          
    TX(CanFrame),                   // Frame converted (to be sent)                                                                                
}                                                                             

// Config type
#[derive(Flags, Debug)]
pub struct C2MFlags {
    #[flag(desc = "which config file to use", short='f', placeholder="", default = "can2mqtt.csv")]
    pub file: String,
    #[flag(desc = "which CAN interface to use", short='c', placeholder="", default = "can0")]
    pub can_interface: String,
    #[flag(desc = "which mqtt-broker to use. Example: tcp://user:password@broker.hivemq.com:1883", short='m', placeholder="", default = "tcp://localhost:1883")]
    pub mqtt_connection: String,
    #[flag(desc = "show (very) verbose debug log", short='v', placeholder="", default = false)]
    pub verbose_output: bool,
    #[flag(desc = "direction mode: 0 - bidirectional, 1 - can2mqtt only, 2 - mqtt2can only", short='d', placeholder="", default = 0)]
    pub dir_mode: usize,
}