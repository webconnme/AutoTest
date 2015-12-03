package main

import (
   "jmclog"
   "strconv"
   "time"
)

//---------------------------------------------------------------------------------------------------------------------
//   
//   주석을 처리 한다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func runScriptCmdSleep( data interface{} ) bool {

	m, ok := data.(map[string]interface{})
    if !ok {
	    jmclog.LogWrite( "fail SLEEP Command Syntax Error\n" );
        return false
    }
	
	wait_time := 1000
	
	wait_time, _ = strconv.Atoi(m["time"].(string))
	jmclog.LogWrite( "wait time [%d] m sec\n", wait_time ); 
	
	time.Sleep( time.Duration(wait_time) * time.Millisecond )
	
    return true
}
