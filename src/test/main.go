package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"log"
	"fmt"
)

type JtlFrameScript struct {
	Title     string                   `json:"title"`        // 시험 항목
	Descript  string                   `json:"descript"`     // 시험 설명
	Commands  []interface{}    `json:"task"`         // 시험 태스크
}

func main() {
	c, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Panic("ReadFile:", err)
	}

	var j JtlFrameScript
	json.Unmarshal(c, &j)

	fmt.Println(j)
}

