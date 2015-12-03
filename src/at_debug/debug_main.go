package main

//import "fmt"

import (
	"os"
)

import (
	"at"
	atd "at/debug"
)

var ad *atd.AtDebug

func main() {

	//  fmt.Println( "This is Main\n" )
	//	at.TestFunc()

	af, _ := at.NewAtFrame(nil)
	defer af.Close()

	af.DispName = atd.AD_DEFAULT_NAME
	af.SetId(atd.AD_DEFAULT_NAME)

	af.OnKill = CallbackKill
	af.OnRxIn = CallbackRxIn
	af.OnReset = CallbackReset

	ad, _ = atd.NewAtDebugServer(atd.AD_DEFAULT_NAME, af)

	ret, _ := af.MainLoop()
	ad.Close()

	if ret != 0 {
		os.Exit(ret)
	}

}

// KILL 명령을 무시한다.
func CallbackKill(af *at.AtFrame) bool {
	return false
}

// ZMQ 수신 이벤트
func CallbackRxIn(af *at.AtFrame, index int, data []byte) bool {

	ad.CallbackServerRxIn(af, index, data)

	return false
}

// RESET 수신 이벤트
func CallbackReset(af *at.AtFrame) bool {

	ad.ResetServer()
	return false
}
