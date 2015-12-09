/**
 * The MIT License (MIT)
 *
 * Copyright (c) 2015 David You <david@webconn.me>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 * 
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 * 
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

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

