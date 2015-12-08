package main


import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/mikepb/go-serial"
	"log"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"time"
)

type jconfig struct {
	Command	string `json:"command"`
	Data	string `json:"data"`
}


var RS232options = serial.Options{

	BitRate	 	: 115200,
	DataBits 	: 8,
	Parity	 	: serial.PARITY_NONE,
	StopBits 	: 1,
	FlowControl : serial.FLOWCONTROL_NONE,
	Mode		: serial.MODE_READ_WRITE,
}

var RS232Path string
var serialPort *serial.Port

var chanmapchan = make(chan map[string]chan []jconfig, 1)

func RS232Open() {

	options := RS232options

	var err error
	if serialPort != nil {
		serialPort.Close()
	}

	serialPort, err = options.Open(RS232Path)
	if err != nil {
		log.Fatal("serial open ",err)
	} else {
		fmt.Println("serial open...")
	}
}

func RS232Close() {
	if serialPort != nil {
		serialPort.Close()
		fmt.Println("serial close...")
	}
}

func runMartini(sending, receiving chan []jconfig) {
	m := martini.Classic()

	chanmap := make(map[string]chan []jconfig)
	chanmap["GET"] = sending
	chanmap["POST"] = receiving
	chanmapchan <- chanmap

	m.Group("/v01", func(r martini.Router) {
		r.Get("/rs232/:port", martiniSender)
		r.Post("/rs232/:port", martiniReceiver)
	})

	m.RunOnAddr(":3006")
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

func RS232Rx(ch chan []jconfig) {

	for {
		remain, err := serialPort.InputWaiting()
		if err != nil {
			fmt.Println(err)
			continue
		}
		if remain != 0 {

			buf := make([]byte, remain)
			_, err = serialPort.Read(buf)
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Println(">>>rs232 read : ", string(buf))

			jconf := jconfig{
				Command:"rx",
				Data:string(buf),
			}
			m := []jconfig{jconf}

			select {
			case toDeliver := <-ch:
				toDeliver = append(toDeliver,m...)
				ch <- toDeliver
			default:
				ch <- m
			}
		}

	}
}

func RS232Tx(ch chan []jconfig) {

	for {
		fmt.Println("waiting channel")
		buf := <-ch

		fmt.Println(">>>Send String: ",buf)
		for _, tmp := range buf {
			serialPort.Write([]byte(tmp.Data))
		}
	}
}


func main() {

	s2m := make(chan []jconfig, 1)
	m2s := make(chan []jconfig, 1)


	RS232Path = "/dev/ttyS1"
	RS232Open()

	go runMartini(s2m,m2s)
	go RS232Rx(s2m)	//get
	go RS232Tx(m2s) //post

	fmt.Scanln()

	RS232Close()

}