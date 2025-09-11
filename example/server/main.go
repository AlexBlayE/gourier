package main

import (
	"fmt"
	"gourier"
	"gourier/example/server/controllers"
	"gourier/example/server/middlewares"
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
	p := gourier.New()
	// p.SetOptions()

	p.Handler(byte(VERSION), func(ctx *gourier.Context) {
		fmt.Println("Response version")
		ctx.Send(nil, byte(V1))
	})

	p.Error(func(ctx *gourier.Context) {
		fmt.Println("Error version")
		// ctx.Send([]byte("VERSION_ERROR"))
	})

	g := p.Group(byte(TEST))
	g.Handler(
		byte(S1),
		middlewares.DeserializeTestS1Middleware,
		controllers.TestS1Controller,
	)

	g.Handler(
		byte(S2),
		middlewares.DeserializeTestS2Middleware,
		controllers.TestS2Controller,
	)

	g.Error(func(ctx *gourier.Context) {
		fmt.Println("Command error")
	})

	err := p.Run(":3000")
	if err != nil {
		fmt.Println("Initiation serveer error")
	}

}
