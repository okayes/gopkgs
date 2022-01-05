package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var robotMessageCh chan string
var robotMessageUrl string
var robotMessageTitle string

type Option struct {
	MsgChCapacity int
	MsgUrl        string
	Title         string
}

type robotMessage struct {
	MsgType string           `json:"msgtype"`
	Text    robotMessageText `json:"text"`
}

type robotMessageText struct {
	Content string `json:"content"`
}

func Init(option Option) {
	robotMessageCh = make(chan string, option.MsgChCapacity)
	robotMessageUrl = option.MsgUrl
	robotMessageTitle = option.Title

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("wechat robot message recover:", r)
			}
		}()

		for {
			select {
			case m, ok := <-robotMessageCh:
				if ok {
					doSend(m)
				} else {
					log.Println("wechat robot message channel close")
					return
				}
			}
		}
	}()
}

func Send(message string) {
	if len(robotMessageCh) >= cap(robotMessageCh) {
		return
	}
	robotMessageCh <- message
}

func doSend(message string) {
	message = fmt.Sprintf("%s\n%s\n%s", robotMessageTitle, time.Now().Format("2006-01-02 15:04:05"), message)
	robotMessage := robotMessage{
		MsgType: "text",
		Text: robotMessageText{
			Content: message,
		},
	}
	data, err := json.Marshal(robotMessage)
	if err != nil {
		log.Println(err)
		return
	}

	resp, err := http.Post(robotMessageUrl, "application/json", bytes.NewReader(data))
	if err != nil {
		log.Println(err)
		return
	}

	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(data))
}
