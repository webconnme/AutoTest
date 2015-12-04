package main
import "fmt"

type Message struct {
	Id string `json:"id"`
	Target string `json:"target"`
	Stage int `json:"stage"`
	Data string `json:"data"`
}

type MessageArray []Message
type MessageMap map[string]MessageArray
type MessageStage map[int]MessageMap


var messages MessageStage

func PopMessage(msg Message) *Message {
	index := -1

	for k, _ := range messages {
		if index == -1 {
			index = k
		} else {
			if index > k {
				index = k
			}
		}
	}

	if index == -1 {
		return nil
	}

	if _, ok := messages[index][msg.Target]; !ok {
		return nil
	}

	if len(messages[index][msg.Target]) <= 0 {
		return nil
	}

	result := messages[index][msg.Target][0]

	messages[index][msg.Target] = messages[index][msg.Target][1:]

	if len(messages[index][msg.Target]) == 0 {
		delete(messages[index], msg.Target)
	}

	fmt.Println("-------------")
	fmt.Println(messages)
	fmt.Println("-------------")

	if len(messages[index]) == 0 {
		delete(messages, index)
	}
	return &result
}

func PushMessage(msg Message) {
	if messages == nil {
		messages = make(MessageStage)
	}

	if _, ok := messages[msg.Stage]; !ok {
		messages[msg.Stage] = make(map[string]MessageArray)
	}

	if _, ok := messages[msg.Stage][msg.Target]; !ok {
		messages[msg.Stage][msg.Target] = make(MessageArray, 0)
	}

	messages[msg.Stage][msg.Target] = append(messages[msg.Stage][msg.Target], msg)


	fmt.Println(messages)
}