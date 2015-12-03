package main

import (
    "fmt"
	"time"
	"encoding/json"
)

import (
    "at"
)

import (
  	"github.com/codeskyblue/go-sh"
)

//  새로운 JEU 를 만든다
func (jtl *JtlFrame) RunScriptCmdRun( cmd JtlFrameCommandJson ) bool {

    data_json, err := json.Marshal( cmd.Data )
	if err != nil {
		reason := fmt.Sprintf( "script file syntax error")
		ar.SetResultError( reason )
        return false
 	} 
	
	data_str := string( data_json )
    ad.Println( "script run command id = [%s] path = [%s] data = [%s]", cmd.Id, cmd.Path, data_str  )
    
	// JEU 실행 프로그램 파일이 있는가를 확인한다. 
    if is_been,_ := at.CheckPathExists( cmd.Path ); !is_been {
    
        ad.Println( "fail do not find exec file [%s]", cmd.Path )
		reason := fmt.Sprintf( "do not find exec file [%s]", cmd.Path )
		ar.SetResultError( reason )
        return false
    	 
    }
	
	// 컴포넌트의 실행 프로그램을 수행한다.
    if err := sh.Command( cmd.Path, cmd.Id, data_str ).Start(); err != nil {
        ad.Println( "fail do not exec file [%s]", cmd.Path )
		reason := fmt.Sprintf( "do not exec file [%s]", cmd.Path )
		ar.SetResultError( reason )
        return false
 	}
	
    // 컴포넌트가 등록되었음을 알리는 메세지를 기다린다. 	
	ad.Println( "wait register JEU [%s]", cmd.Id )
 	var loop_out int
	var been     bool
	
 	for loop_out = 10; loop_out > 0 ; loop_out-- {
	    been, err = js.CheckJEU( cmd.Id ); 
		
        if err != nil {
            ad.Println( "fail check been of JEU [%s]", cmd.Id )
	    	reason := fmt.Sprintf( "do not check been of JEU [%s]", cmd.Id )
	    	ar.SetResultError( reason )
            return false
 	    }
		
 	    if been {
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

