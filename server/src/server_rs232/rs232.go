/**
 * The MIT License (MIT)
 *
 * Copyright (c) 2015 Victor Kim <victor@webconn.me>
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
	"encoding/json"
	"net/http"
	"fmt"
)

import (
	zmq "github.com/alecthomas/gozmq"
	"github.com/go-martini/martini"
	"log"
	"time"
	"sync"
	"io/ioutil"
)

var chanmapchan = make(chan map[string]chan []map[string]interface{}, 1)

func main() {
	z2m := make(chan []map[string]interface{}, 1)
	m2z := make(chan []map[string]interface{}, 1)

	go runZMQ(m2z, z2m)
	runMartini(z2m, m2z)
}

func runZMQ(sending, receiving chan []map[string]interface{}) {
	c, err := zmq.NewContext()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	s, err := c.NewSocket(zmq.PAIR)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	s.Bind("tcp://*:3001")

	wg := sync.WaitGroup{}

	wg.Add(2)
	go zmqSender(s, sending, &wg)
	go zmqReceiver(s, receiving, &wg)

	wg.Wait()
}

func runMartini(sending, receiving chan []map[string]interface{}) {
	m := martini.Classic()

	chanmap := make(map[string]chan []map[string]interface{})
	chanmap["GET"] = sending
	chanmap["POST"] = receiving
	chanmapchan <- chanmap

	m.Group("/v01", func(r martini.Router) {
		r.Get("/rs232/:port", martiniSender)
		r.Post("/rs232/:port", martiniReceiver)
	})

	m.Run()
}

func zmqSender(s *zmq.Socket, c chan []map[string]interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		fmt.Println("waiting channel")
		buf := <-c
		data, err := json.Marshal(buf)
		if err != nil {
			log.Println("json.Marshal():", err)
			continue
		}

		fmt.Println("Send String: " + string(data))
		s.Send(data, zmq.NOBLOCK)
	}
}

func zmqReceiver(s *zmq.Socket, c chan []map[string]interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		buf, err := s.Recv(0)
		if err != nil {
			log.Println("zmq.Socket.Recv():", err)
			continue
		}

		var m []map[string]interface{}
		err = json.Unmarshal(buf, &m)
		if err != nil {
			log.Println("json.Unmarshal():", err)
			continue
		}

		if !check(m) {
			log.Println("Invalid Command:", m)
			continue
		}

		select {
		case toDeliver := <-c:
			toDeliver = append(toDeliver, m...)
			c <- toDeliver
		default:
			c <- m
		}
	}
}

func martiniSender() string {
	m := <-chanmapchan
	c := m["GET"]
	chanmapchan <- m
	select {
	case buf := <-c:
		fmt.Printf("martini send: %v", buf)
		s, err := json.Marshal(buf)
		if err != nil {
			log.Println(err)
			return ""
		}
		return string(s)
	case <-time.After(time.Second * 60):
		return "Timed Out"
	}
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

	fmt.Println("POST:", string(data))

	var u []map[string]interface{}

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
		fmt.Println("append channel")
		toDeliver = append(toDeliver, u...)
		c <- toDeliver
	default:
		fmt.Println("put channel")
		c <- u
	}

	log.Println("OK")

	return "OK"
}

func check(array []map[string]interface{}) bool {
	for _, u := range array {
		cmd, ok := u["command"]
		if !ok {
			return false
		}

		switch cmd {
		case "rx":
			return true
		case "tx":
			return true
		case "rs232_option":
			return true
		default:
			return false
		}
	}

	return false
}
