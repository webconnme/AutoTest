package main

import (
_   "fmt"
    "os"
_    "time"
	"encoding/json"
	zmq "github.com/alecthomas/gozmq"
	
	    "jmclog"
	js  "jmcsvr"
)

// Ctrl SUB Cmd type
const (
    CMD_KILL    = "kill"
    CMD_INIT    = "init"
    CMD_SET     = "set"
    CMD_LINK    = "link"
    CMD_UNLINK  = "unlink"
    CMD_START   = "start"
    CMD_STOP    = "stop"
)

type CtrlCmdMsg struct {
	Cmd     string      `json:"cmd"`     // 시험 진행 명령
	Id      string      `json:"id"`      // 시험 진행 ID 
	Data    interface{} `json:"data"`    // 시험 설정 데이터 
}

var zmqContext   *zmq.Context
var zmqCtrlSUB   *zmq.Socket             // 제어용 SUB 소켓 

var myComponentId string                 // 내 자신의 ID 

var PortInAsciiSUB   *zmq.Socket = nil   // 분석 문자열 수신 SUB 소켓 


//---------------------------------------------------------------------------------------------------------------------
//   
//   JUE 를 초기화 한다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func initData() {
    jmclog.LogWrite( "Call function initData()\n" )
}

//---------------------------------------------------------------------------------------------------------------------
//   
//   데이터를 설정한다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func setData( data interface{} ) {
    jmclog.LogWrite( "Call function setData()\n" )

	// 수행 조건을 얻는다. 
	m, ok := data.(map[string]interface{})
    if !ok {
	    jmclog.LogWrite( "fail Set Command Syntax Error\n" );
        return 
    }

	for k,v := range m {
	    jmclog.LogWrite( " -- [%s] [%s]\n", k, v )
	}
}

//---------------------------------------------------------------------------------------------------------------------
//   
//   LINK 를 처리한다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func linkChannel( data interface{} ) {
    jmclog.LogWrite( "Call function linkChannel()\n" )

	// 수행 조건을 얻는다. 
	m, ok := data.(map[string]interface{})
    if !ok {
	    jmclog.LogWrite( "fail Link Command Syntax Error\n" );
        return 
    }

//	for k,v := range m {
//	    jmclog.LogWrite( " -- [%s] [%s]\n", k, v )
//	}


    link_port    := m["port"].(string)
	link_channel := m["channel"].(string)
	
	if link_port ==  "IN ASCII" {
	
	    jmclog.LogWrite( "Connect PORT [IN ASCII]\n" )
        jmclog.LogWrite( " -- link_channel = [%s]\n", link_channel )
	    
        //	만약 이전에 소켓이 열려 있다면 닫는다. 
	    if PortInAsciiSUB != nil {
	        PortInAsciiSUB.Close()
	    }
        
	    // 제어 수신용 zmq SUB 소켓을 만든다. 
	    var err error
	    
	    PortInAsciiSUB, err = zmqContext.NewSocket(zmq.SUB)
	    if err != nil {
	        jmclog.LogWrite( "fail New Socket PortInAsciiSUB  Error [%s]\n", err );
	    	return
	    }
	    
	    // 채널을 연결한다.
	    PortInAsciiSUB.Connect("ipc:///tmp/jmc_channel" + link_channel )
		PortInAsciiSUB.SetSubscribe("")

	}
	
	return 
}

//---------------------------------------------------------------------------------------------------------------------
//   
//   UNLINK 를 처리한다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func unlinkChannel( data interface{} ) {
    jmclog.LogWrite( "Call function unlinkChannel()\n" )

	// 수행 조건을 얻는다. 
	m, ok := data.(map[string]interface{})
    if !ok {
	    jmclog.LogWrite( "fail Unlink Command Syntax Error\n" );
        return 
    }

//	for k,v := range m {
//	    jmclog.LogWrite( " -- [%s] [%s]\n", k, v )
//	}


    link_port    := m["port"].(string)
	link_channel := m["channel"].(string)
	
	if link_port ==  "IN ASCII" {
	    jmclog.LogWrite( "Connect PORT [IN ASCII]\n" )
        jmclog.LogWrite( " -- unlink_channel = [%s]\n", link_channel )
	    
        //	만약 이전에 소켓이 열려 있다면 닫는다. 
	    if PortInAsciiSUB != nil {
	        PortInAsciiSUB.Close()
			PortInAsciiSUB = nil
	    }
        
	}
	
	return 
}

//---------------------------------------------------------------------------------------------------------------------
//   
//   JUE 의 동작을 시작한다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func start() {
    jmclog.LogWrite( "Call function start()\n" )
	SetReadyCheckMessage( false )
}

//---------------------------------------------------------------------------------------------------------------------
//   
//   JUE 의 동작을 중지 한다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func stop() {
    jmclog.LogWrite( "Call function stop()\n" )
	SetReadyCheckMessage( true )
}

//---------------------------------------------------------------------------------------------------------------------
//   
//   이 프로그램은 RS232 인터페이스 JUE 이다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func main() {

	zmqContext, _ = zmq.NewContext()
    defer zmqContext.Close()
	
	// 로그 시작  
	jmclog.Start( jmclog.CLIENT,            // 클라이언트 모드 
	              "",                       // 로그 파일 이름은 지정하지 않는다. 
				  "ck_message",             // 로그 프롬프트
				  zmqContext )              // zmq 컨텍스트 
				  
    // 컴포넌트 관리 서버를 시작한다. 
	js.Start( js.COMPONENT,                 // 컴포넌트 모드 
				  zmqContext )              // zmq 컨텍스트 
	
    jmclog.LogWrite( "ck_message component start" );
	
	component_class     := "ck_message"
	component_option    := ""
	component_info_json := ""
	
	if len(os.Args) < 2 {
	    jmclog.LogWrite( "fail Have No Component Argument" )
		goto error_no_argument
	} 
	
	myComponentId     = os.Args[1]
	
	if len(os.Args) > 2 {
	   component_option = os.Args[2]
	} 
	
	jmclog.LogWrite( "id = [%s]",     myComponentId );
    jmclog.LogWrite( "option = [%s]", component_option );
	

	// 제어 수신용 zmq SUB 소켓을 만든다. 
	zmqCtrlSUB, _ = zmqContext.NewSocket(zmq.SUB)
	
	zmqCtrlSUB.Connect("ipc:///tmp/jmc_ctrl")
	zmqCtrlSUB.SetSubscribe("")
		
    // 컴포넌트 등록 요청을 한다. 
	component_info_json  = `{`
	component_info_json +=     `"id":"`    + myComponentId    + `",`
	component_info_json +=     `"class":"` + component_class + `"`
	component_info_json += `}`
	
	js.RegisterComponet( component_info_json )
	
	// 시험 진행 쓰레드를 수행한다. 
	
	go MainCheckMessage()
	
	// 제어 명령을 대기 하고 처리 한다. 
	jmclog.LogWrite( "[%s:%s] Wait\n", component_class, myComponentId );
    for {
	    buf, err := zmqCtrlSUB.Recv(0)
   	    if err != nil {
    		jmclog.LogWrite( "fail zmq %s recv: %s", component_class, err  );	
			break;
    	}
		
		str := string(buf)
		jmclog.LogWrite( str );
		
		var ctrl_cmd_msg  CtrlCmdMsg
								   
        err = json.Unmarshal( buf, &ctrl_cmd_msg )
	    if err != nil {
        
	        jmclog.LogWrite( "fail Convert JSON to Ctrl  [%s]\n", string(buf) )
	        jmclog.LogWrite( "%s\n", err )
        
	    }
		
		jmclog.LogWrite( "Ctrl Rx Cmd  = [%s]\n", ctrl_cmd_msg.Cmd )
		jmclog.LogWrite( "Ctrl Rx Id   = [%s]\n", ctrl_cmd_msg.Id  ) 

		if ctrl_cmd_msg.Id != myComponentId {
		   continue
		}
		
		if ctrl_cmd_msg.Data  != nil { 
		    jmclog.LogWrite( "Ctrl Rx Data = [%s]\n", ctrl_cmd_msg.Data  )
		}
		
		
		switch ctrl_cmd_msg.Cmd {
        case CMD_KILL :  js.UnregisterComponet( myComponentId )
		                 goto error_no_argument
						 
        case CMD_STOP  :  stop()
        case CMD_START :  start()
		
        case CMD_INIT  :  initData()
        case CMD_SET   :  setData( ctrl_cmd_msg.Data )
		
		case CMD_LINK  :  linkChannel( ctrl_cmd_msg.Data )
		case CMD_UNLINK:  unlinkChannel( ctrl_cmd_msg.Data )
		
		}
		
	}
	
error_no_argument:

	if PortInAsciiSUB != nil {
	    PortInAsciiSUB.Close()
	}
	
    // 시험 진행 쓰레드를 중지 한다.  
     StopCheckMessage()
	
	jmclog.LogWrite( "End [%s:%s]\n", component_class, myComponentId );
	// zmq 의 대기열 처리를 위해서 1 초간 잠든다. 
//	time.Sleep(1000 * time.Millisecond)
	
	zmqCtrlSUB.Close()
	js.End()
	
	jmclog.End()

}

