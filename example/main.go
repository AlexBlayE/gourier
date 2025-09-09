package main

import (
	"fmt"
	"gourier"
)

// TODO: Fer que a part del send també pugui enviar peticions a altres ip
// TODO: fer que si no coincideix res donar la opció de resposta d'error
// TODO: he fet tota la part que rep. Hara fer la part que envia(aunque crec que hamb el send ya hi ha suficient)
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
	g.Handler(byte(V1), func(ctx *gourier.Context) {

	})
	g.Handler(byte(V2), func(ctx *gourier.Context) {
		// ctx.Send() // Per enviar resposta o a altres ip
	})

	g5 := g.Group(5)
	g5.Handler(6, func(ctx *gourier.Context) {})

	err := p.Run(":3000")
	if err != nil {
		fmt.Println("Initiation serveer error")
	}
}
