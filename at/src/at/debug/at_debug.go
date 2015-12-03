package debug

import (
    "os"
	"log"
_     "fmt"
 _	"time"
)

import (
 	zmq "github.com/alecthomas/gozmq"
)

import (
    "at"
)

const AD_ZMQ_PROXY_PULL  = "ipc:///tmp/at_debug_pull"
const AD_ZMQ_PROXY_PUSH  = "ipc:///tmp/at_debug_pull"

const AD_DEFAULT_NAME    = "atd"

type AtDebug struct {

	ZmqContext     *zmq.Context
	cmdPULL        *zmq.Socket      			// 디버그 데이터 수신 PULL 소켓 
	cmdPUSH        *zmq.Socket      			// 디버그 데이터 송신 PULL 소켓
	pollIndex      int                          // POLL 에 추가된 cmdPULL 소켓에 해당하는 인덱스 번호
	
	ServerMode     bool                         // 서버에서 동작하는가 클라이언트가 동작하는가?
	
	RootPath       string                       // 가장 상위 디렉토리를 갖고 있다. 
	LogPath        string                       // RootPath 하부에 Log를 기록하는 디렉토리 
	FileName       string                       // Log 파일 명 
	
	logFile        *os.File                     // 로그 파일 
	logger         *log.Logger                  // 로거
	
	AF             *at.AtFrame
}

func (ad *AtDebug) Close() {
   
   if ad.ServerMode {
       ad.CloseServer()
   } else {
       ad.CloseClient()
   }
}

