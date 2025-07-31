use crate::config::ToCanMap;
use crate::config::ToMqttMap;
use can_socket::CanFrame;
use rumqttc::Publish;
/// MQTTMngEvent is the data structure handled by the MQTT Manager                                                                                
///                                                                                                                                               
/// it can have three different states with a fourth one as an idea:    
/// 1. Pairlist (a list of pairs to subscribe and their convertmodes)   
/// 2. Message In (a Message that has been received that needs to be processed now)                                                               
/// 3. Message Out (a Message that shall be send (on behalf of the CAN Manager))                                                                  
/// (4.) Exit, disconnect and shutdown                                                                                                            
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