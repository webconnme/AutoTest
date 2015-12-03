package at


// AT 프레임 주기적 호출 
func onPeriod( af *AtFrame )(bool){
    return false
}

func onKill( af *AtFrame )(bool){
    return true
}

// AT 프레임 리셋 
func onReset( af *AtFrame )(bool){
    return false
}

// AT 프레임간의 연결
func onLink( af *AtFrame , data interface{} )(bool){
    return false
}

// AT 프레임간의 연결 해제
func onUnlink( af *AtFrame , data interface{} )(bool){
    return false
}

// AT 프레임 동작 시작
func onStart( af *AtFrame )(bool){
    return false
}

// AT 프레임 동작 강제 중지
func onStop( af *AtFrame )(bool){
    return false
}

// AT 프레임 설정
func onSet( af *AtFrame, data interface{} )(bool){
    return false
}

// AT 프레임 값얻기 
func onCall( af *AtFrame, data interface{} )(bool){
    return false
}

// AT 프레임 값 반환
func onAck( af *AtFrame )(bool){
    return false
}

// ZMQ 수신 이벤트 
func onRxIn( af *AtFrame, index int, data []byte )(bool){
    return false
}
