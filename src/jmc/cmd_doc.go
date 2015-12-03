package main

import (
   "jmclog"
)

//---------------------------------------------------------------------------------------------------------------------
//   
//   주석을 처리 한다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func runScriptCmdDoc( data interface{} ) bool {

	m, ok := data.(map[string]interface{})
    if !ok {
	    jmclog.LogWrite( "fail DOC Command Syntax Error\n" );
        return false
    }

	jmclog.LogWrite( "%s\n", m["content"] );
	
    return true
}
