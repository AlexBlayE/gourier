package controllers

import (
	"fmt"
	"gourier"
	"gourier/example/server/domain"
)

func TestS2Controller(ctx *gourier.Context) {
	elem := ctx.Get("tests2")
	t2 := elem.(domain.TestS2)

	fmt.Println("Ok test s2 -> ", t2)

	ctx.Send([]byte("OK-S2"))
}
