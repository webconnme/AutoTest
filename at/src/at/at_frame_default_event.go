/**
 * The MIT License (MIT)
 *
 * Copyright (c) 2015 David You <david@webconn.me>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 * 
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 * 
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

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
