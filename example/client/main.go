package main

import (
	"fmt"

	"github.com/AlexBlayE/gourier"
	"github.com/AlexBlayE/gourier/example/client/domain"
)

type ProtocolCommands byte

const (
	_ ProtocolCommands = iota

	VERSION
	TEST
)

type Versions byte

const (
	_ Versions = iota + 10

	V1
	V2
)

type TestSubcommands byte

const (
	_ TestSubcommands = iota + 100

	S1
	S2
)

func main() {
	var ch chan struct{}

	p := gourier.New()
	p.Error(func(ctx *gourier.Context) {
		fmt.Println("eeeee")
	})

	go upHandlers(p)

	payload := []byte{byte(VERSION)}
	err := p.Send("127.0.0.1:3000", payload)

	if err != nil {
		fmt.Println("Mal -> ", err)
	}

	if err != nil {
		fmt.Println("Mal -> ", err)
	}

	<-ch
}

func upHandlers(p *gourier.Server) {
	p.Handler(byte(V1), func(ctx *gourier.Context) {
		ts1 := domain.TestS1{
			Id:      [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			Emmiter: [16]byte{16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		}

		data := domain.SerializeTestS1(ts1)

		fmt.Println("Send data S1 -> ", ts1, "\n", data)
		fmt.Println("------")
		err := ctx.Send(data, byte(TEST), byte(S1))
		if err != nil {
			fmt.Println("Error S1")
			return
		}
	})

	p.Handler(byte(V2), func(ctx *gourier.Context) {
		ts1 := domain.TestS1{
			Id:      [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			Emmiter: [16]byte{16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		}

		ts2 := domain.TestS2{
			TestS1: ts1,
			Count:  uint32(666),
		}

		data := domain.SerializeTestS2(ts2)

		fmt.Println("Send data S2 -> ", ts2, "\n", data)
		// err := ctx.Send(data, byte(TEST), byte(S2))
		err := ctx.Send(data, byte(S2))
		if err != nil {
			fmt.Println("Error s2")
		}

	})

	p.Error(func(ctx *gourier.Context) {
		fmt.Println(ctx.GetPayload())
		// ctx.Send([]byte("BYE_BYE"))
	})

	err := p.Run("127.0.0.1:3001")
	fmt.Println("B")
	if err != nil {
		fmt.Println("Mal")
		fmt.Println(err)
		return
	}
}
