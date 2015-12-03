package main

import "fmt"

import (
	zmq "github.com/alecthomas/gozmq"
)

const AF_ZMQ_PROXY_XPUB = "ipc:///tmp/at_frame_pub"
const AF_ZMQ_PROXY_XSUB = "ipc:///tmp/at_frame_sub"
const AF_ZMQ_PROXY_PUB  = "ipc:///tmp/at_frame_pub"
const AF_ZMQ_PROXY_SUB  = "ipc:///tmp/at_frame_sub"

// const AF_ZMQ_PROXY_XPUB = "tcp://*:6000"
// const AF_ZMQ_PROXY_XSUB = "tcp://*:6001"
// const AF_ZMQ_PROXY_PUB  = "tcp://*:6000"
// const AF_ZMQ_PROXY_SUB  = "tcp://*:6001"


//---------------------------------------------------------------------------------------------------------------------
//   
//   AT 용 ZMQ 프록시 서버 앱이다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func main() {

    fmt.Println( "Proxy Start..." )
	zmqContext, _ := zmq.NewContext()
    defer zmqContext.Close()
	
    sub, _ := zmqContext.NewSocket(zmq.SUB)
	defer sub.Close()
	
	pub, _ := zmqContext.NewSocket(zmq.PUB)
	defer pub.Close()

//	sub.Connect( "ipc:///tmp/at_frame_sub" )
//	pub.Connect( "ipc:///tmp/at_frame_pub" )
	pub.Bind( AF_ZMQ_PROXY_XPUB )

	sub.Bind( AF_ZMQ_PROXY_XSUB )
    sub.SetSubscribe("")	
	
	zmq.Device(zmq.FORWARDER, sub, pub )
//	for {
//		message, _ := sub.Recv(0)
//        fmt.Println( "Proxy Read Ok" );		
//		pub.Send(message, 0)
//	}
	
}
