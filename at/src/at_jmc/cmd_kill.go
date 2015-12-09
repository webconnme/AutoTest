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
)

//  JEU 를 종료 시킨다. 
func (jtl *JtlFrame) RunScriptCmdKill( cmd JtlFrameCommandJson ) bool {

    ad.Println( "script kill command id = [%s]", cmd.Id )
	af.SendCommandKill( cmd.Id )

    // 컴포넌트가 등록되었음을 알리는 메세지를 기다린다. 	
	ad.Println( "wait unregister JEU [%s]", cmd.Id )
 	var loop_out int
	
 	for loop_out = 10; loop_out > 0 ; loop_out-- {
	    been, err := js.CheckJEU( cmd.Id ); 
		
        if err != nil {
            ad.Println( "fail check been of JEU [%s]", cmd.Id )
	    	reason := fmt.Sprintf( "do not check been of JEU [%s]", cmd.Id )
	    	ar.SetResultError( reason )
            return false
 	    }
		
 	    if !been {
 		    break;
 		}
        time.Sleep(10 * time.Millisecond)
 	}
 	if loop_out == 0 {
 	    ad.Println( "fail register JEU [%s]", cmd.Id ) 
		reason := fmt.Sprintf( "do not register JEU [%s]", cmd.Id ) 
		ar.SetResultError( reason )
 		return false
 	}
	
    return true
}

