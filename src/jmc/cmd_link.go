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
//   컴포넌트를 연결한다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func runScriptCmdLink( data interface{} ) bool {

    jmclog.LogWrite( "Link Command\n" );
	 
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

	jmclog.LogWrite( "id      = [%s]\n", component_id );
	jmclog.LogWrite( "data   = [%s]\n", component_data_str );

	link_json  :=   `{ "cmd"  : "link",` ;
	link_json  +=     `"id"   : "` + component_id        + `",`;
	link_json  +=     `"data" : `  + component_data_str  + ` `;
	link_json  +=   `}`
	sendCtrlCmdToComponent( link_json )

    return true
}

