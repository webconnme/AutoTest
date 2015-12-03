package main

import (
 	"github.com/codeskyblue/go-sh"
	"time"
	"encoding/json"
	
	"jmclog"
	js "jmcsvr"
)

//---------------------------------------------------------------------------------------------------------------------
//   
//   새로운 컴포넌트를 만든다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func runScriptCmdRun( data interface{} ) bool {

   jmclog.LogWrite( "Run Command\n" );

	// 수행 조건을 얻는다. 
	m, ok := data.(map[string]interface{})
    if !ok {
	    jmclog.LogWrite( "fail RUN Command Syntax Error\n" );
        return false
    }

	component_id              := m["id"].(string)
	component_path            := m["path"].(string)
	component_option          := m["option"]
    component_option_json, err := json.Marshal( component_option )
	if err != nil {

        jmclog.LogWrite( "fail Convert Comport Option JSON" )
	    RunResult.success  = false;     
	    return false

	}	
	component_option_str := string(component_option_json)
	
	jmclog.LogWrite( "id     = [%s]\n", component_id );
	jmclog.LogWrite( "path   = [%s]\n", component_path );
	jmclog.LogWrite( "option = [%s]\n", component_option_str );
		
    // 컴포넌트의 실행 프로그램 파일이 있는가를 확인한다. 	
	jmclog.LogWrite( "check component path [%s]", component_path );
    if is_been,_ := path_exists( component_path ); !is_been {
	
		 jmclog.LogWrite( "fail Find Not Componenet File [%s]", component_path )
		 RunResult.success  = false;     // 빌드 실패
		 return false
			 
    }
	
	// 컴포넌트의 실행 프로그램을 수행한다. 
	jmclog.LogWrite( "exec component program [%s]", component_path ) 

    err = sh.Command( component_path, component_id, component_option_str ).Start()
//    err = sh.Command( component_path, component_id, "'" + component_option_str + "'"  ).Start()
	if err != nil {
	    jmclog.LogWrite( "fail exec component program [%s]", component_path ) 
        RunResult.success  = false;     // 실행 실패
		return false
	}
	
    // 컴포넌트가 등록되었음을 알리는 메세지를 기다린다. 	
	jmclog.LogWrite( "wait register component [%s:%s]", component_path, component_id ) 
	var loop_out int
	for loop_out = 10; loop_out > 0 ; loop_out-- {
	    if js.CheckBeenComponet( component_id ) {
		    break;
		}
        time.Sleep(10 * time.Millisecond)
	}
	if loop_out == 0 {
	    jmclog.LogWrite( "fail exec component program [%s] No Register Component", component_path ) 
        RunResult.success  = false;     // 실행 실패
		return false
	}
	
    jmclog.LogWrite( "Sucess Execute Component Program [%s:%s]", component_path, component_id ) 	
	
    return true
}

