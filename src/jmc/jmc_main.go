package main

import (
	"fmt"
	"os"
	"time"
	"io/ioutil"
	"path/filepath"
	"net/smtp"
_ 	"github.com/codeskyblue/go-sh"
	"github.com/scorredoira/email"
	"encoding/json"
	zmq "github.com/alecthomas/gozmq"
	
	"jmclog"
	"jmcsvr"
	
)

type Config struct {
    title                  string      // 빌드 제목
	requestTime            string      // 빌드 요청일 : 2015-10-15 10:30:27
	requestEMail           string      // 빌드 요청자 : frog@falinux.com
	ccEMail                []string    // 빌드 참조자 : boggle70@falinux.com
	
	testerEMail            string     // 빌드 이메일 
	smtpGmailWithPort      string     // Gmail SMTP 정보     "smtp.gmail.com:587"
	smtpGmail              string     // Gmail SMTP 정보     "smtp.gmail.com"
	userGmail              string     // Gmail USER 정보     "webconn.autobuild@gmail.com" 
	passwdGmail            string     // Gmail PASSWORD 정보 "2001May09"
	
}

type JtlCmdDoc struct {
	Content   string
}

type JtlTask struct {
	Cmd     string      `json:"cmd"`      // 시험 진행 명령
	Data    interface{} `json:"data"`     // 시험 진행 데이터
}

type JtlContent struct {
	Title     string       `json:"title"`        // 시험 항목
	Descript  string       `json:"descript"`     // 시험 설명
	Task     []JtlTask     `json:"task"`         // 시험 태스크
}

type RunConfig struct {
    prompt               string      // 커널 테스트 중 임을 알리는 프롬프트   
	emailTitle           string      // 메일 제목 
	testLogFileName      string      // 테스트 로그 파일 이름 : 예) test-2015-13-02T13:12:11.log	
	                     
	testTopPath          string      // 테스트 가장 상위 디렉토리 
	testLogFilePath      string      // 테스트 로그 파일 패쓰 : 예) /home/frog/webconn_auto_test/test-2015-13-02T13:12:11.log	
	
	jtlFilename          string      // 테스트 진행 스크립트 파일 이름 
	jtlContent           JtlContent  // 수행 스크립트            
	
	testTitle            string      // 가장 상위의 시험 항목 이름 
	testDescript         string      // 가장 상위의 시험 항목 설명 
}

type Result struct {
    success        bool       // 커널 빌드 중 임을 알리는 프롬프트   
	
	testStartTime string     // 빌드 시  작
	testEndTime   string     // 빌드 종  료
	
}

var Cfg           Config     // 실행 조건 
var RunEnv        RunConfig  // 실행 설정 
var RunResult     Result     // 빌드 결과 

var zmqContext   *zmq.Context
var zmqCtrlPUB   *zmq.Socket                                 // 제어용 PUB 소켓 

//---------------------------------------------------------------------------------------------------------------------
//   
//   시작 시간 마크
//   
//---------------------------------------------------------------------------------------------------------------------
func markTestStartTime() {

    t := time.Now()

	RunResult.testStartTime = fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
                                      t.Year(), t.Month(), t.Day(),
                                      t.Hour(), t.Minute(), t.Second())

}

//---------------------------------------------------------------------------------------------------------------------
//   
//   종료 시간 마크
//   
//---------------------------------------------------------------------------------------------------------------------
func markTestEndTime() {

    t := time.Now()

	RunResult.testEndTime = fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
                                      t.Year(), t.Month(), t.Day(),
                                      t.Hour(), t.Minute(), t.Second())
}

//---------------------------------------------------------------------------------------------------------------------
//   
//   빌드 상위 패쓰를 구한다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func getTestTopPath() {
    
	
    // 현재 실행되는 디렉토리를 구한다. 
	if dir, err := filepath.Abs(filepath.Dir(os.Args[0])); err  != nil {
	
		jmclog.LogWrite( "fail run get current directory\n" );
		os.Exit(1)
		
    } else {
	
	    RunEnv.testTopPath =  dir + "/"
		
	}	
}

//---------------------------------------------------------------------------------------------------------------------
//   
//   실행 조건을 준비한다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func initConfig() {

    Cfg.title             = "jmc"                        // 테스트 제목
	
// ------------------------------------------------------------------------------------	
	
	Cfg.requestTime       = "2015-10-15T10:30:27"                // 테스트 요청일
	Cfg.requestEMail      = "erunsee@falinux.com"                // 테스트 요청자
	Cfg.ccEMail           = []string { "frog@falinux.com" }  // 빌드 참조자
	
	Cfg.testerEMail         = "webconn.autobuild@gmail.com"       // 테스트 이메일 
	
	Cfg.smtpGmailWithPort = "smtp.gmail.com:587"           // Gmail SMTP 정보 
	Cfg.smtpGmail         = "smtp.gmail.com"               // Gmail SMTP 정보 
	Cfg.userGmail         = "webconn.autobuild@gmail.com"  // Gmail USER 정보 
	Cfg.passwdGmail       = "2001May09"                    // Gmail PASSWORD 정보 
	
}

//---------------------------------------------------------------------------------------------------------------------
//   
//   작업 준비를 한다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func prepare() bool {

	RunResult.success        = true;                                      // 테스트 최종 결과

	jmclog.LogWrite( "Run Top Directory [%s]\n", RunEnv.testTopPath )
	jmclog.LogWrite( "Test Log File  [%s]\n",   RunEnv.testLogFilePath )


    // 명령 라인에 스크립트 파일명이 있는가를 검사 한다. 
	jmclog.LogWrite( "Check Script File\n" )

	if len(os.Args) != 2 {
	    jmclog.LogWrite( "fail No Command Line Argument Script Filename\n" )
		RunResult.success  = false;     // 테스트 실패
		return false
	} 
	
	// 수행 시험 스크립트 파일 이름 
    RunEnv.jtlFilename = os.Args[1]	
	
	jmclog.LogWrite( "SCRIPT Filename [%s]\n", RunEnv.jtlFilename );

	// 스크립트가 존재 하는가를 확인한다. 
    if is_been,_ := path_exists( RunEnv.jtlFilename ); is_been != true {

	    jmclog.LogWrite( "fail Not Found Script File  [%s] \n", RunEnv.jtlFilename )
		RunResult.success  = false;     // 테스트 실패
		return false
			 
    }
	
	// 스크립트 파일을 읽어 들인다. 

    if b,err := ioutil.ReadFile( RunEnv.jtlFilename ); err != nil {
	
	    jmclog.LogWrite( "fail Read  Script File  [%s] \n", RunEnv.jtlFilename )
		RunResult.success  = false;     // 테스트 실패
		return false
		
	} else {
	
	    // 내부 자료 구조체로 변경한다.  
        err := json.Unmarshal( b, &(RunEnv.jtlContent) )
	    if err != nil {

	        jmclog.LogWrite( "fail Convert Script File  [%s] \n", RunEnv.jtlFilename )
	        jmclog.LogWrite( "%s\n", err )
		    RunResult.success  = false;     // 테스트 실패
		    return false

	    }	
	}
	
	// 기본 정보를 재 설정 한다. 
	RunEnv.testTitle    =  RunEnv.jtlContent.Title       // 가장 상위의 시험 항목 이름 
	RunEnv.testDescript =  RunEnv.jtlContent.Descript    // 가장 상위의 시험 항목 설명 
	
	jmclog.LogWrite( "Top Test Title    : [%s]\n", RunEnv.testTitle );
	jmclog.LogWrite( "Top Test Descript : [%s]\n", RunEnv.testDescript );
	

	RunEnv.emailTitle        = fmt.Sprintf( "[자동시험] %s", RunEnv.testTitle )   // 메일 제목

    return true	
}

//---------------------------------------------------------------------------------------------------------------------
//   
//   스크립트를 실행한다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func runScript() bool {

 	// 테스트 시  작
 	jmclog.LogWrite( "Scrtipt Start\n" );
	
    jtlTask := RunEnv.jtlContent.Task 
	
	var success bool
	
	success = true
	
	for _,task := range jtlTask {
		
		if success == false {
		    break
		}
		
		jmclog.LogWrite( "CMD : [%s]\n", task.Cmd );
		
		switch task.Cmd {
		
		case "doc"    : success = runScriptCmdDoc( task.Data ) 
		case "sleep"  : success = runScriptCmdSleep( task.Data ) 
		
		case "end"    : goto runScriptEnd
					 
		case "run"    : jmclog.LogWrite( "component create\n" );
                        success = runScriptCmdRun( task.Data ) 
						 
		case "kill"   : jmclog.LogWrite( "component destroy\n" );
		                success = runScriptCmdKill( task.Data )  
		
		case "link"   : jmclog.LogWrite( "component link\n" );
		                success = runScriptCmdLink( task.Data )
						
		case "unlink" : jmclog.LogWrite( "component unlink\n" );
		                success = runScriptCmdUnlink( task.Data )
		
		case "init"   : jmclog.LogWrite( "component data init\n" );
		                success = runScriptCmdInit( task.Data )
						
		case "set"    : jmclog.LogWrite( "component data setting\n" );
		                success = runScriptCmdSet( task.Data ) 
		
		case "start"  : jmclog.LogWrite( "component start\n" );
		                success = runScriptCmdStart( task.Data )
						
		case "stop"   : jmclog.LogWrite( "component stop\n" );
		                success = runScriptCmdStop( task.Data ) 
						
		case "check"  : jmclog.LogWrite( "test end wait\n" );
		                success = runScriptCmdCheck( task.Data ) 
		
		default       : jmclog.LogWrite( "unknow command\n" );
	                    success = false	
		}
		
	}
   
runScriptEnd:   
   
    return true
}

//---------------------------------------------------------------------------------------------------------------------
//   
//   GMail 을 통해서 빌드 결과를 메일로 보낸다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func sendGmail() {

	// 메일 제목 : "[자동 빌드] : 시험용 빌드 - 성공"
	
	var MailTitle   string;
	
	if RunResult.success {
	    MailTitle = RunEnv.emailTitle + " - 성공";
	} else {
	    MailTitle = RunEnv.emailTitle + " - 실패";
	}

	// 메일 내용 : 아래는 예
    //             테스트 요청일 : 2015-10-15T10:30:27
    //             테스트 요청자 : frog@falinux.com
    //             테스트 시  작 : 2015-10-15T10:30:27 
    //             테스트 종  료 : 2015-10-15T10:30:27 
    //             테스트 결  과 : 성공

	var MailContent string;
	
	MailContent += "테스트 요청일 : " + Cfg.requestTime          + "\n";
	MailContent += "테스트 요청자 : " + Cfg.requestEMail         + "\n";
	MailContent += "테스트 시  작 : " + RunResult.testStartTime + "\n";
	MailContent += "테스트 종  료 : " + RunResult.testEndTime   + "\n";
	
	if RunResult.success {
	    MailContent += "테스트 결  과 : " + "성공";
	} else {
	    MailContent += "테스트 결  과 : " + "실패";
	}

    m := email.NewMessage( MailTitle, MailContent )

    m.From = Cfg.testerEMail
    m.To = []string{ Cfg.requestEMail }
    m.Cc = Cfg.ccEMail
    
    err := m.Attach( RunEnv.testLogFilePath )
    if err != nil {
        jmclog.LogWrite("fail Can not Attach Log file to Gmail [%s] \n", err )
    }
    
    err = email.Send( Cfg.smtpGmailWithPort, 
	                  smtp.PlainAuth("", Cfg.userGmail, Cfg.passwdGmail, Cfg.smtpGmail), 
   					  m )
    if err != nil {
		jmclog.LogWrite("fail Can not Send to Gmail [%s] \n", err )
   	} 

}

//---------------------------------------------------------------------------------------------------------------------
//   
//   컴포넌트에 명령을 보낸다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func sendCtrlCmdToComponent( json_str string ) {

    jmclog.LogWrite( "sendCtrlCmdToComponent() json_str [%s]\n", json_str );
	
    err := zmqCtrlPUB.Send([]byte(json_str), 0)
    if err != nil {
    	jmclog.LogWrite( "zmq ctrl send: ", err  );	
    }
    
}

//---------------------------------------------------------------------------------------------------------------------
//   
//   현재 등록된 모든 컴포넌트에 종료 명령 KILL 을 보낸다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func sendCtrlCmdKillToAllComponent() {

	id_list, err := jmcsvr.GetAllComponetIdList()
	if err != nil {
	    jmclog.LogWrite( "jmcsvr.GetAllComponetIdList() err [%s]\n", err );
	}	
	jmclog.LogWrite( "jmcsvr.GetAllComponetIdList() id_list [%s]\n", id_list );
	
	for _, id := range id_list {
	    kill_json := `{ "cmd" : "kill", "id" : "` + id + `" }`
		sendCtrlCmdToComponent( kill_json )
	}

}

//---------------------------------------------------------------------------------------------------------------------
//   
//   모든 등록된 컴포넌트가 제거 되었는가를 확인한다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func checkEmptyRegisterComponent() {

    for {
	    id_list, err := jmcsvr.GetAllComponetIdList()
	    if err != nil {
	        jmclog.LogWrite( "jmcsvr.GetAllComponetIdList() err [%s]\n", err );
	    }	
	    jmclog.LogWrite( "jmcsvr.GetAllComponetIdList() id_list [%s]\n", id_list );
		if len( id_list ) == 0 {
		    break
		}
		time.Sleep(10 * time.Millisecond)
	}
    
}

//---------------------------------------------------------------------------------------------------------------------
//   
//   이 프로그램은 커널 동작을 시험하는 프로그램이다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func main() {

	zmqContext, _ = zmq.NewContext()
    
    //  수행 초기화를 처리 한다. 
    initConfig()
	
    markTestStartTime()  // 동작 시작 시간을 표기 한다. 
    getTestTopPath()     // 수행 가장 상위 패스를 구한다. 
	  
 	// 로그를 할 준비를 한다. 
    //     로그 파일 이름 예 : test-2015-13-02T13:12:11.log	
    RunEnv.prompt            = Cfg.title                                       // 실행 프로그램 프롬프트   
 	RunEnv.testLogFileName   = "jmc-" + RunResult.testStartTime + ".log"
 	RunEnv.testLogFilePath   = RunEnv.testTopPath + RunEnv.testLogFileName

	// 로그 시작  
	jmclog.Start( jmclog.SERVER,            // 서버 모드 
	              RunEnv.testLogFilePath ,  // 로그 파일 이름 지정
				  RunEnv.prompt,            // 로그 프롬프트 지정
				  zmqContext )              // zmq 컨텍스트 
	            
    // 컴포넌트 관리 서버를 시작한다. 
	
	jmcsvr.Start( jmcsvr.SERVER,            // 서버 모드 
				  zmqContext )              // zmq 컨텍스트 

				  
	// 컴포넌트 관리 서버 에코 테스트 
	jmcsvr.Echo()


    // 제어 방송용 zmq 	를 만든다. 
	zmqCtrlPUB, _ = zmqContext.NewSocket(zmq.PUB)
	zmqCtrlPUB.Bind("ipc:///tmp/jmc_ctrl")
	
 	// 테스트 시작 
	
 	jmclog.LogWrite( "Test Start [%s]\n", RunResult.testStartTime );

 	var success bool = true
	
 	if success   {  success    =   prepare()      }  // 스크립트 파일 체크 및 작업 디렉토리 준비 
	if success   {  success    =   runScript()    }  // 스크립트를 실행한다. 

 	// 수행 종  료 
 	markTestEndTime()  // 빌드 종료 시간을 표기 한다. 
 
 	jmclog.LogWrite( "test End [%s]\n", RunResult.testEndTime );
 
	// 결과를 메일로 전달한다. 
//	sendGmail()


    // 현재 등록된 모든 컴포넌트에 종료 명령 DEL 을 보낸다. 
//	echo_json := `{ "cmd" : "echo" }`
//	sendCtrlCmdToComponent( echo_json )	
    sendCtrlCmdKillToAllComponent()
	
    // 모든 레지스터가 종료되었는가를 확인한다. 
    checkEmptyRegisterComponent()	

	// 잔여 처리를 수행한다. 
	jmclog.LogWrite( "jmc All Close\n" );
	
	// zmq 의 대기열 처리를 위해서 1 초간 잠든다. 
//	time.Sleep(1000 * time.Millisecond)
	
	zmqCtrlPUB.Close()

	jmcsvr.End()
    jmclog.End()
	
	zmqContext.Close()
 	if success == false {
 	    os.Exit(1)
 	} 

}

