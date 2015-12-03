package main

import (
    "fmt"
	"time"
)

//  JEU 를 종료 시킨다. 
func (jtl *JtlFrame) RunScriptCmdKill( cmd JtlFrameCommandJson ) bool {

    ad.Println( "script kill command id = [%s]", cmd.Id )
	af.SendCommandKill( cmd.Id )

    // 컴포넌트가 등록되었음을 알리는 메세지를 기다린다. 	
	ad.Println( "wait unregister JEU [%s]", cmd.Id )
 	var loop_out int
	
 	for loop_out = 10; loop_out > 0 ; loop_out-- {
	    been, err := js.CheckJEU( cmd.Id ); 
		
        if err != nil {
            ad.Println( "fail check been of JEU [%s]", cmd.Id )
	    	reason := fmt.Sprintf( "do not check been of JEU [%s]", cmd.Id )
	    	ar.SetResultError( reason )
            return false
 	    }
		
 	    if !been {
 		    break;
 		}
        time.Sleep(10 * time.Millisecond)
 	}
 	if loop_out == 0 {
 	    ad.Println( "fail register JEU [%s]", cmd.Id ) 
		reason := fmt.Sprintf( "do not register JEU [%s]", cmd.Id ) 
		ar.SetResultError( reason )
 		return false
 	}
	
    return true
}

