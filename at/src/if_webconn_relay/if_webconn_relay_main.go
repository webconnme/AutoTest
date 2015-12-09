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

package main

import (
	"fmt"
	"os"
//	"time"
//	"encoding/json"
)

import (
	zmq "github.com/alecthomas/gozmq"
)

import (
	"at"
	atd "at/debug"
	atr "at/report"
	atj "at/jeus"
	"encoding/json"
	"time"
)

var af  *at.AtFrame;
var ad  *atd.AtDebug;
var ar  *atr.AtReport;
var js  *atj.AtJEUServer;

var serverIP string
var serverPort string

func main() {


	myJEUId := os.Args[1]
	option  := []byte(os.Args[2])

	af, _ = at.NewAtFrame(nil);
	af.DispName = "if_webconn_relay"
	af.SetId( myJEUId )

	af.OnKill   = CallbackKill;

	af.OnPeriod = CallbackPeriod;

	af.OnReset  = CallbackReset
	af.OnSet    = CallbackSet
	af.OnLink   = CallbackLink
	af.OnUnlink = CallbackUnlink
	af.OnStart  = CallbackStart
	af.OnStop   = CallbackStop

	ad, _ = atd.NewAtDebugClient( atd.AD_DEFAULT_NAME, af )
	ad.Println( "Program Start..." )

	ar, _ = atr.NewAtReportClient( atr.AR_DEFAULT_NAME, af )

	js, _ = atj.NewAtJEUClient( af )

	var argv_options map[string]interface{}
	if err := json.Unmarshal(option, &argv_options); err != nil {
		ad.Println( "fail option syntex error [%s]", err )
		reason := fmt.Sprintf( "option syntex error [%s]", err )
		ar.SetResultError( reason )
		close()
		return
	}

	if err := js.RegisterJEU( af.GetId() ); err != nil {
		ad.Println( "fail register JEU [%s]", af.GetId() )
		reason := fmt.Sprintf( "do not register JEU [%s]", af.GetId() )
		ar.SetResultError( reason )

		close()
		return
	}

	serverIP = argv_options["ip"].(string)
	serverPort = argv_options["port"].(string)
	ad.Println( "Relay Server = [%s:%s]", serverIP, serverPort )

	PairSocket, _ = af.ZmqContext.NewSocket(zmq.PAIR)
	PairSocket.Connect("tcp://" + serverIP + ":"+ serverPort)

	//	OpenRS232()

	js.SetJEUStateReady( af.GetId() )
	ret,_ := af.MainLoop()

	if ret != 0 {
		close()
		os.Exit( ret )
	}

	close()
}

func close() {
	ad.Println( "Program Close..." )

	CloseRelay()

	if err := js.UnRegisterJEU( af.GetId() ); err != nil {
		ad.Println( "fail unregister JEU [%s]", af.GetId() )
		reason := fmt.Sprintf( "do not unregister JEU [%s]", af.GetId() )
		ar.SetResultError( reason )
	}

//	if PortTxSUB != nil {
//		PortTxSUB.Close()
//	}
//
//	if PortRxPUB != nil {
//		PortRxPUB.Close()
//	}

	js.Close()
	ar.Close()
	ad.Close()
	af.Close()
}

var callbackCount = 0

func CallbackKill( af *at.AtFrame )(bool){
	ad.Println( "callback kill" )

	js.SetJEUStateEnding( af.GetId() )

	return true
}

func CallbackPeriod( af *at.AtFrame )(bool){
	callbackCount++;
	if callbackCount % 1000 == 0 {
		ad.Println( "callback live = %d", ThreadRelayLive )
	}

	return false
}

// RESET 수신 이벤트
func CallbackReset( af *at.AtFrame )(bool){
	ad.Println( "callback reset" )


	m := Message{ "power", "off" }
	str, _ := json.Marshal([]Message{m})
	PairSocket.Send(str, zmq.NOBLOCK)

	return false
}

type CkmCommand struct {
	Cmd     string     // 처리 명령
	Value   string     // 값
}
// SET 수신 이벤트
func CallbackSet( af *at.AtFrame, p_data interface{} )(bool){
	ad.Println( "callback set" )
	//	ad.Println( "set data [%s]", p_data )

	// 맵 데이터로 변환한다.
	m, ok := p_data.(map[string]interface{})
	if !ok {
		ad.Println( "fail SET command to covert to json")
		return false
	}

	cmd    := m["cmd"  ].(string)
	value  := m["value"].(string)

	commands = append( commands, CkmCommand{ Cmd : cmd, Value : value } )
	i := len( commands )
	ad.Println( "commands[%d].Cmd = [%s] commands[%d].Value = [%s]", i-1,commands[i-1].Cmd, i-1,commands[i-1].Value )

	return false
}

// LINK 수신 이벤트
func CallbackLink( af *at.AtFrame, p_data interface{} )(bool){
	ad.Println( "callback link data = [%s]", p_data )

	// 맵 데이터로 변환한다.
	m, ok := p_data.(map[string]interface{})
	if !ok {
		ad.Println( "fail LINK command to covert to json")
		return false
	}

	channel    := m["channel"  ].(string)
	port       := m["port"].(string)

	ad.Println( "channel  = [%s]", channel   )
	ad.Println( "port     = [%s]", port      )

	if port ==  "RX DATA" {

//		if PortRxPUB != nil {
//			PortRxPUB.Close()
//		}
//
//		// 제어 수신용 zmq SUB 소켓을 만든다.
//		var err error
//		PortRxPUB, err = af.ZmqContext.NewSocket(zmq.PAIR)
//		if err != nil {
//
//			ad.Println( "fail do not create channel socket [%s:%s]", channel, port )
//			reason := fmt.Sprintf( "do not create channel socket [%s:%s]", channel, port )
//			ar.SetResultError( reason )
//			return false
//		}
//
//		// 채널에 바인드 한다.
//		PortRxPUB.Bind( at.AF_ZMQ_CHANNEL + channel )

	}

	if port ==  "TX DATA" {

//		if PortTxSUB != nil {
//			PortTxSUB.Close()
//		}
//
//		// 제어 수신용 zmq SUB 소켓을 만든다.
//		var err error
//		PortTxSUB, err = af.ZmqContext.NewSocket(zmq.PAIR)
//		if err != nil {
//
//			ad.Println( "fail do not create channel socket [%s:%s]", channel, port )
//			reason := fmt.Sprintf( "do not create channel socket [%s:%s]", channel, port )
//			ar.SetResultError( reason )
//			return false
//		}
//
//		// 채널에 바인드 한다.
//		PortTxSUB.Connect( at.AF_ZMQ_CHANNEL + channel )
//		PortTxSUB.SetSubscribe("")
	}

	return false
}

// UNLINK 수신 이벤트
func CallbackUnlink( af *at.AtFrame, p_data interface{}  )(bool){
	ad.Println( "callback unlink data = [%s]", p_data )

	// 맵 데이터로 변환한다.
	m, ok := p_data.(map[string]interface{})
	if !ok {
		ad.Println( "fail UNLINK command to covert to json")
		return false
	}

	channel    := m["channel"  ].(string)
	port       := m["port"].(string)

	ad.Println( "channel  = [%s]", channel   )
	ad.Println( "port     = [%s]", port      )

//	if port ==  "RX DATA" {
//		if PortRxPUB != nil {
//			PortRxPUB.Close()
//			PortRxPUB = nil
//		}
//	}
//
//	if port ==  "TX DATA" {
//		if PortTxSUB != nil {
//			PortTxSUB.Close()
//			PortTxSUB = nil
//		}
//	}

	return false
}


// START 수신 이벤트
func CallbackStart( af *at.AtFrame )(bool){
	ad.Println( "callback start" )
	js.SetJEUStateStarting( af.GetId() )

	ThreadCheckMsgRun    = false
	go ThreadCheckMsg()
	time.Sleep( time.Millisecond )

	for {
		time.Sleep( time.Millisecond )
		if ThreadCheckMsgRun {
			break
		}
	}

	js.SetJEUStateRun( af.GetId() )

	return false
}

// STOP 수신 이벤트
func CallbackStop( af *at.AtFrame )(bool){
	ad.Println( "callback stop" )
	js.SetJEUStateStopping( af.GetId() )

	ThreadCheckMsgReqEnd = true

	for {
		time.Sleep( time.Millisecond )
		if !ThreadCheckMsgRun {
			break
		}
	}

	js.SetJEUStateReady( af.GetId() )

	return false
}

func CloseRelay() {
	PairSocket.Close()
}
