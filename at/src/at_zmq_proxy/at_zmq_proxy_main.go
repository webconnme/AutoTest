/**
 * The MIT License (MIT)
 *
 * Copyright (c) 2015 David You <david@webconn.me>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 * 
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 * 
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

﻿package main

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
