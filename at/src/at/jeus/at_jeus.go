package jeus

import (
 	zmq "github.com/alecthomas/gozmq"
)

import (
    "at"
)

const JS_DEFAULT_NAME    = "jus"

// JS Frame Command 
const (
    JS_CMD_REGISTER_JEU              = "register_jeu"
    JS_CMD_UNREGISTER_JEU            = "unregister_end"
    JS_CMD_CHECK_JEU                 = "check"
    JS_CMD_SET_READY_JEU             = "set_ready_jeu"
    JS_CMD_SET_STARTING_JEU          = "set_starting_jeu"
    JS_CMD_SET_RUN_JEU               = "set_run_jeu"
    JS_CMD_SET_STOPPING_JEU          = "set_stopping_jeu"
    JS_CMD_SET_ENDING_JEU            = "set_ending_jeu"
	JS_CMD_CHECK_STATE_STARTING_JEU  = "check_state_starting_jeu"
)

// JS Frame Command JSON

// CALL
type AtjJsonRegisterJEU struct {
    Cmd     string      `json:"cmd"`     // 명령 
	Id      string      `json:"id"`      // 등록 ID
}

type AtjJsonUnRegisterJEU struct {
    Cmd     string      `json:"cmd"`     // 명령 
	Id      string      `json:"id"`      // 등록 해제 ID
}

type AtjJsonCheckJEU struct {
    Cmd     string      `json:"cmd"`     // 명령 
	Id      string      `json:"id"`      // 체크 ID
}

type AtjJEUInfo struct {
	Id        string                     // JEU ID    
	Class     string                     // 현재 사용하지 않음
    State     int                        //	JEU 실행 상태
}

// ACK 
type AtjJsonDone struct {
	Result   string      `json:"result"`   // 처리 결과  "ok", "fail", "been"
}

type AtJEUServer struct {

	ZmqContext     *zmq.Context
	
	ServerMode     bool                         // 서버에서 동작하는가 클라이언트가 동작하는가?
	
	AF             *at.AtFrame

    JEUs           map[string]*AtjJEUInfo
	
}

func (js *AtJEUServer) Close() {
   
   if js.ServerMode {
       js.CloseServer()
   } else {
       js.CloseClient()
   }
}

