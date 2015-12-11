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

	s.Bind("tcp://*:3007")

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
		r.Get("/do1ch/:port", martiniSender)
		r.Post("/do1ch/:port", martiniReceiver)
	})

	m.RunOnAddr(":3006")
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
		log.Println(string(buf))
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
		case "do":
			return true
		default:
			return false
		}
	}

	return false
}
