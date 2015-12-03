package main

//  JEU 를 리셋한다.
func (jtl *JtlFrame) RunScriptCmdReset( cmd JtlFrameCommandJson ) bool {

    ad.Println( "script reset command id = [%s]", cmd.Id )
	af.SendCommandReset( cmd.Id );
	
    return true
}
