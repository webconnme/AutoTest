package main

import (
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"ioctl"
	"unsafe"
	"os"
)

type jconfig struct {
	Command	string `json:"command"`
	Data	string `json:"data"`
}

type GpioS struct {
	io int
	mode int
	value int
}

const (
	IOCTL_GPIO_SET_OUTPUT = 7239681
	IOCTL_GPIO_GET_OUTPUT = 7239682
)

var gpioPath string
var file *os.File
var url string

func GPIOOpen() {
	var err error
	file, err = os.OpenFile(gpioPath, os.O_RDWR | os.O_SYNC, 0777)
	if err != nil {
		log.Fatal("open", err)
	}
	log.Println("GPIO Open...");

}

func GPIOClose() {
	file.Close()
	fmt.Println("GPIO Close...")
}

func httpGet(ch chan<- int) {

	var jconf []jconfig

	for {
		client := &http.Client{}

		resp, err := client.Get(url)

		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		contexts, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(">>>contexts : ",string(contexts), " for the url : ", url)
		json.Unmarshal(contexts,&jconf)
		for _, j:=range jconf {
			if j.Command == "power" {
				if j.Data == "on" {
					ch <- 1 //high
				}else if j.Data == "off" {
					ch <- 0    //low
				}
			}
		}
	}
}

func GpioOut(ch <-chan int) {

	g := GpioS {0, 0, 0}
	header := unsafe.Pointer(&g)
	g.io = 28

	for {

		g.value = <-ch
		ioctl.IOCTL(uintptr(file.Fd()), IOCTL_GPIO_SET_OUTPUT, uintptr(header))

		fmt.Println(">>>power stat : ", g.value)
	}
}

func main() {

	channel := make(chan bool)
	txchan := make(chan int)

	url = "Http://nor.kr:3002/v01/relay/80"
	gpioPath = "/dev/ioctl_gpio"

	GPIOOpen()

	go httpGet(txchan)
	go GpioOut(txchan)

	<-channel

	GPIOClose()
}
