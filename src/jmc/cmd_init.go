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
//   컴포넌트의 설정된 데이터를 초기화 한다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func runScriptCmdInit( data interface{} ) bool {
    jmclog.LogWrite( "Init Command\n" );
	 
	// 수행 조건을 얻는다. 
	m, ok := data.(map[string]interface{})
    if !ok {
	    jmclog.LogWrite( "fail Init Command Syntax Error\n" );
        return false
    }

	component_id              := m["id"].(string)
	 
	jmclog.LogWrite( "id     = [%s]\n", component_id );
	
	init_json := `{ "cmd" : "init", "id" : "` + component_id  + `" }`
	sendCtrlCmdToComponent( init_json )

    return true
}
