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
_   "fmt"
    "os"
_	"strconv"
)

import (
        "at"
	atd "at/debug"
	atj "at/jeus" 
)

var ad *atd.AtDebug;
var js *atj.AtJEUServer;

func main() {


	af, _ := at.NewAtFrame(nil); defer af.Close()
	af.DispName    = atj.JS_DEFAULT_NAME
	af.SetId( atj.JS_DEFAULT_NAME )
	
	af.OnKill  = CallbackKill
	af.OnCall  = CallbackCall
	
    ad, _ = atd.NewAtDebugClient( atd.AD_DEFAULT_NAME, af )
	ad.Println( "Program Start..." );
	
	js, _ = atj.NewAtJEUServer( af )
	af.OnReset = CallbackReset

	ret,_ := af.MainLoop()
	
	js.Close()
	ad.Close()
	
	if ret != 0 {
	  os.Exit( ret )
	}

}

// KILL 명령을 무시한다. 
func CallbackKill( af *at.AtFrame )(bool){
	return false
}

// CALL 수신 이벤트 
func CallbackCall( af *at.AtFrame, data interface{} )(bool){

	// 맵 데이터로 변환한다. 
	m, ok := data.(map[string]interface{})
    if !ok {
	    ad.Println( "fail SET command to covert to json")
        return false
    }
	
	cmd := m["cmd"].(string) 
    ad.Println( "cmd = [%s]", cmd )

	switch  cmd {
	case atj.JS_CMD_REGISTER_JEU              : js.RegisterJEUOnServer        ( m["id"].(string) )
	case atj.JS_CMD_UNREGISTER_JEU            : js.UnRegisterJEUOnServer      ( m["id"].(string) )
	case atj.JS_CMD_CHECK_JEU                 : js.CheckJEUOnServer           ( m["id"].(string) ) 
	                                          
	case atj.JS_CMD_SET_READY_JEU             : js.SetJEUStateReadyOnServer   ( m["id"].(string) ) 
	case atj.JS_CMD_SET_STARTING_JEU          : js.SetJEUStateStartingOnServer( m["id"].(string) ) 
	case atj.JS_CMD_SET_RUN_JEU               : js.SetJEUStateRunOnServer     ( m["id"].(string) ) 
	case atj.JS_CMD_SET_STOPPING_JEU          : js.SetJEUStateStoppingOnServer( m["id"].(string) ) 
	case atj.JS_CMD_SET_ENDING_JEU            : js.SetJEUStateEndingOnServer  ( m["id"].(string) ) 
	
	case atj.JS_CMD_CHECK_STATE_STARTING_JEU  : js.CheckEUStateStartingOnServer() 
	}
	
	return false
}

// RESET 수신 이벤트 
func CallbackReset( af *at.AtFrame )(bool){

    js.ResetOnServer()
	return false
}

