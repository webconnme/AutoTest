package main

//  서브 시험 항목 시작
func (jtl *JtlFrame) RunScriptCmdSubEnd( cmd JtlFrameCommandJson ) bool {

    ad.Println( "script subend command" )
	ar.EndSub()
	
    return true
}
