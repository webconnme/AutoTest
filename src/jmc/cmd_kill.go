package main

import (
_ 	"github.com/codeskyblue/go-sh"
	"time"
_	"encoding/json"
	
	"jmclog"
js	"jmcsvr"

)

//---------------------------------------------------------------------------------------------------------------------
//   
//   생성되었던 컴포넌트를 파괴한다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func runScriptCmdKill( data interface{} ) bool {
    jmclog.LogWrite( "Kill Command\n" );
	 
	// 수행 조건을 얻는다. 
	m, ok := data.(map[string]interface{})
    if !ok {
	    jmclog.LogWrite( "fail KILL Command Syntax Error\n" );
        return false
    }

	component_id              := m["id"].(string)
	 
	jmclog.LogWrite( "id     = [%s]\n", component_id );
	
	kill_json := `{ "cmd" : "kill", "id" : "` + component_id  + `" }`
	sendCtrlCmdToComponent( kill_json )

    // 컴포넌트가 해제되었음을 알리는 메세지를 기다린다. 	
	jmclog.LogWrite( "wait unregister component [%s]", component_id ) 
	var loop_out int
	for loop_out = 10; loop_out > 0 ; loop_out-- {
	    if !js.CheckBeenComponet( component_id ) {
		    break;
		}
        time.Sleep(10 * time.Millisecond)
	}
	if loop_out == 0 {
	    jmclog.LogWrite( "fail kill component program [%s] No Unregister Component", component_id ) 
        RunResult.success  = false;     // 실행 실패
		return false
	}
	
    jmclog.LogWrite( "Sucess Kill Component Program [%s]", component_id ) 		
	 
    return true
}
