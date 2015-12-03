package jeus

import (
_    "fmt" 
    "at"
	"time"
)

func NewAtJEUClient( af *at.AtFrame ) (*AtJEUServer, error) {

    js := &AtJEUServer{}
	
	js.AF         = af
	js.ZmqContext = af.ZmqContext
	
	js.ServerMode = false
	
	return js, nil
}

func (js *AtJEUServer) CloseClient() {

}

func (js *AtJEUServer) Reset() {
    js.AF.SendCommandReset( JS_DEFAULT_NAME )
	time.Sleep( time.Millisecond)
}

func (js *AtJEUServer) RegisterJEU( id string ) (error) {

	js_data  := AtjJsonRegisterJEU{ Cmd : JS_CMD_REGISTER_JEU, Id : id }
	
    result, err := js.AF.SendCall( JS_DEFAULT_NAME , js_data, 1000 )
    _ = result
	if err != nil {
	    return err
	}
	
    return nil
}

func (js *AtJEUServer) UnRegisterJEU( id string ) (error) {

	js_data  := AtjJsonUnRegisterJEU{ Cmd : JS_CMD_UNREGISTER_JEU, Id : id }
	
    result, err := js.AF.SendCall( JS_DEFAULT_NAME , js_data, 1000 )
    _ = result
	if err != nil {
	    return err
	}

    return nil
}

func (js *AtJEUServer) CheckJEU( id string ) ( bool, error) {

	js_data  := AtjJsonUnRegisterJEU{ Cmd : JS_CMD_CHECK_JEU, Id : id }
	
    result, err := js.AF.SendCall( JS_DEFAULT_NAME , js_data, 1000 )
	if err != nil {
	    return false, err
	}
	
    // 맵 데이터로 변환한다. 
    m, ok := result.(map[string]interface{})
      if !ok {
        return false, nil
    }
	
	done, ok := m["result"]
	if !ok {
        return false, nil
	}
	
	if done != "been" {
	    return false, nil
	}
	
    return true, nil
}

func (js *AtJEUServer) SetJEUStateReady( id string ) (error) {

	js_data  := AtjJsonRegisterJEU{ Cmd : JS_CMD_SET_READY_JEU, Id : id }
	
    result, err := js.AF.SendCall( JS_DEFAULT_NAME , js_data, 1000 )
    _ = result
	if err != nil {
	    return err
	}
	
    return nil
}

func (js *AtJEUServer) SetJEUStateStarting( id string ) (error) {

	js_data  := AtjJsonRegisterJEU{ Cmd : JS_CMD_SET_STARTING_JEU, Id : id }
	
    result, err := js.AF.SendCall( JS_DEFAULT_NAME , js_data, 1000 )
    _ = result
	if err != nil {
	    return err
	}
	
    return nil
}

func (js *AtJEUServer) SetJEUStateRun( id string ) (error) {

	js_data  := AtjJsonRegisterJEU{ Cmd : JS_CMD_SET_RUN_JEU, Id : id }
	
    result, err := js.AF.SendCall( JS_DEFAULT_NAME , js_data, 1000 )
    _ = result
	if err != nil {
	    return err
	}
	
    return nil
}

func (js *AtJEUServer) SetJEUStateStopping( id string ) (error) {

	js_data  := AtjJsonRegisterJEU{ Cmd : JS_CMD_SET_STOPPING_JEU, Id : id }
	
    result, err := js.AF.SendCall( JS_DEFAULT_NAME , js_data, 1000 )
    _ = result
	if err != nil {
	    return err
	}
	
    return nil
}

func (js *AtJEUServer) SetJEUStateEnding( id string ) (error) {

	js_data  := AtjJsonRegisterJEU{ Cmd : JS_CMD_SET_ENDING_JEU, Id : id }
	
    result, err := js.AF.SendCall( JS_DEFAULT_NAME , js_data, 1000 )
    _ = result
	if err != nil {
	    return err
	}
	
    return nil
}

func (js *AtJEUServer) CheckJEUStateStarting() ( bool, error) {

	js_data  := AtjJsonUnRegisterJEU{ Cmd : JS_CMD_CHECK_STATE_STARTING_JEU, Id : "" }
	
    result, err := js.AF.SendCall( JS_DEFAULT_NAME , js_data, 1000 )
	if err != nil {
	    return false, err
	}
	
    // 맵 데이터로 변환한다. 
    m, ok := result.(map[string]interface{})
      if !ok {
        return false, nil
    }
	
	done, ok := m["result"]
	if !ok {
        return false, nil
	}
	
	if done != "been" {
	    return false, nil
	}
	
    return true, nil
}


/*
// RUN STATE
const (
    AF_JEU_INIT       = 0        // 프로그램 시작 중 - RegisterJEU() 함수 호출에 의해 설정
    AF_JEU_READY      = 1        // 명령 대기 중     - SetJEUStateReady() 함수 호출에 의해 설정
    AF_JEU_STARTING   = 2        // 시작 중          - SetJEUStateStarting() 함수 호출에 의해 설정 시작은 START 명령에 의해서 시작
    AF_JEU_RUN        = 3        // 동작 중          - SetJEUStateRun() 함수 호출에 의해 설정 START 명령에 따른 동작이 모두 시작되었을때 설정
    AF_JEU_STOPPING   = 4        // 중지 중  여기서 AF_STATE_READY 로 이동 
	                             //                  - SetJEUStateStopping() 함수 호출에 의해 설정  STOP 명령에 따른 동작이 모두 종료 되었을때 설정 
    AF_JEU_ENDING     = 5        // 종료 중          - SetJEUStateEnding() 함수 호출에 의해 설정    KILL 명령에 따라서 UnRegisterJEU() 함수가 호출될때까지 유지
    AF_JEU_END        = 6        // 종료             - UnRegisterJEU() 함수 호출로 이 상태 발생 
	                             //                    jes 에서 이상태는 관리하지 않으므로 실제로는 존재하지 않는 상태  
)

*/

