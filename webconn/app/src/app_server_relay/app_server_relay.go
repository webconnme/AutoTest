/**
 * The MIT License (MIT)
 *
 * Copyright (c) 2015 Jane Lee <jane@webconn.me>
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

package main


import (
	"fmt"
	"github.com/go-martini/martini"
	"log"
	"net/http"
	"encoding/json"
	"io/ioutil"
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

var chanmapchan = make(chan map[string]chan []jconfig, 1)

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

func runMartini(receiving chan []jconfig) {
	m := martini.Classic()

	chanmap := make(map[string]chan []jconfig)
//	chanmap["GET"] = sending
	chanmap["POST"] = receiving
	chanmapchan <- chanmap

	m.Group("/v01", func(r martini.Router) {
//		r.Get("/relay/:port", martiniSender)
		r.Post("/relay/:port", martiniReceiver)
	})

	m.RunOnAddr(":3007")
}

func martiniReceiver(r *http.Request, params martini.Params) string {
	if r.Body == nil {
		log.Println("http.Request.Body is empty")
		return ""
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("ioutil.ReadAll():", err)
		return ""
	}

	fmt.Println(">>>POST:", string(data))

	var u []jconfig

	err = json.Unmarshal(data, &u)
	if err != nil {
		log.Println("json.Unmarshal():", err)
		return ""
	}

	if !check(u) {
		log.Println("Invalid Command", u)
		return ""
	}

	m := <-chanmapchan
	c := m["POST"]
	chanmapchan <- m

	select {
	case toDeliver := <-c:
		fmt.Println(">>>append channel")
		toDeliver = append(toDeliver, u...)
		c <- toDeliver
	default:
		fmt.Println(">>>put channel")
		c <- u
	}

	log.Println("OK")

	return "OK"
}

func check(array []jconfig) bool {
	for _, u := range array {
		cmd := u.Command

		switch cmd {
		case "power":
			return true
		default:
			return false
		}
	}

	return false
}

func GpioOut(ch chan []jconfig) {

	g := GpioS {0, 0, 0}
	header := unsafe.Pointer(&g)
	g.io = 28

	for {

		buf := <-ch

		fmt.Println(">>>Send String: ",buf)
		for _, tmp := range buf {
			if tmp.Data == "on" {
				g.value = 1
			} else if tmp.Data == "off" {
				g.value = 0
			}
			ioctl.IOCTL(uintptr(file.Fd()), IOCTL_GPIO_SET_OUTPUT, uintptr(header))

		}

		fmt.Println(">>>power stat : ", g.value)
	}
}

func main() {

//	s2m := make(chan []jconfig, 1)
	m2g := make(chan []jconfig, 1)

	gpioPath = "/dev/ioctl_gpio"
	GPIOOpen()

	go runMartini(m2g)
	go GpioOut(m2g)

	fmt.Scanln()

	GPIOClose()

}