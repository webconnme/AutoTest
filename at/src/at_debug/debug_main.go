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
