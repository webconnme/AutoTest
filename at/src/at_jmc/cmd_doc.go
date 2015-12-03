package main

//   주석 명령을 처리 한다. 
func (jtl *JtlFrame) RunScriptCmdDoc( cmd JtlFrameCommandJson ) bool {

    ad.Println( "script doc command txt = [%s]", cmd.Txt )
	ar.WriteDocument( cmd.Txt );

    return true
}
