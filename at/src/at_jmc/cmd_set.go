package main

//  JEU 를 리셋한다.
func (jtl *JtlFrame) RunScriptCmdSet( cmd JtlFrameCommandJson ) bool {

    ad.Println( "script set command id = [%s] data = [%s]", cmd.Id, cmd.Data )
	af.SendCommandSet( cmd.Id , cmd.Data );

    return true
}
