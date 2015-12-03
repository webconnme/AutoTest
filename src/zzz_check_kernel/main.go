package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"io"
	"path/filepath"
	"net/smtp"
_ 	"github.com/codeskyblue/go-sh"
	"github.com/scorredoira/email"
        
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
	
	testKernelImagePath    string    // 커널 이미지 패스
	
}

type RunConfig struct {
    prompt               string      // 커널 테스트 중 임을 알리는 프롬프트   
	emailTitle           string      // 메일 제목 
	                     
	testLogFileName      string      // 테스트 로그 파일 이름 : 예) test-2015-13-02T13:12:11.log	
	                     
	testTopPath          string      // 테스트 가장 상위 디렉토리 
	testLogFilePath      string      // 테스트 로그 파일 패쓰 : 예) /home/frog/webconn_auto_test/test-2015-13-02T13:12:11.log	
	                     
}

type Result struct {
    success        bool       // 커널 빌드 중 임을 알리는 프롬프트   
	
	testStartTime string     // 빌드 시  작
	testEndTime   string     // 빌드 종  료
	
}

var Cfg           Config;     // 실행 조건 
var RunEnv        RunConfig;  // 실행 설정 
var RunResult     Result;     // 빌드 결과 

var	logger               *log.Logger                                    // 로거  
 
//---------------------------------------------------------------------------------------------------------------------
//   
//   패쓰 존재 검사 함수 
//   
//---------------------------------------------------------------------------------------------------------------------
func path_exists( path string ) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return true, err
}

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
	
		logger.Printf( "fail run get current directory\n" );
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

    Cfg.title             = "Kernel Test"                        // 테스트 제목
	Cfg.requestTime       = "2015-10-15T10:30:27"                // 테스트 요청일
	Cfg.requestEMail      = "erunsee@falinux.com"                // 테스트 요청자
	Cfg.ccEMail           = []string { "frog@falinux.com" }  // 빌드 참조자
	
	Cfg.testKernelImagePath = "./uImage"
	
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

	RunEnv.emailTitle        = fmt.Sprintf( "[자동빌드] %s", Cfg.title )   // 메일 제목
	
	RunResult.success        = true;                                      // 테스트 최종 결과

	logger.Printf( "Run Top Directory [%s]\n", RunEnv.testTopPath )
	logger.Printf( "test Log File  [%s]\n",   RunEnv.testLogFilePath )
	
    // 커널 이미지가 있는가를 확인한다. 
	logger.Printf( "Check kernel Image [%s]\n", Cfg.testKernelImagePath  )
	
    if is_been,_ := path_exists( Cfg.testKernelImagePath ); is_been != true {
	
	    logger.Printf( "fail Not Found Kernel Image  [%s] \n", Cfg.testKernelImagePath )
		RunResult.success  = false;     // 테스트 실패
		return false
			 
    }
	
    return true	
}

//---------------------------------------------------------------------------------------------------------------------
//   
//   커널을 보드에 다운로드 한다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func kernel_download() bool {

	// 기트 서버에서 커널을 다운로드 한다.
	logger.Printf( "Kernel Image Download [%s]\n", Cfg.testKernelImagePath ) 
	
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
        logger.Println(err)
    }
    
    err = email.Send( Cfg.smtpGmailWithPort, 
	                  smtp.PlainAuth("", Cfg.userGmail, Cfg.passwdGmail, Cfg.smtpGmail), 
   					  m )
    if err != nil {
  		logger.Println(err)
   	} 

}

//---------------------------------------------------------------------------------------------------------------------
//   
//   이 프로그램은 커널 동작을 시험하는 프로그램이다. 
//   
//---------------------------------------------------------------------------------------------------------------------
func main() {

    // 테스트 조건을 초기화 한다. - 초기화 파일에서 읽어 온 후 내부적인 잔여 초기화를 한다. 
    initConfig()

	markTestStartTime()  // 테스트 시작 시간을 표기 한다. 
	getTestTopPath()     // 빌드 상위 패쓰를 구한다. 
	
	// 로그를 할 준비를 한다. 
    //     로그 파일 이름 예 : test-2015-13-02T13:12:11.log	
    RunEnv.prompt            = fmt.Sprintf( ">> %s : ", Cfg.title )          // 테스트 중 임을 알리는 프롬프트   
	RunEnv.testLogFileName  = "test-" + RunResult.testStartTime + ".log"
	RunEnv.testLogFilePath  = RunEnv.testTopPath + RunEnv.testLogFileName
	
    logFile, err := os.OpenFile( RunEnv.testLogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalln("Failed to open log file", RunEnv.testLogFilePath, ":", err)
    } 	
	
	multiLog := io.MultiWriter( logFile, os.Stdout)
	
    logger = log.New( multiLog, RunEnv.prompt, log.Ldate|log.Ltime )
	

	// 테스트 시  작
	logger.Printf( "Test Start [%s]\n", RunResult.testStartTime );
	
	var success bool = true

	if success   {  success    =   prepare()      }    // 작업 디렉토리 준비 
//	if success   {  success    =   kernel_download() } // 커널을 다운로드 한다. 
//	if success   {  success    =   kernel_test() } // 커널을 빌드 한다.  	

	// 빌드 종  료 
	markTestEndTime()  // 빌드 종료 시간을 표기 한다. 

	logger.Printf( "test End [%s]\n", RunResult.testEndTime );
	logFile.Sync()        // 최소한 여기 까지 로그 파일을 저장하고 동기화 한다. 

//	// 결과를 메일로 전달한다. 
//	sendGmail()
//	
	logFile.Close()
	
	if success == false {
	    os.Exit(1)
	} 

}

