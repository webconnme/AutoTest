package main

import (
	"time"
	"encoding/json"
)

import (
	zmq "github.com/webconnme/zmq4"
	uuid "github.com/satori/go.uuid"
	"log"
)

const AF_ZMQ_MSG_ADD = "ipc:///tmp/at_msg_add"
const AF_ZMQ_MSG = "ipc:///tmp/at_msg"


type Network struct {
	context *zmq.Context
	pull_socket *zmq.Socket
	req_socket *zmq.Socket
	dealer *zmq.Socket
	router *zmq.Socket
}

var messageMap map[string]MessageStage
var timeoutMap map[string]int

func tickerHandler(data interface{}) error {
	for k, _ := range timeoutMap {
		timeoutMap[k]++
	}
	return nil
}

func resetTimeout(testId string) {
	if _, o := timeoutMap[testId]; o {
		timeoutMap[testId] = 0
	}

}

func requestHandler(sock *zmq.Socket, st zmq.State) error {
	buf, err := sock.RecvBytes(0)
	if err != nil {
		return err
	}

	var req Request
	var msg Message
	err = json.Unmarshal(buf, &req)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, &msg)
	if err != nil {
		return err
	}

	switch req.Command {
	case "new":
		var uid uuid.UUID
		for {
			uid = uuid.NewV4()
			if _, ok := messageMap[uid.String()]; !ok {
				break
			}
		}

		messageMap[uid.String()] = make(MessageStage)
		timeoutMap[uid.String()] = 0

	case "pop":
		if msgMap, ok := messageMap[msg.Test]; ok {
			resetTimeout(msg.Test)
			m := msgMap.PopMessage(msg)
			response, err := json.Marshal(m)
			if err != nil {
				return err
			}

			_, err = sock.SendBytes(response, 0)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func addHandler(sock *zmq.Socket, st zmq.State) error {
	buf, err := sock.RecvBytes(0)
	if err != nil {
		return err
	}

	var msg Message
	err = json.Unmarshal(buf, &msg)
	if err != nil {
		return err
	}

	if msgMap, ok := messageMap[msg.Test]; ok {
		resetTimeout(msg.Test)
		m := msgMap.PopMessage(msg)
		response, err := json.Marshal(m)
		if err != nil {
			return err
		}

		_, err = sock.SendBytes(response, 0)
		if err != nil {
			return err
		}
	}

	return nil
}

func (network *Network) Run() error {
	var err error

	network.context, err = zmq.NewContext()
	if err != nil {
		return err
	}
	defer network.context.Term()

	network.pull_socket, err = network.context.NewSocket(zmq.PULL)
	if err != nil {
		return err
	}
	defer network.pull_socket.Close()
	network.pull_socket.Bind(AF_ZMQ_MSG_ADD)

	network.dealer, err = network.context.NewSocket(zmq.DEALER)
	if err != nil {
		return err
	}
	defer network.dealer.Close()
	network.dealer.Bind("inproc://router")

	network.router, err = network.context.NewSocket(zmq.ROUTER)
	if err != nil {
		return err
	}
	defer network.router.Close()
	network.router.Bind(AF_ZMQ_MSG)

	network.req_socket, err = network.context.NewSocket(zmq.REQ)
	if err != nil {
		return err
	}
	defer network.req_socket.Close()
	network.req_socket.Connect("inproc://router")

	go zmq.Proxy(network.router, network.dealer, nil)

	reactor := zmq.NewReactor()
	reactor.AddSocket(network.pull_socket, zmq.POLLIN, addHandler)
	reactor.AddSocket(network.req_socket, zmq.POLLIN, requestHandler)
	reactor.AddChannelTime(time.Tick(time.Second), 10, tickerHandler)

	err = reactor.Run(time.Second)

	if err != nil {
		log.Panic(err)
	}

	return err
}
