package main
import "fmt"


func main() {

	msg := Message{
		Id: "myId",
		Target: "yourId1",
		Stage: 0,
		Data: "",
	}

	PushMessage(Message{
		Id: "myId",
		Target: "yourId1",
		Stage: 1,
		Data: "data1",
	})
	PushMessage(Message{
		Id: "myId",
		Target: "yourId1",
		Stage: 2,
		Data: "data2",
	})
	PushMessage(Message{
		Id: "myId",
		Target: "yourId2",
		Stage: 2,
		Data: "data2",
	})
	PushMessage(Message{
		Id: "myId",
		Target: "yourId1",
		Stage: 3,
		Data: "data3",
	})


	result := PopMessage(msg)
	fmt.Println(result)

	result = PopMessage(msg)
	fmt.Println(result)

	result = PopMessage(msg)
	fmt.Println(result)

	network := new(Network)
	network.Run()
}