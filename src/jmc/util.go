package main

import (
	"os"
)

//---------------------------------------------------------------------------------------------------------------------
//   
//   패쓰 존재 검사 함수 
//   
//---------------------------------------------------------------------------------------------------------------------
func path_exists( path string ) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return true, err
}

