package main

import (
    "fmt"
    "os" 
	"errors"
	"io/ioutil"
    "path/filepath"
	"encoding/json"
)	

import (
    "at" 
)	

// JTL Frame Command JSON
type JtlFrameCommandJson struct {
	Cmd     string      `json:"cmd"`     // 명령
	Id      string      `json:"id"`      // Cmd 명령 대상 구별 ID  "ALL" 은 예약어 이며 모든 프레임워크가 받아 들인다.  
	Txt     string      `json:"txt"`     // Cmd 의 명령에 필요한 텍스트 정보가 필요할 때 사용 
	Time    string      `json:"time"`    // Cmd 의 명령에 시간적인 값이 필요할 때 사용
	Path    string      `json:"path"`    // Cmd 의 명령에 파일 패쓰가 필요할 때 사용
	Channel string      `json:"channel"` // Cmd 의 link , unlink 의 채널 구별 ID 
	Port    string      `json:"port"`    // Cmd 의 link , unlink 의 채널 연결 Port 
	Data    interface{} `json:"data"`    // 명령에 따른 데이터 
}

type JtlFrameScript struct {
	Title     string                   `json:"title"`        // 시험 항목
	Descript  string                   `json:"descript"`     // 시험 설명
	Commands  []JtlFrameCommandJson    `json:"task"`         // 시험 태스크
}

type JtlFrame struct {
   TopPath  string             // 가장 상위 디렉토리 
   FileName string             // JTL 파일 이름 
   
   Script   JtlFrameScript     // JTL 스크립트
   Current  int                // 현재 진행 중인 스크립트 항목 인덱스
}

func NewJtlFrame() (*JtlFrame, error) {

	jtl := &JtlFrame{}
	
	if dir, err := filepath.Abs(filepath.Dir(os.Args[0])); err  != nil {
		return nil, err;
    } else {
	    jtl.TopPath = dir + "/"
	}	
	
	ad.Println( "jtl.TopPath = [%s]", jtl.TopPath )

	jtl.Current = 0
	
	return jtl, nil
}

func (jtl *JtlFrame) Close() {
    
	
}

func (jtl *JtlFrame) LoadScript( filename string ) ( error ) {

    jtl.FileName = filename
	ad.Println( "script filename  = [%s]", jtl.FileName )

    if is_been, _ := at.CheckPathExists( jtl.FileName ); !is_been {

		ad.Println( "fail do not exist script file = [%s]", jtl.FileName )
		return errors.New( "fail read script file" )
			 
    }	

    if b,err := ioutil.ReadFile( jtl.FileName ); err != nil {
	
	    ad.Println( "fail read script file = [%s]", jtl.FileName )
		return errors.New( "fail read script file" )
		
	} else {

	    if err := json.Unmarshal( b, &(jtl.Script) ); err != nil {
             ad.Println( "fail convert script file = [%s]", jtl.FileName )
	    	 ad.Println( "err : [%s]", err )
	    	 return errors.New( "fail read script file" )
        }
	}	
	
    return nil
}

func (jtl *JtlFrame) RunScript() ( error ) {

    ad.Println( "run script" )

	var err error = nil
	success  := true
	commands := jtl.Script.Commands
	
	for index,one_cmd := range commands {
	    if success == false {
	    	err = errors.New("")
			break
	    }
		
		ar.SetCurrent ( index )

		switch one_cmd.Cmd {
		
		case "doc"    : success = jtl.RunScriptCmdDoc   ( one_cmd ) 
		case "sleep"  : success = jtl.RunScriptCmdSleep ( one_cmd ) 
		
		case "sub"    : success = jtl.RunScriptCmdSub   ( one_cmd ) 
		case "subend" : success = jtl.RunScriptCmdSubEnd( one_cmd ) 
		
		case "end"    : ad.Println( "script end command"  )
		                goto test_end
					 
		case "run"    : success = jtl.RunScriptCmdRun   ( one_cmd )
		case "kill"   : success = jtl.RunScriptCmdKill  ( one_cmd )

		case "reset"  : success = jtl.RunScriptCmdReset ( one_cmd ) 
		case "set"    : success = jtl.RunScriptCmdSet   ( one_cmd ) 
		
		case "link"   : success = jtl.RunScriptCmdLink  ( one_cmd ) 
		case "unlink" : success = jtl.RunScriptCmdUnlink( one_cmd ) 

		case "start"  : success = jtl.RunScriptCmdStart ( one_cmd )
		case "stop"   : success = jtl.RunScriptCmdStop  ( one_cmd )

		case "check"  : success = jtl.RunScriptCmdCheck ( one_cmd )

		default       : ad.Println( "script unknow command"  )
						reason := fmt.Sprintf( "script unknow command")
						ar.SetResultError( reason )
						err = errors.New("")
		                goto test_end
		}
    
	}
	
test_end:	

    ad.Println( "end script" )
    return err
}

func (jtl *JtlFrame) GetTitle() ( string ) {

    return jtl.Script.Title

}

func (jtl *JtlFrame) GetTotal() ( int ) {

    return len( jtl.Script.Commands )

}

func (jtl *JtlFrame) GetCurrent() ( int ) {

    return jtl.Current

}
