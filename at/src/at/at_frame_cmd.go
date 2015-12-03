package at

import (
_    "fmt"
_	"time"
	"encoding/json"
)

func (af *AtFrame) CmdMain( index int ) ( bool, error ) {

    buf, rx_err := af.ZmqPollItems[index].Socket.Recv(0)
	if rx_err != nil {
	    return false, rx_err
	}

    err := json.Unmarshal( buf, &af.cmdJSON )
    if err != nil {
//	    fmt.Printf( "Error : [%s]\n", err );
    }	
		
//	str := string(buf)
//	fmt.Printf( "RX CMD : [%s]\n", str );
//
//	fmt.Printf( "cmd.Cmd  : [%s]\n", af.cmdJSON.Cmd  );
//	fmt.Printf( "cmd.Src  : [%s]\n", af.cmdJSON.Src  );
//	fmt.Printf( "cmd.Dsc  : [%s]\n", af.cmdJSON.Dsc  );
//	fmt.Printf( "cmd.Data : [%s]\n", af.cmdJSON.Data );
	
	if af.cmdJSON.Dsc != "ALL" {
	    if af.cmdJSON.Dsc != af.id {
		    return true, nil
		}
	}
	
	cmd := af.cmdJSON.Cmd
	switch cmd {
    case AF_CMD_KILL    : return af.CmdKill  ()
    case AF_CMD_RESET   : return af.CmdReset ()
    case AF_CMD_SET     : return af.CmdSet   ()
    case AF_CMD_LINK    : return af.CmdLink  ()
    case AF_CMD_UNLINK  : return af.CmdUnlink()
    case AF_CMD_START   : return af.CmdStart ()
    case AF_CMD_STOP    : return af.CmdStop  ()
    case AF_CMD_CALL    : return af.CmdCall  ()
    case AF_CMD_ACK     : return af.CmdAck   ()
	}
	
	return true, nil
}

func (af *AtFrame) CmdKill  () ( bool, error ) {

	if af.OnKill( af ) == true {
	    af.ReqEnd = true
	}

	return true, nil
}

func (af *AtFrame) CmdReset () ( bool, error ) {

    af.OnReset( af )
	
	return true, nil
}

func (af *AtFrame) CmdSet   () ( bool, error ) {

    af.OnSet( af , af.cmdJSON.Data )
	
	return true, nil
}

func (af *AtFrame) CmdLink  () ( bool, error ) {

    af.OnLink( af , af.cmdJSON.Data )
	
	return true, nil
}

func (af *AtFrame) CmdUnlink() ( bool, error ) {

    af.OnUnlink( af , af.cmdJSON.Data )
	
	return true, nil
}

func (af *AtFrame) CmdStart () ( bool, error ) {

    af.OnStart( af )
	
	return true, nil
}

func (af *AtFrame) CmdStop  () ( bool, error ) {

    af.OnStop( af )
	
	return true, nil
}

func (af *AtFrame) CmdCall  () ( bool, error ) {

    af.OnCall( af, af.cmdJSON.Data )
	
	return true, nil
}

func (af *AtFrame) CmdAck   () ( bool, error ) {

    af.OnAck( af )
	
	return true, nil
}
