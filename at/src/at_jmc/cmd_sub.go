package main

//  서브 시험 항목 시작
func (jtl *JtlFrame) RunScriptCmdSub( cmd JtlFrameCommandJson ) bool {

    ad.Println( "script sub command txt = [%s]", cmd.Txt )
	ar.StartSub( cmd.Txt )
	
    return true
}
