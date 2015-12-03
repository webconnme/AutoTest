package main

//  JEU 의 연결된 채널을 해제한다. 
func (jtl *JtlFrame) RunScriptCmdUnlink( cmd JtlFrameCommandJson ) bool {

    ad.Println( "script unlink command id = [%s] channel = [%s] port = [%s]", cmd.Id, cmd.Channel, cmd.Port  )
	
	af.SendCommandUnlink( cmd.Id, cmd.Channel, cmd.Port )
	
    return true
}

