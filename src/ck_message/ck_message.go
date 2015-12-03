package main

import (
    "fmt"
	"time"
	"strconv"
	"bytes"
_    "runtime"
)

import (
 	zmq "github.com/alecthomas/gozmq"
)

var ThreadCheckMsgReqEnd bool    = false
var ThreadCheckMsgRun    bool    = false
var ThreadCheckMsgLive   uint64  = 0 

var PortInAsciiSUB    *zmq.Socket = nil   // 분석 문자열 수신 SUB 소켓 
var PortOutAsciiPUB   *zmq.Socket = nil   // 문자열 송신 SUB 소켓 

var commands []CkmCommand 
var cmdIndex  int
var curCommand CkmCommand 
var timeOut    int

var checkMsg []byte
var currentMsg []byte

func SleepCheckMsg(  sleep_time int ) {
    
	ad.Println( "SleepCheckMsg() start" )

	// 1 m sec 마다 끝났는가를 확인한다. 
	start_time    := time.Now()
	time_out_msec := time.Duration( sleep_time ) * time.Millisecond 

	for !ThreadCheckMsgReqEnd {
	
	    ThreadCheckMsgLive++  
		
        current_time := time.Now()
		pass_time    := current_time.Sub( start_time )
		
		if pass_time > time_out_msec {
		    break
		}
		
        pi := zmq.PollItems{
    	     zmq.PollItem{ Socket: PortInAsciiSUB, Events: zmq.POLLIN},
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
		                   
		        str := string(buf)
		        ad.Println( "IN ASCII : [%s]\n", str )
			}	
		}
		
	}
	
	ad.Println( "SleepCheckMsg() end" )

}

func checkSameMessage( b byte) (bool) {
//    ad.Println( "check byte %d" , b )

    if len( checkMsg ) == 0 {
	    return false
	}

	if len( checkMsg ) <= len( currentMsg ) {
		diff := len(currentMsg) - len(checkMsg) + 1
	    currentMsg = currentMsg[diff:]
	}

	currentMsg = append( currentMsg, b )
	fmt.Printf( "Check(%v): [%v]\n", len(checkMsg), string(checkMsg ))
	fmt.Printf( "Read (%v) : [%v]\n", len(currentMsg), string(currentMsg ))

	if bytes.Equal(checkMsg, currentMsg) {

		fmt.Printf("Checked: %v\n", string(checkMsg))
		ad.Println( "check ok [%s:%s]" , string(checkMsg), string(currentMsg) )
	
		return true
	}
	
	return false
}

func ThreadCheckMsg() {
    
	ad.Println( "ThreadCheckMsg() start" )
	
	ThreadCheckMsgReqEnd = false
	ThreadCheckMsgRun    = true 
	cmdIndex             = 0
	
	checkMsg    = []byte{}
	currentMsg  = []byte{}
	concat_buf := []byte{}
	remainBuffer := []byte{}
	
	for !ThreadCheckMsgReqEnd {
	
	    if cmdIndex >= len( commands ) {
		    ThreadCheckMsgReqEnd = true
			break
		}
		curCommand = commands[cmdIndex]
		
		
		switch curCommand.Cmd {
		
		case "doc"      : ad.Println( "doc command [%s]", curCommand.Value )
		                  ar.WriteDocument( curCommand.Value )
						  cmdIndex++
						  continue
						  
		case "time"     : ad.Println( "time command [%s]", curCommand.Value )
		                  timeOut,_ = strconv.Atoi( curCommand.Value )
						  cmdIndex++
						  continue

		case "sleep"    : ad.Println( "sleep command [%s]", curCommand.Value ) 
		                  sleepTime ,_ := strconv.Atoi( curCommand.Value )
		                  SleepCheckMsg( sleepTime ) 
						  cmdIndex++
						  continue
		
		case "send"     : ad.Println( "send command [%s]", curCommand.Value )
		                  if err := PortOutAsciiPUB.Send([]byte(curCommand.Value), 0); err != nil {
                              ad.Println( "fail do not send to OUT ASCII [%s]", err )
	                          reason := fmt.Sprintf( "do not send to OUT ASCII [%s]", err )
	                          ar.SetResultError( reason )
 		                  	  break
    	                  }
	                      ad.Println( "send command OK" ) 
		                  cmdIndex++ 
						  continue
		
		case "check"    : ad.Println( "check command [%s]", curCommand.Value ) 
		                  checkMsg    = []byte(curCommand.Value)
		
		default         : ad.Println( "fail unknow command [%s]", curCommand.Cmd )
	                      reason := fmt.Sprintf( "unknow command [%s]", curCommand.Cmd )
	                      ar.SetResultError( reason )
						  break
		}


		found := false
		for index, b := range remainBuffer {
			if checkSameMessage( b ) {
				currentMsg  = []byte{}
				cmdIndex++

				if index < len(remainBuffer) {
					remainBuffer = remainBuffer[(index + 1):]
				}
				fmt.Printf("Remain buffer: [%s]\n", string(remainBuffer) )
				found = true
				break
			}
		}

		if (found) {
			continue
		}


		ThreadCheckMsgLive++

        pi := zmq.PollItems{
    	     zmq.PollItem{ Socket: PortInAsciiSUB, Events: zmq.POLLIN },
    	}

        event_count, err := zmq.Poll( pi, time.Duration(timeOut) * time.Millisecond )
   		if err != nil {
            ad.Println( "fail do not poll[%s]", err )
	        reason := fmt.Sprintf( "do not poll[%s]", err )
	        ar.SetResultError( reason )
 			break
			   
    	}
			
		if event_count == 0 {
		    cmdIndex++ 
            ad.Println( "check time out[%s]", err )
//	        reason := fmt.Sprintf( "do not read [%s]", err )
//	        ar.SetResultError( reason )
// 		    break;
			
		//	   // 시간 초과가 발생 했다.
		//        count++
		//        if count > 10000 {
		//	        jmcsvr.SendTestEndResult( `{ "result" : "ok" }` )
		//            ReadyTestState      = true
		//	        break
		//        }
			   
		} else {
            if pi[0].REvents&zmq.POLLIN != 0 {
	            buf, err := pi[0].Socket.Recv(0)
   	            if err != nil {
                     ad.Println( "fail do not read [%s]", err )
	                 reason := fmt.Sprintf( "do not read [%s]", err )
	                 ar.SetResultError( reason )
 		    	     break;
                }

				concat_buf = append(remainBuffer, buf...)
				remainBuffer = []byte{}

		        str := string(buf)
		        ad.Println( "IN ASCII : [%s]\n", str )
			}	
		}

		for index, b := range concat_buf {
			if checkSameMessage( b ) {
				currentMsg  = []byte{}
				cmdIndex++

				if index < len(concat_buf) {
					remainBuffer = concat_buf[(index + 1):]
				}
				fmt.Printf("Remain buffer: [%s]\n", string(remainBuffer) )
				break
			}
		}

	}

	ad.Println( "ThreadCheckMsg() End" )
	ThreadCheckMsgRun    = false

}

