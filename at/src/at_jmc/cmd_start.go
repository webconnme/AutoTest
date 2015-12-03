package main

//  JEU 의 시험을 시작한다.
func (jtl *JtlFrame) RunScriptCmdStart( cmd JtlFrameCommandJson ) bool {

    ad.Println( "script start command id = [%s]", cmd.Id )
	af.SendCommandStart( cmd.Id );
	
    return true
}
