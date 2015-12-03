package main

import (
_	"log"
   "time" 
_	"os"
_	"encoding/json"
)

import (
    zmq "github.com/alecthomas/gozmq"
)

import (
   "jmclog"
   "jmcsvr"
)


// var zmqContext           *zmq.Context                             // zmq 컨텍스트 
var ReqStopCheckMessage   bool                                       // 메세지 검사 쓰레드 종료 요구  
var ReadyTestState        bool                                       // 검사 대기 상태로 만든다. 

//---------------------------------------------------------------------------------------------------------------------
//   
//  RS232 수신 쓰레드 
//   
//---------------------------------------------------------------------------------------------------------------------
func MainCheckMessage() {

   jmclog.LogWrite( "start MainCheckMessage\n" ); 
   ReqStopCheckMessage = false
   ReadyTestState      = true
   
   count := 0;
   
   for ReqStopCheckMessage != true {
   
//        jmclog.LogWrite( "CALL MainCheckMessage Thread Loop\n" ); 

		if ReadyTestState {
		
		    count = 0;
			time.Sleep(1 * time.Millisecond)
			
		} else {

            pi := zmq.PollItems{
    		    zmq.PollItem{ Socket: PortInAsciiSUB, Events: zmq.POLLIN},
    	    }

            event_count, err := zmq.Poll( pi, 1 * time.Millisecond )
//            jmclog.LogWrite( "zmq Poll event_count = %d\n", event_count );
    		if err != nil {
    		   jmclog.LogWrite( "zmq Poll error =[%s]\n", err );
			   
			   jmcsvr.SendTestEndResult( `{ "result" : "fail" }` )
		       ReadyTestState      = true
			   
			   break
			   
    		}
			
			if event_count == 0 {
			
			   // 시간 초과가 발생 했다.
		        count++
		        if count > 10000 {
			        jmcsvr.SendTestEndResult( `{ "result" : "ok" }` )
		            ReadyTestState      = true
			        break
		        }
			   
			} else {
			   // 데이터가 수신되었다. 
			   
                  switch {
		          case pi[0].REvents&zmq.POLLIN != 0:
				  
	                       buf, err := pi[0].Socket.Recv(0)
   	                       if err != nil {
    	                   	jmclog.LogWrite( "fail zmq PortInAsciiSUB recv: %s", err  );	
		                   	break;
    	                   }
		                   
		                   str := string(buf)
		                   jmclog.LogWrite( "IN ASCII : [%s]\n", str );
					
		          }
			}

		}	
		
   }
   
//   mainRxDone<- ture;

}

//---------------------------------------------------------------------------------------------------------------------
//   
// 시험 진행 쓰레드를 중지 한다.  
//   
//---------------------------------------------------------------------------------------------------------------------
func StopCheckMessage() {
    ReqStopCheckMessage = true
}

//---------------------------------------------------------------------------------------------------------------------
//   
// 시험 진행 쓰레드를 중지 한다.  
//   
//---------------------------------------------------------------------------------------------------------------------
func SetReadyCheckMessage( enable bool ) {

    ReadyTestState = enable
}

