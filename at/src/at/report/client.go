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

package report

import (
    "fmt"
	"time"
	"strconv"
	"runtime"
)

import (
 	zmq "github.com/alecthomas/gozmq"
)

import (
    "at"
)

func NewAtReportClient( id string, af *at.AtFrame ) (*AtReport, error) {

    ar := &AtReport{}
	
	ar.AF         = af
	ar.ZmqContext = af.ZmqContext
	
	ar.ServerMode = false
	
	ar.cmdPUSH ,  _ = ar.ZmqContext.NewSocket(zmq.PUSH)	
	ar.cmdPUSH.Connect( AR_ZMQ_PROXY_PUSH + id )

	return ar, nil
}

func (ar *AtReport) CloseClient() {

    ar.cmdPUSH.Close()
	
}

func (ar *AtReport) Println( format string, args ...interface{} ) {

    name_tag := "[" + ar.AF.DispName + "] : "
    msg := name_tag + fmt.Sprintf( format, args... ) + "\n" 

    err := ar.cmdPUSH.Send([]byte(msg), 0)
	if err != nil {
	    // fmt.Println( err );
	}
	
}

func (ar *AtReport) Reset() {
	ar.AF.SendCommandReset( AR_DEFAULT_NAME )
}

func (ar *AtReport) Sync() {
	runtime.Gosched()
	time.Sleep( time.Millisecond)
}

func (ar *AtReport) SendSet( data interface{} ) {
	ar.AF.SendCommandSet( AR_DEFAULT_NAME, data )
	ar.Sync()
}

func (ar *AtReport) StartReport( title string ) {

	ar_data            := AtrJsonStart{ Cmd : AR_CMD_START, Title : title }
	ar.SendSet( ar_data )
	
}

func (ar *AtReport) EndReport() {

	ar_data            := AtrJsonEnd{ Cmd : AR_CMD_END }
	ar.SendSet( ar_data )
	
}

func (ar *AtReport) SetTotal( value int ) {

	value_str := strconv.Itoa( value )
	ar_data            := AtrJsonSetTotal{ Cmd : AR_CMD_SET_TOTAL, Value : value_str }
	ar.SendSet( ar_data )
	
}

func (ar *AtReport) SetCurrent( value int ) {

    value_str := strconv.Itoa( value )
	ar_data            := AtrJsonSetCurrent{ Cmd : AR_CMD_SET_CURRENT, Value : value_str }
	ar.SendSet( ar_data )
	
}

func (ar *AtReport) WriteDocument( value string ) {
	
	ar_data            := AtrJsonSetCurrent{ Cmd : AR_CMD_DUCMENT, Value : value }
	ar.SendSet( ar_data )
}

func (ar *AtReport) StartSub( title string ) {

	ar_data            := AtrJsonStartSub{ Cmd : AR_CMD_START_SUB, Title : title }
	ar.SendSet( ar_data )
	
}

func (ar *AtReport) EndSub() {

	ar_data            := AtrJsonEndSub{ Cmd : AR_CMD_END_SUB }
	ar.SendSet( ar_data )
	
}

func (ar *AtReport) SetResultPass() {

	ar_data            := AtrJsonSetResultPass{ Cmd : AR_CMD_SET_PASS }
	ar.SendSet( ar_data )
	
}

func (ar *AtReport) SetResultFail( reason string ) {

	ar_data            := AtrJsonSetResultFail{ Cmd : AR_CMD_SET_FAIL, Reason : reason }
	ar.SendSet( ar_data )
	
}

func (ar *AtReport) SetResultError( reason string ) {

	ar_data            := AtrJsonSetResultError{ Cmd : AR_CMD_SET_ERROR, Reason : reason }
	ar.SendSet( ar_data )
	
}
