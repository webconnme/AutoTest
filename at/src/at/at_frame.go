/**
 * The MIT License (MIT)
 *
 * Copyright (c) 2015 
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

package at

import (
    "fmt"
	"time"
	"os"
)

import (
        "github.com/satori/go.uuid"
	zmq "github.com/alecthomas/gozmq"
)

func TestFunc( ) {
    fmt.Println( "Call Test Func" )
}

const AF_ZMQ_PROXY_XPUB = "ipc:///tmp/at_frame_pub"
const AF_ZMQ_PROXY_XSUB = "ipc:///tmp/at_frame_sub"
const AF_ZMQ_PROXY_PUB  = "ipc:///tmp/at_frame_pub"
const AF_ZMQ_PROXY_SUB  = "ipc:///tmp/at_frame_sub"

const AF_ZMQ_BASE_REP   = "ipc:///tmp/at_frame_rep_"

const AF_ZMQ_CHANNEL    = "ipc:///tmp/at_channel_"

// AT Frame Command 
const (
    AF_CMD_KILL    = "kill"
    AF_CMD_RESET   = "reset"
    AF_CMD_SET     = "set"
    AF_CMD_LINK    = "link"
    AF_CMD_UNLINK  = "unlink"
    AF_CMD_START   = "start"
    AF_CMD_STOP    = "stop"
    AF_CMD_CALL    = "call"
    AF_CMD_ACK     = "ack"
)

// AT Frame Command JSON
type AtFrameCommandJson struct {
	Src     string      `json:"src"`     // 송신자 구별 ID         
	Dsc     string      `json:"dsc"`     // 수신자 구별 ID         "ALL" 은 예약어 이며 모든 프레임워크가 받아 들인다.  
	Cmd     string      `json:"cmd"`     // 명령
	Data    interface{} `json:"data"`    // 명령에 따른 데이터 
}

type AtFrameLinkJson struct {
	Channel string      `json:"channel"` // 채널 구분 ID         
	Port    string      `json:"port"`    // 포트 구분 ID
}

type AtFrameOnCallPeriodFunc    func( af *AtFrame )(bool)

type AtFrameOnCallKillFunc      func( af *AtFrame )(bool)
type AtFrameOnCallResetFunc     func( af *AtFrame )(bool)

type AtFrameOnCallSetFunc       func( af *AtFrame , data interface{} )(bool)
type AtFrameOnCallCallFunc      func( af *AtFrame , data interface{} )(bool)
type AtFrameOnCallAckFunc       func( af *AtFrame )(bool)

type AtFrameOnCallLinkFunc      func( af *AtFrame , data interface{} )(bool)
type AtFrameOnCallUnlinkFunc    func( af *AtFrame , data interface{} )(bool)

type AtFrameOnCallStartFunc     func( af *AtFrame )(bool)
type AtFrameOnCallStopFunc      func( af *AtFrame )(bool)


type AtFrameOnRxInFunc          func( af *AtFrame, index int, data []byte )(bool)

// RUN STATE
const (
    AF_JEU_INIT       = 0        // 프로그램 시작 중 - RegisterJEU() 함수 호출에 의해 설정
    AF_JEU_READY      = 1        // 명령 대기 중     - SetJEUStateReady() 함수 호출에 의해 설정
    AF_JEU_STARTING   = 2        // 시작 중          - SetJEUStateStarting() 함수 호출에 의해 설정 시작은 START 명령에 의해서 시작
    AF_JEU_RUN        = 3        // 동작 중          - SetJEUStateRun() 함수 호출에 의해 설정 START 명령에 따른 동작이 모두 시작되었을때 설정
    AF_JEU_STOPPING   = 4        // 중지 중  여기서 AF_STATE_READY 로 이동 
	                             //                  - SetJEUStateStopping() 함수 호출에 의해 설정  STOP 명령에 따른 동작이 모두 종료 되었을때 설정 
    AF_JEU_ENDING     = 5        // 종료 중          - SetJEUStateEnding() 함수 호출에 의해 설정    KILL 명령에 따라서 UnRegisterJEU() 함수가 호출될때까지 유지
    AF_JEU_END        = 6        // 종료             - UnRegisterJEU() 함수 호출로 이 상태 발생 
	                             //                    jes 에서 이상태는 관리하지 않으므로 실제로는 존재하지 않는 상태  
)

type AtFrame struct {

    Parent         *AtFrame                     // 부모
    Uuid           uuid.UUID        			// 고유 구별 객체
    DispName       string           			// 표출용 프레임 메인 이름
	id             string                       // 프레임 구별 ID 
	Period         time.Duration    			// 주기적 호출 주기 디폴트 1msec
	               
	ZmqContext     *zmq.Context
	cmdSUB         *zmq.Socket      			// 프레임 명령 수신 SUB 소켓 
	cmdPUB         *zmq.Socket      			// 프레임 명령 송신 PUB 소켓 

	cmdREP         *zmq.Socket      			// 프레임 명령 수신 REP 소켓 
    lastREQId      string                       // 마지막 연결 REQ 바인드 Id 스트링
	cmdREQ         *zmq.Socket      			// 프레임 명령 송신 REQ 소켓 
	
	ZmqPollItems   []zmq.PollItem               // ZMQ 폴 구조
	               
	ReqEnd         bool                         // 프레임 동작 종료 요청
	
	cmdJSON        AtFrameCommandJson           // 프레임에 수신된 명령 
	                                            
	OnPeriod	   AtFrameOnCallPeriodFunc      // 주기적 호출 콜백 
	                                            
    OnKill         AtFrameOnCallKillFunc        // AT 프레임 종료 명령 처리전 호출 
    OnReset        AtFrameOnCallResetFunc       // AT 프레임 리셋 
    OnSet          AtFrameOnCallSetFunc         // AT 프레임 설정
    OnLink         AtFrameOnCallLinkFunc        // AT 프레임간의 연결
    OnUnlink       AtFrameOnCallUnlinkFunc      // AT 프레임간의 연결 해제
    OnStart        AtFrameOnCallStartFunc       // AT 프레임 동작 시작
    OnStop         AtFrameOnCallStopFunc        // AT 프레임 동작 강제 중지
    OnCall         AtFrameOnCallCallFunc        // AT 프레임 ACK 응답 형 요구 
    OnAck          AtFrameOnCallAckFunc         // AT 프레임 값 요구에 대한 응답
	
	OnRxIn         AtFrameOnRxInFunc            // ZMQ 수신 이벤트 

	RunState       int                          // 수행 상태 
}

func (af *AtFrame) DispInfo() {

    fmt.Printf( "AF Information\n" )
    fmt.Printf( "  Name   = [%s]\n"      , af.DispName                  )
    fmt.Printf( "  Id     = [%s]\n"      , af.id                        )
    fmt.Printf( "  Uuid   = [%s]\n"      , af.Uuid                      )
    fmt.Printf( "  Period = [%d] msec\n" , af.Period / time.Millisecond )

}

func NewAtFrame( parent *AtFrame ) (*AtFrame, error) {

	af := &AtFrame{}
	
	af.Uuid          = uuid.NewV4()
	af.Period        = time.Millisecond
	
	af.Parent = parent
	if af.Parent == nil {
	    af.ZmqContext, _ = zmq.NewContext()
	} else {
	   af.ZmqContext = af.Parent.ZmqContext
	}
	
	af.cmdSUB ,    _ = af.ZmqContext.NewSocket(zmq.SUB)	
	af.cmdPUB ,    _ = af.ZmqContext.NewSocket(zmq.PUB)	
	
	af.cmdPUB.Connect( AF_ZMQ_PROXY_SUB  )
	af.cmdSUB.Connect( AF_ZMQ_PROXY_PUB  )
	af.cmdSUB.SetSubscribe("")
	
	af.ZmqPollItems   = []zmq.PollItem{ zmq.PollItem{ Socket: af.cmdSUB, Events: zmq.POLLIN} }

    af.lastREQId      = ""
	
	af.ReqEnd        = false

	af.OnPeriod      = onPeriod

    af.OnKill        = onKill
    af.OnReset       = onReset
    af.OnSet         = onSet
    af.OnLink        = onLink
    af.OnUnlink      = onUnlink
    af.OnStart       = onStart 
    af.OnStop        = onStop  
    af.OnCall        = onCall
    af.OnAck         = onAck

	af.OnRxIn        = onRxIn
	
	return af, nil
}

func (af *AtFrame) Close() {

    if af.cmdREQ != nil {
	    af.cmdREQ.Close()
	}
	
    if af.cmdREP != nil {
	    af.cmdREP.Close()
	}

	if af.cmdSUB != nil {
	    af.cmdSUB.Close()
	}
	
	if af.cmdPUB != nil {
	    af.cmdPUB.Close()
	}

	if af.Parent == nil {
        if af.ZmqContext != nil {
	        af.ZmqContext.Close()
	    }
	} 
	
}

func (af *AtFrame) MainLoop() ( int, error) {

    af.ReqEnd = false
	for !af.ReqEnd {
	
	    pi := af.ZmqPollItems
		
		event_count, err := zmq.Poll( pi, af.Period )
		if err != nil {
			break
		}
		
		if event_count == 0 {
		
			if af.OnPeriod != nil {
				if af.OnPeriod( af ) {
				    af.ReqEnd = true
				}
			}
			
		} else {
		
            if pi[0].REvents&zmq.POLLIN != 0 {
		    	if _, err = af.CmdMain(0); err != nil {
		    		af.ReqEnd = true
		    	}
		    }
			
            if pi[1].REvents&zmq.POLLIN != 0 {
		    	if _, err = af.CmdMain(1); err != nil {
		    		af.ReqEnd = true
		    	}
		    }

		    for i := 2; i < len( pi); i++ {
                if pi[i].REvents&zmq.POLLIN != 0 {
		    	    if _, err = af.RxIn( i ); err != nil {
		    	    	af.ReqEnd = true
		    	    }
		        }
			}
			
		}
	}	
			
    return 0, nil
}

func (af *AtFrame) RxIn( index int ) ( bool, error ) {

    buf, rx_err := af.ZmqPollItems[index].Socket.Recv(0)
	if rx_err != nil {
	    return false, rx_err
	}

	if af.OnRxIn( af, index, buf ) == true {
	    af.ReqEnd = true
	}

	return true, nil
}

func (af *AtFrame) Sleep( msec int ) {

    time.Sleep( time.Duration(msec) * time.Millisecond)	
	
}

func CheckPathExists( path string ) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return true, err
}

func (af *AtFrame) SetId( id string ) {

    af.id = id;
	af.cmdREP ,    _ = af.ZmqContext.NewSocket(zmq.REP); 
	af.cmdREP.Bind( AF_ZMQ_BASE_REP + af.id  )
	
	af.ZmqPollItems  = append( af.ZmqPollItems, zmq.PollItem{ Socket: af.cmdREP, Events: zmq.POLLIN} )
	
}

func (af *AtFrame) GetId() ( string ) {

    return af.id
}

func (af *AtFrame) AppendZmqPollItem( item zmq.PollItem )(int,error){

    af.ZmqPollItems    = append( af.ZmqPollItems, item )
    
	index := len( af.ZmqPollItems ) -1
	
    return index, nil
	
}
