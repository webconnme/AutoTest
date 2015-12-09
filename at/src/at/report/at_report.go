/**
 * The MIT License (MIT)
 *
 * Copyright (c) 2015 David You <david@webconn.me>
 * Copyright (c) 2015 Victor Kim <victor@webconn.me>
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
    "os"
	"log"
)

import (
 	zmq "github.com/alecthomas/gozmq"
)

import (
    "at"
	"net"
)

const AR_ZMQ_PROXY_PULL  = "ipc:///tmp/at_report_pull"
const AR_ZMQ_PROXY_PUSH  = "ipc:///tmp/at_report_pull"

const AR_DEFAULT_NAME    = "atr"

// AR Frame Command 
const (
    AR_CMD_START        = "start"
    AR_CMD_END          = "end"
    AR_CMD_SET_TOTAL    = "set_total"
	AR_CMD_SET_CURRENT  = "set_current"
	AR_CMD_DUCMENT      = "document"
    AR_CMD_START_SUB    = "start_sub"
    AR_CMD_END_SUB      = "end_sub"
	
	AR_CMD_SET_PASS     = "set_pass"   // 시험을 통과 했다. 
	AR_CMD_SET_FAIL     = "set_fail"   // 시험에 실패 했다.
	AR_CMD_SET_ERROR    = "set_error"  // 시험중 에러가 발생했다.
)

// AR Frame Command JSON
type AtrJsonStart struct {
    Cmd     string      `json:"cmd"`     // 명령 
	Title   string      `json:"title"`   // 시험 타이틀
}

type AtrJsonEnd struct {
    Cmd     string      `json:"cmd"`     // 명령 
}

type AtrJsonSetTotal struct {
    Cmd     string      `json:"cmd"`     // 명령 
	Value   string      `json:"value"`   // 전체 항목 수
}

type AtrJsonSetCurrent struct {
    Cmd     string      `json:"cmd"`     // 명령 
	Value   string      `json:"value"`   // 현재 진행 수
}

type AtrJsonDocument struct {
    Cmd     string      `json:"cmd"`     // 명령 
	Value   string      `json:"value"`   // 도큐먼트 
}

type AtrJsonStartSub struct {
    Cmd     string      `json:"cmd"`     // 명령 
	Title   string      `json:"title"`   // 항목 타이틀
}

type AtrJsonEndSub struct {
    Cmd     string      `json:"cmd"`     // 명령 
}

type AtrJsonSetResultPass struct {
    Cmd     string      `json:"cmd"`     // 명령 
}

type AtrJsonSetResultFail struct {
    Cmd     string      `json:"cmd"`     // 명령 
	Reason   string     `json:"reason"`   // 실패에 대한 데이터
}

type AtrJsonSetResultError struct {
    Cmd     string      `json:"cmd"`     // 명령 
	Reason   string     `json:"reason"`   // 에러에 대한 데이터
}

type AtReport struct {

	ZmqContext     *zmq.Context
	cmdPULL        *zmq.Socket      			// 디버그 데이터 수신 PULL 소켓 
	cmdPUSH        *zmq.Socket      			// 디버그 데이터 송신 PULL 소켓
	pollIndex      int                          // POLL 에 추가된 cmdPULL 소켓에 해당하는 인덱스 번호
	
	ServerMode     bool                         // 서버에서 동작하는가 클라이언트가 동작하는가?
	
	RootPath       string                       // 가장 상위 디렉토리를 갖고 있다. 
	ReportPath     string                       // RootPath 하부에 레포트를 기록하는 디렉토리 
	FileName       string                       // Report 파일 명 
	
	rptFile        *os.File                     // 레포트 파일 
	reporter       *log.Logger                  // 레포트

	AF             *at.AtFrame

	toWeb		   chan []byte
	manager		   chan net.Conn

	Title          string                       // 레포트 제목 
	Total          int                          // 총 항목 수 
	Current        int                          // 진행 항목 수 
	
	Deep           int                          // 시험 단계 깊이
}

func (ar *AtReport) Close() {
   
   if ar.ServerMode {
       ar.CloseServer()
   } else {
       ar.CloseClient()
   }
}

