package at

import (
_    "fmt"
	"time"
	"errors"
	"encoding/json"
)

import (
	zmq "github.com/alecthomas/gozmq"
)

func (af *AtFrame) SendCommand( dsc string , cmd string , data interface{} )( error ) {

	af_cmd             := AtFrameCommandJson{ Src : af.id,  Dsc : dsc , Cmd : cmd, Data : data }
    af_cmd_json , err  := json.Marshal( af_cmd )
	if err != nil {
	    return err
	}
	af_cmd_json_str := string(af_cmd_json)
	
    err = af.cmdPUB.Send([]byte(af_cmd_json_str), 0)
	return err
}

func (af *AtFrame) SendCommandKill( id string )( error ) {
    
	err := af.SendCommand( id, AF_CMD_KILL , nil )
    return err
}

func (af *AtFrame) SendCommandReset( id string )( error ) {
    
	err := af.SendCommand( id, AF_CMD_RESET , nil )
    return err
}

func (af *AtFrame) SendCommandSet( id string, data interface{} )( error ) {
    
	err := af.SendCommand( id, AF_CMD_SET , data )
    return err
}

func (af *AtFrame) SendCommandLink( id string, channel string, port string )( error ) {
 
    data := AtFrameLinkJson{ Channel : channel, Port : port }
	err := af.SendCommand( id, AF_CMD_LINK , data )
    return err
}

func (af *AtFrame) SendCommandUnlink( id string, channel string, port string )( error ) {
    
    data := AtFrameLinkJson{ Channel : channel, Port : port }
	err := af.SendCommand( id, AF_CMD_UNLINK , data )
    return err
}


func (af *AtFrame) SendCommandStart( id string )( error ) {
    
	err := af.SendCommand( id, AF_CMD_START , nil )
    return err
}

func (af *AtFrame) SendCommandStop( id string )( error ) {
    
	err := af.SendCommand( id, AF_CMD_STOP , nil )
    return err
}

func (af *AtFrame) SendCall( dsc string, data interface{} , timeout int )( interface{}, error ) {
    
	if dsc != af.lastREQId {
	
	    if af.cmdREQ != nil {
		    af.cmdREQ.Close()
		}
		af.lastREQId  = dsc
		af.cmdREQ , _ = af.ZmqContext.NewSocket(zmq.REQ); 
		af.cmdREQ.Connect( AF_ZMQ_BASE_REP + af.lastREQId )
		
	}
	
	af_cmd  := AtFrameCommandJson{ Src : af.id,  Dsc : dsc, Cmd : AF_CMD_CALL, Data : data }
    af_cmd_json , err  := json.Marshal( af_cmd )
	if err != nil {
	    return nil, err
	}
	
	af_cmd_json_str := string(af_cmd_json)
	
    err = af.cmdREQ.Send([]byte(af_cmd_json_str), 0)
	if err != nil {
	    return nil, err
	}
	
	pi := []zmq.PollItem{ zmq.PollItem{ Socket: af.cmdREQ, Events: zmq.POLLIN} }

    event_count, err := zmq.Poll( pi, time.Millisecond * time.Duration(timeout) )
	if err != nil {
		return nil, err
	}
		
	if event_count == 0 {
		return nil, errors.New( "af call wait timeout" )
	}
	
	buf, rx_err := af.cmdREQ.Recv(0) 
	if rx_err != nil {
	    return nil, rx_err
	}
	
    err = json.Unmarshal( buf, &af.cmdJSON )
    if err != nil {
        return nil, err
    }

//	str := string(buf)
//	fmt.Printf( "CALL RX CMD : [%s]\n", str );
//
//	fmt.Printf( "cmd.Cmd  : [%s]\n", af.cmdJSON.Cmd  );
//	fmt.Printf( "cmd.Src  : [%s]\n", af.cmdJSON.Src  );
//	fmt.Printf( "cmd.Dsc  : [%s]\n", af.cmdJSON.Dsc  );
//	fmt.Printf( "cmd.Data : [%s]\n", af.cmdJSON.Data );	
//
//	fmt.Printf( "CALL END\n" );
	return af.cmdJSON.Data, err
	
}

func (af *AtFrame) SendAck( data interface{} )( error ) {

	af_cmd             := AtFrameCommandJson{ Src : af.id,  Dsc : af.cmdJSON.Src , Cmd : AF_CMD_ACK, Data : data }
    af_cmd_json , err  := json.Marshal( af_cmd )
	if err != nil {
	    return err
	}
	af_cmd_json_str := string(af_cmd_json)
	
    err = af.cmdREP.Send([]byte(af_cmd_json_str), 0)
    return err
}

