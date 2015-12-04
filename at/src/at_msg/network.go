package main

import (
	zmq "github.com/alecthomas/gozmq"
	"time"
	"encoding/json"
)

const AF_ZMQ_MSG_PULL = "ipc:///tmp/at_msg_pull"
const AF_ZMQ_MSG_REQ = "ipc:///tmp/at_msg_req"


type NetworkOnQuery    func( network *Network, msg Message)
type NetworkOnAdd    func( network *Network, msg Message)

type Network struct {
	context *zmq.Context
	pull_socket *zmq.Socket
	req_socket *zmq.Socket

	OnQuery NetworkOnQuery
	OnAdd NetworkOnAdd
}

func (network *Network) Run() error {
	var err error

	network.context, err = zmq.NewContext()
	if err != nil {
		return err
	}
	defer network.context.Close()

	network.pull_socket, err = network.context.NewSocket(zmq.PULL)
	if err != nil {
		return err
	}
	defer network.pull_socket.Close()

	network.req_socket, err = network.context.NewSocket(zmq.REQ)
	if err != nil {
		return err
	}
	defer network.req_socket.Close()

	poll_items := []zmq.PollItem {
		zmq.PollItem{ Socket: network.pull_socket, Events: zmq.POLLIN},
		zmq.PollItem{ Socket: network.req_socket, Events: zmq.POLLIN},
	}

	for {
		event_count, err := zmq.Poll(poll_items, 1 * time.Second )
		if err != nil {
			return err
		}

		var msg Message

		if event_count > 0 {
			if poll_items[0].REvents & zmq.POLLIN != 0 {
				buf, err := poll_items[0].Socket.Recv(0)
				if err != nil {
					// handle error
				} else {
					err := json.Unmarshal(buf, &msg)
					if err != nil {
						// handle error
					} else {
						if network.OnQuery != nil {
							network.OnQuery(network, msg)
						}
					}
				}
			}

			if poll_items[1].REvents & zmq.POLLIN != 0 {
				buf, err := poll_items[1].Socket.Recv(0)
				if err != nil {
					// handle error
				} else {
					err := json.Unmarshal(buf, &msg)
					if err != nil {
						// handle error
					} else {
						if network.OnAdd != nil {
							network.OnAdd(network, msg)
						}
					}
				}
			}
		}
	}

	return nil
}

func SendResponse(sock *zmq.Socket, msg Message) {
	buf, err := json.Marshal(msg)
	if err != nil {

	} else {
		sock.Send(buf, 0)
	}
}
