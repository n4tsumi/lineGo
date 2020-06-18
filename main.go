package main

import (
	"fmt"
	"strconv"
	"time"

	. "./line"
	talk "./talkservice"
)

var client = NewClient("token", "IOS")

func main() {
	operation()
}

func operation() {
	for {
		ops, err := client.FetchOperations()
		for _, op := range ops {
			if op.Type != 0 {
				if op.Revision > client.Revision {
					client.Revision = op.Revision
				}
				bot(op)
			}
		}
		if err != nil {
			fmt.Println(err)
		}
	}
}

func bot(op *talk.Operation) {
	switch op.Type {
	case 26:
		msg := op.Message
		switch msg.Text {
		case ".speed":
			start := time.Now()
			client.SendMessage(msg.To, "Start")
			end := time.Now()
			t := strconv.FormatFloat(end.Sub(start).Seconds(), 'f', 8, 64)
			client.SendMessage(msg.To, t)
		case "^^":
			client.SendMessage(msg.To, "^^")
		case ".test":
			client.SendMessage(msg.To, "Hello Go")
		case ".mid":
			client.SendMessage(msg.To, client.Profile.Mid)
		}
	}
}
