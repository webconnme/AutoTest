package main

import (
_ 	"github.com/codeskyblue/go-sh"
_	"time"
_	"encoding/json"
	
	"jmclog"
_ 	"jmcsvr"

)

//---------------------------------------------------------------------------------------------------------------------
//   
//   컴포넌트 처리를 중지한다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func runScriptCmdStop( data interface{} ) bool {

    jmclog.LogWrite( "Stop Command\n" );
	 
	// 수행 조건을 얻는다. 
	m, ok := data.(map[string]interface{})
    if !ok {
	    jmclog.LogWrite( "fail Stop Command Syntax Error\n" );
        return false
    }

	component_id              := m["id"].(string)
	 
	jmclog.LogWrite( "id     = [%s]\n", component_id );
	
	start_json := `{ "cmd" : "stop", "id" : "` + component_id  + `" }`
	sendCtrlCmdToComponent( start_json )

    return true
}
