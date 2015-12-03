package main

import (
	"log"
	"net/http"
	"github.com/googollee/go-socket.io"
	"github.com/codeskyblue/go-sh"
)


func main() {

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")

		so.Emit("message",  "ready........!!!!")

		so.On("start", func(msg string) {
			log.Println("execute jmc start")
//			so.BroadcastTo("chat", "chat message", msg)

			so.Emit("message",  "execute jmc start........!!!!")
			sh.Command("../../bin/at_jmc", "01test.jtl").Run()

		})

		so.On("stop", func(msg string) {
			log.Println("execute jmc stop")
			//			so.BroadcastTo("chat", "chat message", msg)

			so.Emit("message",  "execute jmc stop........!!!!")
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
	log.Println("Listen 2.....!!!")

}
