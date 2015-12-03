package main

import (
    "time"
	"fmt"
    "strconv"
)	

//  시험을 진행을 한다.
func (jtl *JtlFrame) RunScriptCmdCheck( cmd JtlFrameCommandJson ) bool {

    // 시작 요청된 JEU 가 모든 RUN 으로 바뀌었는가를 확인한다.
	ad.Println( "wait run JEUs" )
 	var loop_out int
	var been     bool
	var err      error
	
 	for loop_out = 100; loop_out > 0 ; loop_out-- {
	    been, err = js.CheckJEUStateStarting()
		
        if err != nil {
            ad.Println( "fail check run of JEU" )
	    	reason := fmt.Sprintf( "do not check run of JEU" )
	    	ar.SetResultError( reason )
            return false
 	    }
		
 	    if !been {
 		    break;
 		}
        time.Sleep(1 * time.Millisecond)
 	}
 	if loop_out == 0 {
 	    ad.Println( "fail run state JEU") 
		reason := fmt.Sprintf( "do not run state JEU") 
		ar.SetResultError( reason )
 		return false
 	}	

    wait_time,_ := strconv.Atoi( cmd.Time )
    ad.Println( "script check command timeout = [%d] msec", wait_time  )
	
	time.Sleep( time.Duration(wait_time) * time.Millisecond )
	
    return true
}

