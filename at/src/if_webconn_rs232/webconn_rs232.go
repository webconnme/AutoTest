package main
//
//import (
//    "fmt"
//    "time"
//    "runtime"
//)
//
import (
//	"github.com/mikepb/go-serial"
	zmq "github.com/alecthomas/gozmq"
//	"sync"
	"sync"
	"runtime"
	"time"
	"fmt"
	"encoding/json"
)

type cmd struct {
	command string
	data string
}


//
//var ThreadRS232RxReqEnd bool    = false
//var ThreadRS232TxReqEnd bool    = false
//var ThreadRS232RxRun    bool    = false
//var ThreadRS232TxRun    bool    = false
var ThreadRS232Live   uint64  = 0

var PortRxPUB    *zmq.Socket = nil
var PortTxSUB    *zmq.Socket = nil

var PairSocket   *zmq.Socket = nil

var stopRX = make(chan bool)
var stopTX = make(chan bool)

func ThreadWebConnRS232Rx(wg *sync.WaitGroup) {
	wg.Done()

	ad.Println( "ThreadWebConnRS232Rx() start" )
	for {
		ThreadRS232Live++
		runtime.Gosched()

		select {
		case <-stopRX:
			ad.Println( "ThreadWebConnRS232Rx() End" )
			return
		default:
			pi := zmq.PollItems{
				zmq.PollItem{ Socket: PairSocket, Events: zmq.POLLIN},
			}

			event_count, err := zmq.Poll( pi, 1 * time.Millisecond )
			if err != nil {
				ad.Println( "fail do not poll[%s]", err )
				reason := fmt.Sprintf( "do not poll[%s]", err )
				ar.SetResultError( reason )
				break

			}

			if event_count == 0 {

			} else {
				if pi[0].REvents&zmq.POLLIN != 0 {
					buf, err := pi[0].Socket.Recv(0)

					if err != nil {
						ad.Println( "fail do not read [%s]", err )
						reason := fmt.Sprintf( "do not read [%s]", err )
						ar.SetResultError( reason )
						break;
					}

					var data []map[string]string

					json.Unmarshal(buf, &data)

					for _, m := range data {
						PortRxPUB.Send( []byte(m["data"]), 0 )

						ad.Println( "OUT JSON : [%v]\n", m["data"] )
					}
				}
			}
		}
	}

}

////
////func ThreadRS232Rx() {
////
////	ad.Println( "ThreadRS232Rx() start" )
////
////	ThreadRS232RxReqEnd = false
////    ThreadRS232RxRun    = true
////
////	for !ThreadRS232RxReqEnd {
////
////	    ThreadRS232Live++
////	    runtime.Gosched()
////
////        remain, err := RS232Port.InputWaiting()
//// 		if err != nil {
////            ad.Println( "fail do not read wait[%s]", err )
////	        reason := fmt.Sprintf( "do not read wait[%s]", err )
////	        ar.SetResultError( reason )
//// 			break;
//// 		}
////
//// 		if remain != 0 {
////
////			ad.Println( "remain = %d", remain );
////
////			for i:=0; i < remain; i++ {
//// 		        buf := make([]byte, 1)
////
//// 		        length, err := RS232Port.Read(buf)
//// 		        if err != nil {
////                    ad.Println( "fail do not read wait[%s]", err )
////	                reason := fmt.Sprintf( "do not read wait[%s]", err )
////	                ar.SetResultError( reason )
//// 		        	break;
//// 		        }
//////			    ad.Println( "read length [%d]", length );
////
////			    err = PortRxPUB.Send(buf[:length], 0)
////                if err != nil {
////                    ad.Println( "fail do not send channel[%s]", err )
////	                reason := fmt.Sprintf( "do not send channel[%s]", err )
////	                ar.SetResultError( reason )
//// 		        	break;
////                }
////			}
////
////// 		    buf := make([]byte, remain)
//////
////// 		    length, err := RS232Port.Read(buf)
////// 		    if err != nil {
//////                ad.Println( "fail do not read wait[%s]", err )
//////	            reason := fmt.Sprintf( "do not read wait[%s]", err )
//////	            ar.SetResultError( reason )
////// 		    	break;
////// 		    }
//////			ad.Println( "read length [%d]", length );
//////
//////			err = PortRxPUB.Send(buf[:length], 0)
//////            if err != nil {
//////                ad.Println( "fail do not send channel[%s]", err )
//////	            reason := fmt.Sprintf( "do not send channel[%s]", err )
//////	            ar.SetResultError( reason )
////// 		    	break;
//////            }
////
//// 		}
////
////	}
////
////	ad.Println( "ThreadRS232Rx() End" )
////    ThreadRS232RxRun    = false
////
////}
////

type Message struct {
	Command string `json:"command"`
	Data string `json:"data"`
}
func ThreadWebConnRS232Tx(wg *sync.WaitGroup) {
	wg.Done()

	ad.Println( "ThreadWebConnRS232Tx() start" )

	for {
		select {
		case <-stopTX:
			ad.Println( "ThreadWebConnRS232Tx() End" )
			return
		default:
			// ad.Println( "wait read PortTxSUB" )
			pi := zmq.PollItems{
				 zmq.PollItem{ Socket: PortTxSUB, Events: zmq.POLLIN},
			}

			event_count, err := zmq.Poll( pi, 1 * time.Millisecond )
			if err != nil {
				ad.Println( "fail do not poll[%s]", err )
				reason := fmt.Sprintf( "do not poll[%s]", err )
				ar.SetResultError( reason )
				break

			}

			if event_count == 0 {

			} else {
				if pi[0].REvents&zmq.POLLIN != 0 {
					buf, err := pi[0].Socket.Recv(0)
					if err != nil {
						 ad.Println( "fail do not read [%s]", err )
						 reason := fmt.Sprintf( "do not read [%s]", err )
						 ar.SetResultError( reason )
						 break;
					}

					msg := Message{
						Command: "tx",
						Data: string(buf)}

					data, err := json.Marshal([]Message{msg})
					if err != nil {
						ad.Println("JSON PARSE: " + err.Error())
					}

					PairSocket.Send( data, zmq.NOBLOCK )

					str := string(data)
					ad.Println( "OUT ASCII : [%s]\n", str )
				}
			}
		}
	}
}


////func ThreadRS232Tx() {
////
////	ad.Println( "ThreadRS232Tx() start" )
////
////	ThreadRS232TxReqEnd = false
////    ThreadRS232TxRun    = true
////
////	for !ThreadRS232TxReqEnd {
////
////	    ThreadRS232Live++
////
//////		ad.Println( "wait read PortTxSUB" )
////        pi := zmq.PollItems{
////    	     zmq.PollItem{ Socket: PortTxSUB, Events: zmq.POLLIN},
////    	}
////
////        event_count, err := zmq.Poll( pi, 1 * time.Millisecond )
////   		if err != nil {
////            ad.Println( "fail do not poll[%s]", err )
////	        reason := fmt.Sprintf( "do not poll[%s]", err )
////	        ar.SetResultError( reason )
//// 			break
////
////    	}
////
////		if event_count == 0 {
////
////		} else {
////            if pi[0].REvents&zmq.POLLIN != 0 {
////	            buf, err := pi[0].Socket.Recv(0)
////   	            if err != nil {
////                     ad.Println( "fail do not read [%s]", err )
////	                 reason := fmt.Sprintf( "do not read [%s]", err )
////	                 ar.SetResultError( reason )
//// 		    	     break;
////                }
////
////				RS232Port.Write( buf )
////
////		        str := string(buf)
////		        ad.Println( "OUT ASCII : [%s]\n", str )
////			}
////		}
////
////	}
////
////	ad.Println( "ThreadRS232Tx() End" )
////	ThreadRS232TxRun    = false
////
////}
////
////func OpenRS232() {
////	c, err := zmq.NewContext()
////	if err != nil {
////		ad.Println("if_rs232, OpenRS232(), zmq.NewContext() %s", err)
////	}
////
////	s, err := c.NewSocket(zmq.PAIR)
////	if err != nil {
////		ad.Println("if_rs232, OpenRS232(), zmq.Context.NewSocket() %s", err)
////	}
////
////	err = s.Connect(RS232Path)
////	if err != nil {
////		ad.Println("if_rs232, OpenRS232(), zmq.Socket.Connect() %s", err)
////	}
////
////	RS232Port = s
////}
////func OpenRS232() {
////
////	options := RS232Options
////
////	var err error
////
////	if RS232Port != nil {
////	    RS232Port.Close()
////	}
////
////    RS232Port, err = options.Open( RS232Path )
////    if err != nil {
////        ad.Println( "fail do not open RS232[%s]", RS232Path )
////	    reason := fmt.Sprintf( "do not open RS232[%s]", RS232Path )
////	    ar.SetResultError( reason )
////    }
////
////}
//
//func CloseRS232() {
//
//	if RS232Port != nil {
//	    RS232Port.Close()
//	}
//
//}
//
