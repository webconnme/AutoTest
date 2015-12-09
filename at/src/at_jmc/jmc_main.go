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

// import "fmt"

import (
    "os"
)

import (
        "at"
	atd "at/debug"
	atr "at/report"
	atj "at/jeus"
)

var af  *at.AtFrame
var ad  *atd.AtDebug
var ar  *atr.AtReport
var js  *atj.AtJEUServer

var jtl *JtlFrame

func main() {

	af, _ = at.NewAtFrame(nil);

	af.DispName = "jmc"
	af.SetId( "jmc" )

//	af.OnPeriod = CallbackPeriod;

    ad, _ = atd.NewAtDebugClient( atd.AD_DEFAULT_NAME, af )
	ad.Reset()
	ad.Println( "Program Start..." )

    ar, _ = atr.NewAtReportClient( atr.AR_DEFAULT_NAME, af )

	js, _ = atj.NewAtJEUClient( af )
	js.Reset()

	jtl, _ = NewJtlFrame(); defer jtl.Close()

    if len(os.Args) != 2 {
		ad.Println( "fail do not exists script file name on command line argument" );
		goto jmc_end
	}

	if err := jtl.LoadScript( os.Args[1] ); err != nil {
		ad.Println( "fail load script [%s]" , err )
		goto jmc_end
	}

	ar.StartReport( jtl.GetTitle() )
	ar.SetTotal   ( jtl.GetTotal() )

	if err := jtl.RunScript(); err != nil {
		ad.Println( "fail run script [%s]" , err )
		goto jmc_end
	}

	ar.SetResultPass()

//	ret,_ := af.MainLoop()
//	
//	if ret != 0 {
//	  ad.Close()
//	  os.Exit( ret )
//	}

jmc_end:

	ar.EndReport()
	close()
	
}

func close() {
    ad.Println( "Program Close..." )
	jtl.Close()
    js.Close()
	ar.Close()
	ad.Close()
	af.Close()
}


func CallbackPeriod( af *at.AtFrame )(bool){
	return true
}

