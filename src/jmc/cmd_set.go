package main

import (
_ 	"github.com/codeskyblue/go-sh"
_	"time"
	"encoding/json"
	
	"jmclog"
_ 	"jmcsvr"

)

//---------------------------------------------------------------------------------------------------------------------
//   
//   컴포넌트에 데이터를 설정한다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func runScriptCmdSet( data interface{} ) bool {

    jmclog.LogWrite( "Set Command\n" );
	 
	// 수행 조건을 얻는다. 
	m, ok := data.(map[string]interface{})
    if !ok {
	    jmclog.LogWrite( "fail Init Command Syntax Error\n" );
        return false
    }

	component_id              := m["id"].(string)
	component_data            := m["data"]
    component_data_json, err  := json.Marshal( component_data )
	if err != nil {

        jmclog.LogWrite( "fail Convert Comport Data JSON" )
	    RunResult.success  = false;     
	    return false

	}	
	component_data_str := string(component_data_json)

	jmclog.LogWrite( "id     = [%s]\n", component_id );
	jmclog.LogWrite( "data   = [%s]\n", component_data_str );

	set_json  :=   `{ "cmd"  : "set",` ;
	set_json  +=     `"id"   : "` + component_id        + `",`;
	set_json  +=     `"data" : `  + component_data_str  + ` `;
	set_json  +=   `}`
	sendCtrlCmdToComponent( set_json )

    return true
}

