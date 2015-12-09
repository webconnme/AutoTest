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
    "time"
	"fmt"
    "strconv"
)	

//  시험을 진행을 한다.
func (jtl *JtlFrame) RunScriptCmdCheck( cmd JtlFrameCommandJson ) bool {

    // 시작 요청된 JEU 가 모든 RUN 으로 바뀌었는가를 확인한다.
	ad.Println( "wait run JEUs" )
 	var loop_out int
	var been     bool
	var err      error
	
 	for loop_out = 100; loop_out > 0 ; loop_out-- {
	    been, err = js.CheckJEUStateStarting()
		
        if err != nil {
            ad.Println( "fail check run of JEU" )
	    	reason := fmt.Sprintf( "do not check run of JEU" )
	    	ar.SetResultError( reason )
            return false
 	    }
		
 	    if !been {
 		    break;
 		}
        time.Sleep(1 * time.Millisecond)
 	}
 	if loop_out == 0 {
 	    ad.Println( "fail run state JEU") 
		reason := fmt.Sprintf( "do not run state JEU") 
		ar.SetResultError( reason )
 		return false
 	}	

    wait_time,_ := strconv.Atoi( cmd.Time )
    ad.Println( "script check command timeout = [%d] msec", wait_time  )
	
	time.Sleep( time.Duration(wait_time) * time.Millisecond )
	
    return true
}

