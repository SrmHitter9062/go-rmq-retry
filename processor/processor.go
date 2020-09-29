package processor

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Message struct {
	Years     int    `json:"years"`
	Profile   string `json:"profile"`
	MessageID int    `json:message_id`
}

/**
Message received {

  "years" : 5,
  "profile": "SDE",
  "messageID": 1
}
**/

func Process(data []byte) error {
	var msg *Message
	err := json.Unmarshal(data, &msg)
	fmt.Println("processing the message", msg)
	if err != nil {
		fmt.Println("Error in decoding:", err)
		return err
	}
	if msg.Years >= 5 {
		return errors.New("Need to be requeue")
	}
	return nil
}
