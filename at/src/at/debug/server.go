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
	"log"
    "io"
	"os"
	"path/filepath"
	"time"
)

import (
 	zmq "github.com/alecthomas/gozmq"
)

import (
    "at"
)

func NewAtDebugServer( id string, af *at.AtFrame ) (*AtDebug, error) {

    ad := &AtDebug{}
	
	ad.AF         = af
	ad.ZmqContext = af.ZmqContext
	
	ad.ServerMode = true
	
	ad.cmdPULL ,  _ = ad.ZmqContext.NewSocket(zmq.PULL)	
	ad.cmdPULL.Bind( AD_ZMQ_PROXY_PULL + id )
	
	ad.pollIndex, _ = af.AppendZmqPollItem( zmq.PollItem{ Socket: ad.cmdPULL, Events: zmq.POLLIN} )

	if dir, err := filepath.Abs(filepath.Dir(os.Args[0])); err  != nil {
		return nil, err;
    } else {
	    ad.RootPath =  dir + "/"
	}	

	ad.LogPath   = "log/"
	ad.FileName  = "log"

	ad.logFile = nil
	
	ad.ResetServer()
	
	return ad, nil
}

func (ad *AtDebug) CloseServer() {

   if ad.logFile != nil {
       ad.logFile.Close()
   }
   
   ad.cmdPULL.Close()
}

func (ad *AtDebug) CallbackServerRxIn( af *at.AtFrame, index int, data []byte )(bool){

	if ad.pollIndex == index {
	    ad.logger.Printf( string(data) )
	}

    return false	
}

func (ad *AtDebug) ReopenLogFile()( error ) {

    if ad.logFile != nil {
        ad.logFile.Close()
    }

    t := time.Now()
	ad.FileName  = fmt.Sprintf("log-%d-%02d-%02dT%02d:%02d:%02d",
                               t.Year(), t.Month(), t.Day(),  t.Hour(), t.Minute(), t.Second())
	
	log_path     := ad.RootPath + ad.LogPath
    log_filename := ad.RootPath + ad.LogPath + ad.FileName
    fmt.Printf( "\nlog_file_name = [%s]\n", log_filename )	
	
	var err error
	
    if is_been,_ := at.CheckPathExists( log_path ); !is_been {
	    if err = os.MkdirAll( log_path, 0777 ); err != nil {
	    	 return err
        }
    }
	
    ad.logFile, err = os.OpenFile( log_filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
	    return err
    }

    return nil   
}

func (ad *AtDebug) ResetServer()( error ) {

    err := ad.ReopenLogFile()
	if err != nil {
	    return err
	}
	
 	multiLog := io.MultiWriter( ad.logFile, os.Stdout)
 	
    ad.logger = log.New( multiLog, ">> ", log.Ldate|log.Ltime )
	
	return err
}

