package main

//  새로운 JEU 를 만든다
func (jtl *JtlFrame) RunScriptCmdLink( cmd JtlFrameCommandJson ) bool {

    ad.Println( "script link command id = [%s] channel = [%s] port = [%s]", cmd.Id, cmd.Channel, cmd.Port  )
	
	af.SendCommandLink( cmd.Id, cmd.Channel, cmd.Port )
	
    return true
}

