# lineGo

**Golang lib for line unofficial API**

## Important

Please be warned: lineGo is in a very early beta. You will encounter bugs when using it. In addition, due to the specification of the thrift library, go func cannot be used.

## Installation

`go get github.com/n4tsumi/lineGo`

## Quick start

```go
package main

import (
	line "github.com/n4tsumi/lineGo"
	talk "github.com/n4tsumi/lineGo/talkservice"
	"time"
)

var cl = line.NewClient()

func main() {
	cl.Login(line.Token("Your Auth Token"))
	cl.AddOpInterrupt(talk.Af_RECEIVE_MESSAGE, chatBot)

	cl.Run()
}


func chatBot(op *talk.Operation) {
	msg := op.Message

	switch {
	case msg.Text == "hello":
		cl.SendText(x.To, "Hello!")

	case msg.Text == "time":
		cl.SendText(x.To, time.Now().Format("2006/01/02 15:04:05"))
	}
}
```