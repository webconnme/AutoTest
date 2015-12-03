package jmclog

import (
    "fmt"
_	"time"
	"log"
	"os"
	"io"
	zmq "github.com/alecthomas/gozmq"
)

// server_mode type
const (
    SERVER  = true
	CLIENT  = false
)

// LogCmd type
const (
    STOP    = "<STOP>"
	MESSAGE = "<0000>"
)

var serverMode           bool                                        // 로거가 서버로 동작하는가를 표시
var logFilename          string                                      // 로거 파일 이름
var logPrompt            string                                      // 로거 프롬프트

var	logger               *log.Logger                                 // 로거  

var zmqContext           *zmq.Context                                // zmq 컨텍스트 
var zmqLogSUB            *zmq.Socket                                 // 로그 수신 소켓 
var zmqLogPUB            *zmq.Socket                                 // 로그 송신 소켓 

var logDone              chan bool                                   // 처리 종료 대기 싱크용 변수 

var logger_stop_req      bool                                        // 로거 중지 요청 

//---------------------------------------------------------------------------------------------------------------------
//   
//   로그 쓰래드 메인
//   
//---------------------------------------------------------------------------------------------------------------------
func logServer( log_done chan bool ) {

    // 로그를 기록할 파일을 오픈한다. 
    logFile, err := os.OpenFile( logFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalln("Failed to open log file", logFilename, ":", err)
    } 	
 	defer logFile.Close()
	
	// 로그 파일과 표준 출력을 묶는다.
 	multiLog := io.MultiWriter( logFile, os.Stdout)
 	
    logger = log.New( multiLog, ">> ", log.Ldate|log.Ltime )
 	
	logger_stop_req = false
	
	// 로그 수신 로컬 호스트 네트워크를 연다. 
	
	zmqLogSUB, err = zmqContext.NewSocket(zmq.SUB)
	defer zmqLogSUB.Close()
	
	zmqLogSUB.SetSubscribe("")
	zmqLogSUB.Bind("ipc:///tmp/jmc_log")
	
	log_done <- true
	
	// 로그 수신을 한다. 
	
    for {
	    buf, err := zmqLogSUB.Recv(0)
   	    if err != nil {
    		logger.Printf( "fail zmq log recv: %s", err  );	
			break;
    	}
		str := string(buf)
		cmd := str[:6]
		if cmd == STOP {
		    break
		}
		msg := str[6:] 
	    logger.Printf( msg );
	}
	
	logger.Printf( "log thread end\n" );	
	
	log_done <- true
	
}


//---------------------------------------------------------------------------------------------------------------------
//   
//  로그 서비스를 시작한다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func Start( mode bool, filename string, prompt string, zmq_context *zmq.Context ) {

    serverMode  = mode
	logFilename = filename
	logPrompt   = prompt
	zmqContext  = zmq_context
	
	if serverMode {
	
        logDone = make(chan bool)
        go logServer( logDone )
	    <-logDone
		
	} 
	
	zmqLogPUB, _ = zmqContext.NewSocket(zmq.PUB)
	zmqLogPUB.Connect("ipc:///tmp/jmc_log")
		
}
 

//---------------------------------------------------------------------------------------------------------------------
//   
//  로그를 중지 한다.  
//   
//---------------------------------------------------------------------------------------------------------------------
func End() {

	if serverMode {
	
        msg := STOP
	    
        err := zmqLogPUB.Send([]byte(msg), 0)
        if err != nil {
        	logger.Printf( "zmq send: ", err  );	
        }
		<-logDone
	
	} 
	
    zmqLogPUB.Close()

}

//---------------------------------------------------------------------------------------------------------------------
//   
//  로그를 쓴다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func LogWrite( format string, args ...interface{} ) {

    msg := MESSAGE + logPrompt + " : " + fmt.Sprintf( format, args... )
	
    err := zmqLogPUB.Send([]byte(msg), 0)
    if err != nil {
    	logger.Printf( "zmq send: ", err  );	
    }

}

