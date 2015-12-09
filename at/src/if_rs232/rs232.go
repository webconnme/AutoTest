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

package main

import (
    "fmt"
    "time"
    "runtime"
)

import (
	"github.com/mikepb/go-serial"
	zmq "github.com/alecthomas/gozmq"
)

var ThreadRS232RxReqEnd bool    = false
var ThreadRS232TxReqEnd bool    = false
var ThreadRS232RxRun    bool    = false
var ThreadRS232TxRun    bool    = false
var ThreadRS232Live   uint64  = 0 

var RS232Options = serial.Options{
    BitRate:     115200,
    DataBits:    8,
    Parity:      serial.PARITY_NONE,
    StopBits:    1,
    FlowControl: serial.FLOWCONTROL_NONE,
	Mode: serial.MODE_READ_WRITE,
}

var RS232Path    string
var RS232Port    *serial.Port
var PortRxPUB    *zmq.Socket = nil
var PortTxSUB    *zmq.Socket = nil

func ThreadRS232Rx() {
    
	ad.Println( "ThreadRS232Rx() start" )
	
	ThreadRS232RxReqEnd = false
    ThreadRS232RxRun    = true
	
	for !ThreadRS232RxReqEnd {
	
	    ThreadRS232Live++  
	    runtime.Gosched()

        remain, err := RS232Port.InputWaiting()
 		if err != nil {
            ad.Println( "fail do not read wait[%s]", err )
	        reason := fmt.Sprintf( "do not read wait[%s]", err )
	        ar.SetResultError( reason )
 			break;
 		}

 		if remain != 0 {
		
			ad.Println( "remain = %d", remain ); 

 		    buf := make([]byte, remain)

 		    length, err := RS232Port.Read(buf)
 		    if err != nil {
                ad.Println( "fail do not read wait[%s]", err )
	            reason := fmt.Sprintf( "do not read wait[%s]", err )
	            ar.SetResultError( reason )
 		    	break;
 		    }
			ad.Println( "read length [%d]", length );

			err = PortRxPUB.Send(buf[:length], 0)
            if err != nil {
                ad.Println( "fail do not send channel[%s]", err )
	            reason := fmt.Sprintf( "do not send channel[%s]", err )
	            ar.SetResultError( reason )
 		    	break;
            }

 		}			
		
	}

	ad.Println( "ThreadRS232Rx() End" )
    ThreadRS232RxRun    = false
	
}

func ThreadRS232Tx() {
    
	ad.Println( "ThreadRS232Tx() start" )
	
	ThreadRS232TxReqEnd = false
    ThreadRS232TxRun    = true
	
	for !ThreadRS232TxReqEnd {
	
	    ThreadRS232Live++  
		
//		ad.Println( "wait read PortTxSUB" )      
        pi := zmq.PollItems{
    	     zmq.PollItem{ Socket: PortTxSUB, Events: zmq.POLLIN},
    	}

        event_count, err := zmq.Poll( pi, 1 * time.Millisecond )
   		if err != nil {
            ad.Println( "fail do not poll[%s]", err )
	        reason := fmt.Sprintf( "do not poll[%s]", err )
	        ar.SetResultError( reason )
 			break
			   
    	}
			
		if event_count == 0 {
			
		} else {
            if pi[0].REvents&zmq.POLLIN != 0 {
	            buf, err := pi[0].Socket.Recv(0)
   	            if err != nil {
                     ad.Println( "fail do not read [%s]", err )
	                 reason := fmt.Sprintf( "do not read [%s]", err )
	                 ar.SetResultError( reason )
 		    	     break;
                }
				
				RS232Port.Write( buf )
		                   
		        str := string(buf)
		        ad.Println( "OUT ASCII : [%s]\n", str )
			}	
		}
		
	}

	ad.Println( "ThreadRS232Tx() End" )
	ThreadRS232TxRun    = false

}

func OpenRS232() {

	options := RS232Options

	var err error
	
	if RS232Port != nil {
	    RS232Port.Close()
	}
	
    RS232Port, err = options.Open( RS232Path )
    if err != nil {
        ad.Println( "fail do not open RS232[%s]", RS232Path )
	    reason := fmt.Sprintf( "do not open RS232[%s]", RS232Path )
	    ar.SetResultError( reason )
    }
	
}

func CloseRS232() {

	if RS232Port != nil {
	    RS232Port.Close()
	}

}

