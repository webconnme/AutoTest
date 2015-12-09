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

package jeus

import "fmt"

import (
    "at"
)

func NewAtJEUServer( af *at.AtFrame ) (*AtJEUServer, error) {

    js := &AtJEUServer{}
	
	js.AF         = af
	js.ZmqContext = af.ZmqContext
	
	js.ServerMode = true

    js.JEUs       =  make(map[string]*AtjJEUInfo) 	
	
	return js, nil
}

func (js *AtJEUServer) CloseServer() {
    
}

func (js *AtJEUServer) DispJEUs() {
    fmt.Printf( "JEUs.........\n" ) 
    for id, JEU := range js.JEUs {
	    state_str := "" 
	    switch  JEU.State {
        case at.AF_JEU_INIT     : state_str = "INIT    "
        case at.AF_JEU_READY    : state_str = "READY   "
        case at.AF_JEU_STARTING : state_str = "STARTING"
        case at.AF_JEU_RUN      : state_str = "RUN     "
        case at.AF_JEU_STOPPING : state_str = "STOPPING"
        case at.AF_JEU_ENDING   : state_str = "ENDING  "
        case at.AF_JEU_END      : state_str = "END     "
		default                 : state_str = "Unknow  " 
		}
	    fmt.Printf( "id = [%s] state =[%s]\n", id, state_str )
	}
}

func (js *AtJEUServer) ResetOnServer() (error) {
    js.JEUs       =  make(map[string]*AtjJEUInfo) 	
	return nil
}

func (js *AtJEUServer) RegisterJEUOnServer( id string ) (error) {

	js.JEUs[id] = &AtjJEUInfo{Id : id, State : at.AF_JEU_INIT }
	
	js.DispJEUs()
	
	js_data  := AtjJsonDone{ Result : "done" }
	js.AF.SendAck( js_data )

    return nil
}

func (js *AtJEUServer) UnRegisterJEUOnServer( id string ) (error) {

	delete( js.JEUs , id )
	
	js.DispJEUs()
	
	js_data  := AtjJsonDone{ Result : "done" }
	js.AF.SendAck( js_data )
	
    return nil
}

func (js *AtJEUServer) CheckJEUOnServer( id string ) (error) {

	js.DispJEUs()
	
	result := "none"
	_, been := js.JEUs[id];
	if been {
	    result = "been"
	} 
	
	js_data  := AtjJsonDone{ Result : result  }
	js.AF.SendAck( js_data )
	
    return nil
}

func (js *AtJEUServer) SetJEUStateReadyOnServer( id string ) (error) {

    JEU, been := js.JEUs[id];
	if !been {
	    
	} else {
	    JEU.State = at.AF_JEU_READY
	}
	js.DispJEUs()
	
	js_data  := AtjJsonDone{ Result : "done" }
	js.AF.SendAck( js_data )
	
    return nil
}

func (js *AtJEUServer) SetJEUStateStartingOnServer( id string ) (error) {

    JEU, been := js.JEUs[id];
	if !been {
	    
	} else {
	    JEU.State = at.AF_JEU_STARTING
	}
	js.DispJEUs()
	
	js_data  := AtjJsonDone{ Result : "done" }
	js.AF.SendAck( js_data )
	
    return nil
}

func (js *AtJEUServer) SetJEUStateRunOnServer( id string ) (error) {

    JEU, been := js.JEUs[id];
	if !been {
	    
	} else {
	    JEU.State = at.AF_JEU_RUN
	}
	js.DispJEUs()
	
	js_data  := AtjJsonDone{ Result : "done" }
	js.AF.SendAck( js_data )
	
    return nil
}

func (js *AtJEUServer) SetJEUStateStoppingOnServer( id string ) (error) {

    JEU, been := js.JEUs[id];
	if !been {
	    
	} else {
	    JEU.State = at.AF_JEU_STOPPING
	}
	js.DispJEUs()
	
	js_data  := AtjJsonDone{ Result : "done" }
	js.AF.SendAck( js_data )
	
    return nil
}

func (js *AtJEUServer) SetJEUStateEndingOnServer( id string ) (error) {

    JEU, been := js.JEUs[id];
	if !been {
	    
	} else {
	    JEU.State = at.AF_JEU_ENDING
	}
	js.DispJEUs()
	
	js_data  := AtjJsonDone{ Result : "done" }
	js.AF.SendAck( js_data )
	
    return nil
}

func (js *AtJEUServer) CheckEUStateStartingOnServer( ) (error) {

	js.DispJEUs()

	result := "none"
	
    for _, JEU := range js.JEUs {
	    if JEU.State == at.AF_JEU_STARTING {
		    result = "been"
			break
        }		
	}
	
	js_data  := AtjJsonDone{ Result : result  }
	js.AF.SendAck( js_data )
	
    return nil
}
