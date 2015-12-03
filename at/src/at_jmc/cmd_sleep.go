package main

import (
    "time"
    "strconv"
)	

//   밀리 초 동안 멈춘다. 
func (jtl *JtlFrame) RunScriptCmdSleep( cmd JtlFrameCommandJson ) bool {

    wait_time,_ := strconv.Atoi( cmd.Time )
    ad.Println( "script sleep command time = [%d] msec", wait_time )
	time.Sleep( time.Duration(wait_time) * time.Millisecond )
	
    return true
}
