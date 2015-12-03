package main

import (
	"fmt"
	"time"
)

import (
	zmq "github.com/alecthomas/gozmq"
)

const AF_ZMQ_PROXY_XPUB = "ipc:///tmp/at_frame_pub"
const AF_ZMQ_PROXY_XSUB = "ipc:///tmp/at_frame_sub"
const AF_ZMQ_PROXY_PUB  = "ipc:///tmp/at_frame_pub"
const AF_ZMQ_PROXY_SUB  = "ipc:///tmp/at_frame_sub"

// const AF_ZMQ_PROXY_XPUB = "tcp://*:6000"
// const AF_ZMQ_PROXY_XSUB = "tcp://*:6001"
// const AF_ZMQ_PROXY_PUB  = "tcp://localhost:6000"
// const AF_ZMQ_PROXY_SUB  = "tcp://localhost:6001"

var (
	zmqContext   *zmq.Context
	cmdSUB       *zmq.Socket      			// 프레임 명령 수신 SUB 소켓 
	cmdPUB       *zmq.Socket      			// 프레임 명령 송신 SUB 소켓 
)

func main() {

    fmt.Println( "This is Test ZMQ SUB\n" )
	
	zmqContext, _ = zmq.NewContext()
	cmdPUB , _    = zmqContext.NewSocket(zmq.PUB)
	cmdSUB , _    = zmqContext.NewSocket(zmq.SUB)	
	
	cmdPUB.Connect( AF_ZMQ_PROXY_SUB )
	
	cmdSUB.Connect( AF_ZMQ_PROXY_PUB  )
	cmdSUB.SetSubscribe("")	
	
	count := 0
	for {
	    str := fmt.Sprintf( "SEND PUB %d" , count );
		fmt.Println( str );
		_ = cmdPUB.Send([]byte( str ), 0)
        time.Sleep( time.Duration(1000) * time.Millisecond)	
		count++
	}
	
	cmdPUB.Close()
	zmqContext.Close()
	
}
