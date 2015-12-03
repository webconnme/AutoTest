package main

//  JEU 의 시험을 중지한다.
func (jtl *JtlFrame) RunScriptCmdStop( cmd JtlFrameCommandJson ) bool {

    ad.Println( "script stop command id = [%s]", cmd.Id )
	af.SendCommandStop( cmd.Id );
	
    return true
}
