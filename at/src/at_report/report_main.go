package main

// import "fmt"

import (
    "os"
	"strconv"
)

import (
        "at"
	atd "at/debug"
	atr "at/report"
)

var ad *atd.AtDebug;
var ar *atr.AtReport;

func main() {

	af, _ := at.NewAtFrame(nil); defer af.Close()
	
	af.DispName    = atr.AR_DEFAULT_NAME
	af.SetId( atr.AR_DEFAULT_NAME )
	
	af.OnKill  = CallbackKill
	af.OnReset = CallbackReset
	af.OnSet   = CallbackSet
	af.OnRxIn  = CallbackRxIn
	af.OnCall  = CallbackCall
	
    ad, _ = atd.NewAtDebugClient( atd.AD_DEFAULT_NAME, af )
	ad.Println( "Program Start..." );
	
	ar, _ = atr.NewAtReportServer( atr.AR_DEFAULT_NAME, af )

	ret,_ := af.MainLoop()
	
	ar.Close()
	ad.Close()
	
	if ret != 0 {
	  os.Exit( ret )
	}

}

// KILL 명령을 무시한다. 
func CallbackKill( af *at.AtFrame )(bool){
	return false
}

// ZMQ 수신 이벤트 
func CallbackRxIn( af *at.AtFrame, index int, data []byte )(bool){

	ar.CallbackServerRxIn( af, index, data )

    return false
}

// RESET 수신 이벤트 
func CallbackReset( af *at.AtFrame )(bool){

    ar.ResetServer()
	return false
}

// SET 수신 이벤트 
func CallbackSet( af *at.AtFrame, data interface{} )(bool){

	// 맵 데이터로 변환한다. 
	m, ok := data.(map[string]interface{})
    if !ok {
	    ad.Println( "fail SET command to covert to json")
        return false
    }
	
	cmd := m["cmd"].(string) 
    ad.Println( "cmd = [%s]", cmd )

	switch  cmd {
	case atr.AR_CMD_START       : ar.StartReportOnServer( m["title"].(string) )
	case atr.AR_CMD_END         : ar.EndReportOnServer()
	
	case atr.AR_CMD_SET_TOTAL   : value,_ := strconv.Atoi( m["value"].(string) )
	                              ar.SetTotalOnServer      ( value )
								  
	case atr.AR_CMD_SET_CURRENT : value,_ := strconv.Atoi( m["value"].(string) )
	                              ar.SetCurrentOnServer    ( value )
								  
	case atr.AR_CMD_DUCMENT		: value,_ := m["value"].(string)
	                              ar.WriteDocumentOnServer ( value )
								  
	case atr.AR_CMD_START_SUB   : ar.StartSubOnServer      ( m["title"].(string) )
	case atr.AR_CMD_END_SUB     : ar.EndSubOnServer        ()
	
	case atr.AR_CMD_SET_PASS    : ar.SetResultPassOnServer ()
	case atr.AR_CMD_SET_FAIL    : ar.SetResultFailOnServer ( m["reason"].(string) )
	case atr.AR_CMD_SET_ERROR   : ar.SetResultErrorOnServer( m["reason"].(string) )
								  
	}
	
	return false
}

// CALL 수신 이벤트 
func CallbackCall( af *at.AtFrame, data interface{} )(bool){

	af.SendAck( nil )
	return false
}
