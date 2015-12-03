package jmcsvr

import (
_    "fmt"
_	"time"
    "encoding/json"
	zmq "github.com/alecthomas/gozmq"
	
	"jmclog"
)

// server_mode type
const (
    SERVER     = true
	COMPONENT  = false
)

const (
    STOP                      = "<STOP>"           // 서버를 중지 시킨다.
	ECHO                      = "<ECHO>"           // 에코 테스트를 한다. 
	CHECK_BEEN_COMPONENT      = "<CKBC>"           // 컴포넌트 존재를 확인한다.
	REGISTER_COMPONENT        = "<RGCM>"           // 컴포넌트를 등록 한다.
	UNREGISTER_COMPONENT      = "<URCM>"           // 컴포넌트를 등록 해제 한다.
	GET_ALL_COMPONENT_ID_LIST = "<GCIL>"           // 모든 컴포넌트 리스트를 얻는다.
	
	SET_WAIT_TEST_END         = "<SWTE>"           // 시험이 끝나기를 기다리도록 설정한다. 
	CHECK_WAIT_TEST_END       = "<CWTE>"           // 시험이 끝났는가를 확인한다. 
	SEND_TEST_END_RESULT      = "<STER>"           // 시험 결과를 전송한다. 
	
)

type SvrComponentInfo struct {
	Id        string       `json:"id"`           // 컴포넌트 구별 인자 
	Class     string       `json:"class"`        // 컴포넌트 종류
}

var serverMode           bool                                        // 로거가 서버로 동작하는가를 표시

var zmqContext           *zmq.Context                                // zmq 컨텍스트 
var zmqSvrREP            *zmq.Socket                                 // 서버 수신 소켓 
var zmqSvrREQ            *zmq.Socket                                 // 서버 송신 소켓 

var svrDone              chan bool                                   // 처리 종료 대기 싱크용 변수 

var componentInfos       map[string]SvrComponentInfo = make(map[string]SvrComponentInfo)                // 컴포넌트 정보  

var svrCheckingState      bool                                        // 체크 중임을 나타낸다. 

//---------------------------------------------------------------------------------------------------------------------
//   
//   서버 쓰래드 메인
//   
//---------------------------------------------------------------------------------------------------------------------
func svrServer( svr_done chan bool ) {

	// 서버 수신 로컬 호스트 네트워크를 연다. 

	zmqSvrREP, _ = zmqContext.NewSocket(zmq.REP)
	defer zmqSvrREP.Close()
	
	zmqSvrREP.Bind("ipc:///tmp/jmc_svr")
	
	svr_done <- true
	
	// 처리 요청을 수신을 한다. 
	
    for {
	    buf, err := zmqSvrREP.Recv(0)
   	    if err != nil {
    		jmclog.LogWrite( "fail zmq svr recv: %s", err  );	
			break;
    	}
		str := string(buf)
		cmd := str[:6]
		
		msg := str[6:] 
//		fmt.Println( msg );
	    
		ack_msg := "OK"
		
		switch cmd {
		
		case STOP                 : jmclog.LogWrite( "zmq recv: STOP"  );	
		
		case ECHO                 : jmclog.LogWrite( "zmq recv: ECHO"  );	

					 
		case REGISTER_COMPONENT   : jmclog.LogWrite( "zmq recv: REGISTER_COMPONENT"    );
		                            jmclog.LogWrite( "-- component info %s", msg       );
									
	                               // 내부 자료 구조체로 변경한다.  
								   var component_info SvrComponentInfo
								   b := []byte( msg )
								   
                                   err := json.Unmarshal( b, &component_info )
	                               if err != nil {
                                   
	                                   jmclog.LogWrite( "fail Convert Component info JSON [%s]\n", msg )
	                                   jmclog.LogWrite( "%s\n", err )
									   ack_msg = "FAIL"
                                   
	                               }	
								   componentInfos[component_info.Id] = component_info	
//								   jmclog.LogWrite( "%s\n", component_info );
						
		case UNREGISTER_COMPONENT : jmclog.LogWrite( "zmq recv: UNREGISTER_COMPONENT"  );	
		                            jmclog.LogWrite( "-- component id [%s]", msg       ); 
									
									delete( componentInfos , msg )
		
		case CHECK_BEEN_COMPONENT : jmclog.LogWrite( "zmq recv: CHECK_BEEN_COMPONENT"  );	
		                            jmclog.LogWrite( "-- component id [%s]", msg       );
									
									_, ok := componentInfos[msg]
									
									if !ok {
									    ack_msg = "FAIL"
									}
									
		case GET_ALL_COMPONENT_ID_LIST : jmclog.LogWrite( "zmq recv: GET_ALL_COMPONENT_ID_LIST"  );	
		                            var id_list []string = make( []string, 0 )
									
									for key, _ := range componentInfos {
									    id_list = append( id_list, key )
									}

                                    id_list_json, err := json.Marshal( id_list )
	                                if err != nil {
                                    
                                        jmclog.LogWrite( "fail Convert Compornent List To JSON" )
	                                    ack_msg = "[]"
										
	                                } else {
									    ack_msg = string(id_list_json)
                                    }									
									
									jmclog.LogWrite( "id_list_json = %s" , ack_msg )
									
        case SET_WAIT_TEST_END : jmclog.LogWrite( "zmq recv: SET_WAIT_TEST_END"  );	
		
					             svrCheckingState = true				
								 
		case CHECK_WAIT_TEST_END :  // jmclog.LogWrite( "zmq recv: CHECK_WAIT_TEST_END"  );	
		
									if svrCheckingState {
									    ack_msg = "NO_END"
									}
									
		case SEND_TEST_END_RESULT  : jmclog.LogWrite( "zmq recv: SEND_TEST_END_RESULT"    );
		                            jmclog.LogWrite( "-- result  %s", msg       );
									
	                               // 내부 자료 구조체로 변경한다.  
//								   var component_info SvrComponentInfo
//								   b := []byte( msg )
//								   
//                                   err := json.Unmarshal( b, &component_info )
//	                               if err != nil {
//                                   
//	                                   jmclog.LogWrite( "fail Convert Component info JSON [%s]\n", msg )
//	                                   jmclog.LogWrite( "%s\n", err )
//									   ack_msg = "FAIL"
//                                   
//	                               }	
//								   componentInfos[component_info.Id] = component_info	
//								   jmclog.LogWrite( "%s\n", component_info );
								   
								   svrCheckingState = false
		                    
		}
	    
        err = zmqSvrREP.Send([]byte(ack_msg), 0)
        if err != nil {
        	jmclog.LogWrite( "zmq send: ", err  );	
        }
		
		if cmd == STOP {
		    goto end_svr_server
		}
		
	}
	
end_svr_server:
	
	jmclog.LogWrite( "svr thread end\n" );	
	
	svr_done <- true
	
}


//---------------------------------------------------------------------------------------------------------------------
//   
//  서버 서비스를 시작한다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func Start( mode bool, zmq_context *zmq.Context ) {

    serverMode  = mode
	zmqContext  = zmq_context
	
	if serverMode {
	
        svrDone = make(chan bool)
        go svrServer( svrDone )
	    <-svrDone
		
	} 
	
	zmqSvrREQ, _ = zmqContext.NewSocket(zmq.REQ)
	zmqSvrREQ.Connect("ipc:///tmp/jmc_svr")
		
}
 

//---------------------------------------------------------------------------------------------------------------------
//   
//  서버 서비스를 중지 한다.  
//   
//---------------------------------------------------------------------------------------------------------------------
func End() {

	if serverMode {
	
        msg := STOP
	    
        err := zmqSvrREQ.Send([]byte(msg), 0)
        if err != nil {
        	jmclog.LogWrite( "zmq send: ", err  );	
        }
		
	    zmqSvrREQ.Recv(0)
		
		<-svrDone
	
	} 
	
    zmqSvrREQ.Close()

}

//---------------------------------------------------------------------------------------------------------------------
//   
//  서버의 동작 여부를 확인하기 위한 에코 
//   
//---------------------------------------------------------------------------------------------------------------------
func Echo() bool {

    msg := ECHO
	
    err := zmqSvrREQ.Send([]byte(msg), 0)
    if err != nil {
    	jmclog.LogWrite( "zmq send: ", err  );	
    }
	
	buf, rerr := zmqSvrREQ.Recv(0)
   	if rerr != nil {
    	jmclog.LogWrite( "fail zmq client recv: %s", rerr  );	
    }
	str := string(buf)
	if str != "OK" {
	   return false
	}

	return true
}

//---------------------------------------------------------------------------------------------------------------------
//   
//  특정 컴포넌트가 등록되었는가를 확인한다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func CheckBeenComponet( component_name string ) bool {

    msg := CHECK_BEEN_COMPONENT + component_name

    err := zmqSvrREQ.Send([]byte(msg), 0)
    if err != nil {
    	jmclog.LogWrite( "zmq send: ", err  );	
		return false
    }
	
	buf, rerr := zmqSvrREQ.Recv(0)
   	if rerr != nil {
    	jmclog.LogWrite( "fail zmq client recv: %s", rerr  );	
    }
	str := string(buf)
	if str != "OK" {
	   return false
	}
	
	return true

}

//---------------------------------------------------------------------------------------------------------------------
//   
//  컴포넌트를 등록한다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func RegisterComponet( component_info_json string ) bool {

    msg := REGISTER_COMPONENT + component_info_json

    err := zmqSvrREQ.Send([]byte(msg), 0)
    if err != nil {
    	jmclog.LogWrite( "zmq send: ", err  );	
    }
	
	buf, rerr := zmqSvrREQ.Recv(0)
   	if rerr != nil {
    	jmclog.LogWrite( "fail zmq client recv: %s", rerr  );	
    }
	str := string(buf)
	if str != "OK" {
	   return false
	}

	return true

}

//---------------------------------------------------------------------------------------------------------------------
//   
//  컴포넌트를 해제한다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func UnregisterComponet( component_id string ) bool {

    msg := UNREGISTER_COMPONENT + component_id

    err := zmqSvrREQ.Send([]byte(msg), 0)
    if err != nil {
    	jmclog.LogWrite( "zmq send: ", err  );	
    }
	
	buf, rerr := zmqSvrREQ.Recv(0)
   	if rerr != nil {
    	jmclog.LogWrite( "fail zmq client recv: %s", rerr  );	
    }
	str := string(buf)
	if str != "OK" {
	   return false
	}

	return true

}

//---------------------------------------------------------------------------------------------------------------------
//   
//  등록된 모든 컴포넌트 ID 리스트를 얻는다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func GetAllComponetIdList() ( ids []string, err error ) {

    msg := GET_ALL_COMPONENT_ID_LIST

    err = zmqSvrREQ.Send([]byte(msg), 0)
    if err != nil {
    	jmclog.LogWrite( "zmq svr send: ", err  );	
		return nil, err
		
    }
	
	buf, rerr := zmqSvrREQ.Recv(0)
   	if rerr != nil {
    	jmclog.LogWrite( "fail zmq client recv: %s", rerr  );	
    }

    var id_list []string = make( []string, 0 )
	err = json.Unmarshal( buf, &id_list )
	if err != nil {
    
	    jmclog.LogWrite( "fail Convert JSON To Component ID List [%s]\n", string(buf) )
	    return nil, err
	}	
	
	return id_list, nil

}

//---------------------------------------------------------------------------------------------------------------------
//   
//  시험이 끝나기를 기다리도록 설정한다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func SetWaitTestEnd() bool {

    msg := SET_WAIT_TEST_END

    err := zmqSvrREQ.Send([]byte(msg), 0)
    if err != nil {
    	jmclog.LogWrite( "zmq send: ", err  );	
    }
	
	buf, rerr := zmqSvrREQ.Recv(0)
   	if rerr != nil {
    	jmclog.LogWrite( "fail zmq client recv: %s", rerr  );	
    }
	str := string(buf)
	if str != "OK" {
	   return false
	}

	return true

}

//---------------------------------------------------------------------------------------------------------------------
//   
//  시험이 끝났는가를 체크 한다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func CheckWaitTestEnd() bool {

    msg := CHECK_WAIT_TEST_END

    err := zmqSvrREQ.Send([]byte(msg), 0)
    if err != nil {
    	jmclog.LogWrite( "zmq send: ", err  );	
    }
	
	buf, rerr := zmqSvrREQ.Recv(0)
   	if rerr != nil {
    	jmclog.LogWrite( "fail zmq client recv: %s", rerr  );	
    }
	str := string(buf)
	if str != "OK" {
	   return false
	}

	return true

}

//---------------------------------------------------------------------------------------------------------------------
//   
//  시험이 끝났음을 알린다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func SendTestEndResult( result_json string ) bool {

    msg := SEND_TEST_END_RESULT + result_json

    err := zmqSvrREQ.Send([]byte(msg), 0)
    if err != nil {
    	jmclog.LogWrite( "zmq send: ", err  );	
    }
	
	buf, rerr := zmqSvrREQ.Recv(0)
   	if rerr != nil {
    	jmclog.LogWrite( "fail zmq client recv: %s", rerr  );	
    }
	str := string(buf)
	if str != "OK" {
	   return false
	}

	return true
}

