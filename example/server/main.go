package main

import (
	"fmt"

	"github.com/AlexBlayE/gourier"
	"github.com/AlexBlayE/gourier/example/server/controllers"
	"github.com/AlexBlayE/gourier/example/server/middlewares"
	connectionmanager "github.com/AlexBlayE/gourier/internal/connection_manager"
	pathfinder "github.com/AlexBlayE/gourier/internal/path_finder"
	group "github.com/AlexBlayE/gourier/internal/router_group"
)

type ProtocolCommands = byte

const (
	_ ProtocolCommands = iota

	VERSION
	TEST
)

type Versions = byte

const (
	_ Versions = iota + 10

	V1
	V2
)

type TestSubcommands = byte

const (
	_ TestSubcommands = iota + 100

	S1
	S2
)

func main() {
	radix := &pathfinder.RadixNode{make(map[byte]gourier.PathFinder), nil, nil, 0}

	p := gourier.New(
		connectionmanager.NewConnectionManager(30, 1024, radix),
		&group.RouterGroup{radix},
		make(chan struct{}, 100),
	)

	p.Handler(VERSION, func(ctx gourier.Context) {
		fmt.Println("Response version")
		ctx.Send(nil, byte(V1))
	})

	p.Error(func(ctx gourier.Context) {
		fmt.Println("Error version")
	})

	g := p.Group(TEST)
	g.Handler(
		S1,
		middlewares.DeserializeTestS1Middleware,
		controllers.TestS1Controller,
	)

	g.Handler(
		S2,
		middlewares.DeserializeTestS2Middleware,
		controllers.TestS2Controller,
	)

	g.Error(func(ctx gourier.Context) {
		fmt.Println("Command error")
	})

	err := p.Run(":3000")
	if err != nil {
		fmt.Println("Initiation serveer error")
	}

}
