package main

type Message struct {
	Test string `json:"test"`
	Id string `json:"id"`
	Target string `json:"target"`
	Stage int `json:"stage"`
	Data string `json:"data"`
}

type Request struct {
	Message
	Command string `json:"command"`
}

type MessageArray []Message
type MessageMap map[string]MessageArray
type MessageStage map[int]MessageMap

func (messages MessageStage) GetStage(min int) int {
	index := -1

	for k, _ := range messages {
		if min < k {
			if index == -1 {
				index = k
			} else {
				if index > k {
					index = k
				}
			}
		}
	}

	return index
}

func (messages MessageStage) popMessage(msg Message, minStage int) *Message {
	index := messages.GetStage(minStage)

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

	if len(messages[index]) == 0 {
		delete(messages, index)
	}
	return &result
}

func (messages MessageStage) PopMessage(msg Message) *Message {
	m := messages.popMessage(msg, -1)
	if m != nil {
		return m
	}
	return messages.popMessage(msg, 0)
}

func (messages MessageStage) PushMessage(msg Message) {
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
}