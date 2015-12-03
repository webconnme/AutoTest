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

