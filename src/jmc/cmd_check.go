package main

import (
   "strconv"
   "time"
)

import (
   "jmclog"
   "jmcsvr"
)
	


//---------------------------------------------------------------------------------------------------------------------
//   
//   주석을 처리 한다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func runScriptCmdCheck( data interface{} ) bool {

	m, ok := data.(map[string]interface{})
    if !ok {
	    jmclog.LogWrite( "fail CHECK Command Syntax Error\n" )
        return false
    }
	
	time_out, _ := strconv.Atoi(m["time"].(string))
	jmclog.LogWrite( "wait time [%d] m sec\n", time_out ); 
	
	
	
	ok = jmcsvr.SetWaitTestEnd()
	if ok == false {
	    jmclog.LogWrite( "fail server Set wait test end\n" );
        return false
	}
	
	// 1 m sec 마다 끝났는가를 확인한다. 
	start_time    := time.Now()
	time_out_msec := time.Duration( time_out ) * time.Millisecond 
	
	for {
	    // 시간 초과를 계산한다.  
	    current_time := time.Now()
		pass_time := current_time.Sub( start_time )
		
//		jmclog.LogWrite( "wait time [%d:%d]\n", pass_time , time_out_msec ); 
		
		if pass_time > time_out_msec {
			jmclog.LogWrite( "fail Check Command Time out!\n" ); 
			RunResult.success  = false;     
			return false
		}
		
		// 시험이 종료 되었는가를 확인한다. 
	    if jmcsvr.CheckWaitTestEnd() {
		    break
		}
     	time.Sleep( time.Millisecond )
	}	
	
	jmclog.LogWrite( "Check End ------------------------- \n"); 
	
    return true
}
