package main

import (
_	"log"
_   "time" 
_	"os"
_	"encoding/json"
   "jmclog"
)

import (
	"github.com/mikepb/go-serial"
)

//	zmq "github.com/alecthomas/gozmq"


// var rs232RxDone          chan bool                                // 처리 종료 대기 싱크용 변수 
// var rs232TxDone          chan bool                                // 처리 종료 대기 싱크용 변수 

var DefaultOptions = serial.Options{
    BitRate:     115200,
    DataBits:    8,
    Parity:      serial.PARITY_NONE,
    StopBits:    1,
    FlowControl: serial.FLOWCONTROL_NONE,
}

// var zmqContext           *zmq.Context                             // zmq 컨텍스트 
var serialPort           *serial.Port                                // 시리얼 포트 
var ReqStopRxRS232        bool                                       // 수신 쓰레드 종료 요구  
		
//---------------------------------------------------------------------------------------------------------------------
//   
//  RS232 포트를 연다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func OpenRS232( device_path string ) {

//	zmqContext  = zmq_context
	

	options := DefaultOptions

	var err error
	
    serialPort, err = options.Open( device_path )
    if err != nil {
       
    }
	
}

//---------------------------------------------------------------------------------------------------------------------
//   
//  RS232 포트를 닫는다.
//   
//---------------------------------------------------------------------------------------------------------------------
func CloseRS232() {

	serialPort.Close()

}

//---------------------------------------------------------------------------------------------------------------------
//   
//  RS232 수신 쓰레드 
//   
//---------------------------------------------------------------------------------------------------------------------
func MainRxRS232() {

   jmclog.LogWrite( "start MainRxRS232\n" ); 
   ReqStopRxRS232 = false
   for ReqStopRxRS232 != true {
//        jmclog.LogWrite( "CALL serialPort.InputWaiting()\n" ); 
 		remain, err := serialPort.InputWaiting()
 		if err != nil {
 		    jmclog.LogWrite( "fail Serial Port Read Wait[%s]\n", err ); 
 			break;
 		}
 
        // 수신된 데이터가 있는가?
 		if remain != 0 {
		
 		    jmclog.LogWrite( "remain = %d\n", remain ); 
 		    buf := make([]byte, remain)
 		    length, err := serialPort.Read(buf)
 		    if err != nil {
 		    	jmclog.LogWrite( "fail Serial Port Read [%s]\n", err ); 
 		    	break;
 		    }
 			jmclog.LogWrite( "read length [%d]\n", length ); 
			
			// 수신된 데이터를 발송 한다. 
			err = PortRxPUB.Send(buf[:length], 0)
            if err != nil {
            	jmclog.LogWrite( "zmq channel send: ", err  );	
            }			

 		}	
 		
 //		rx.Send(buf[:length], 0)
    
   }
   
   jmclog.LogWrite( "End MainRxRS232()\n" ); 
   
//   mainRxDone<- ture;
   
}

func StopRxRS232() {
    ReqStopRxRS232 = true
}

