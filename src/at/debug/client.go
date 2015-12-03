package debug

import (
    "fmt"
	"time"
	"runtime"
)

import (
 	zmq "github.com/alecthomas/gozmq"
)

import (
    "at"
)

func NewAtDebugClient( id string, af *at.AtFrame ) (*AtDebug, error) {

    ad := &AtDebug{}
	
	ad.AF         = af
	ad.ZmqContext = af.ZmqContext
	
	ad.ServerMode = false
	
	ad.cmdPUSH ,  _ = ad.ZmqContext.NewSocket(zmq.PUSH)	
	ad.cmdPUSH.Connect( AD_ZMQ_PROXY_PUSH + id )

	return ad, nil
}

func (ad *AtDebug) CloseClient() {

    ad.cmdPUSH.Close()
	
}

func (ad *AtDebug) Println( format string, args ...interface{} ) {

    name_tag := "[" + ad.AF.DispName + "] : "
    msg := name_tag + fmt.Sprintf( format, args... ) + "\n" 

    err := ad.cmdPUSH.Send([]byte(msg), 0)
	if err != nil {
	    // fmt.Println( err );
	}
    
}

func (ad *AtDebug) Reset() {
	ad.AF.SendCommandReset( AD_DEFAULT_NAME )
	time.Sleep( time.Millisecond)
}

func (ad *AtDebug) Trace() {

    pc := make([]uintptr, 10)  // at least 1 entry needed
    runtime.Callers(2, pc)
    f := runtime.FuncForPC(pc[0])
    file, line := f.FileLine(pc[0])

	ad.Println( "%s:%d %s\n", file, line, f.Name() )
    
}

