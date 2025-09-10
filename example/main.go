package main

import (
	"fmt"
	"gourier"
)

// TODO: net.Dial("tcp", "host:port") -> per enviar
// TODO: tests

type ProtocolCommands byte

const (
	_ ProtocolCommands = iota

	VERSION
	TEST
)

type TestSubcommands byte

const (
	_ TestSubcommands = iota

	V1
	V2
)

func main() {
	p := gourier.New()

	p.Handler(byte(VERSION), func(ctx *gourier.Context) {
		// ctx.Next()
	})

	g := p.Group(byte(TEST))
	g.Handler(byte(V1),
		func(ctx *gourier.Context) {
			fmt.Println("EUREKA")
			// ctx.Abort([]byte{})
		},
		func(ctx *gourier.Context) {
			fmt.Println("HOUSTON")
			// fmt.Println(ctx.GetPayload())
			err := ctx.Send([]byte{1, 2, 3, 4, 5, 6, 7}, 3, 2)
			fmt.Println("Error -> ", err)
		},
	)

	// g.Handler(byte(V2), func(ctx *gourier.Context) {
	// 	// ctx.Send() // Per enviar resposta o a altres ip
	// })

	// g5 := g.Group(5)
	// g5.Handler(6, func(ctx *gourier.Context) {})

	err := p.Run(":3000")
	if err != nil {
		fmt.Println("Initiation serveer error")
	}

}
