package main

import (
	"log"
	"net/http"
	"github.com/googollee/go-socket.io"
	"github.com/codeskyblue/go-sh"
	"net"
//	"time"
)

var client net.Conn

func main() {


	chs := make([]socketio.Socket, 0)

	tobuf := make(chan []byte)

	fromSocket := make(chan socketio.Socket)

	server, err := socketio.NewServer(nil)

	if err != nil {
		log.Fatal(err)
	}

	// client
	client, err := net.Dial("tcp", "127.0.0.1:2999")


	if err != nil {
		log.Fatal(err)
	}

	log.Println("client connect!!!!!!!!")

	go func(c net.Conn) {

		data := make([]byte, 4096)

		for {
			n, err := client.Read(data)

			if err != nil {
				log.Println(err)
				return
			}

			tobuf <- data[:n]

			log.Print(string(data[:n]))
			// 브라우저로 보냄
			//so.Emit("message", string(data[:n]))

			//					time.Sleep(1 * time.Second)
		}
	}(client)

	go func() {

		for {
			select {
			case databuf := <-tobuf:

				for _, ch := range chs {
					ch.Emit("message", string(databuf))

				}
			case ch := <-fromSocket:

				chs = append(chs, ch)

			}
		}

	}()


	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")


		fromSocket <- so

		so.Emit("message",  "ready........!!!!\n")

		so.On("start", func(msg string) {
			log.Println("execute jmc start")
//			so.BroadcastTo("chat", "chat message", msg)

			so.Emit("message",  "execute jmc start........!!!!\n")
			sh.Command("../../bin/at_jmc", "relay_test.jtl").Run()

	//			defer client.Close()


		})

		so.On("stop", func(msg string) {
			log.Println("execute jmc stop")
			//			so.BroadcastTo("chat", "chat message", msg)

			so.Emit("message",  "execute jmc stop........!!!!\n")
		})


		so.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})

	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	http.Handle("/socket.io/", server)

	// 고정 HTML파일을 볼수 있도록 파일 서버를 설정 한다.
	http.Handle("/", http.FileServer(http.Dir(".")))

	// private views
	http.HandleFunc("/post", PostOnly(BasicAuth(HandlePost)))
	http.HandleFunc("/json", GetOnly(BasicAuth(HandleJSON)))


	log.Println("Listen.....!!!")
	log.Fatal(http.ListenAndServe(":9000", nil))

}
