package controllers

import (
	"fmt"

	"github.com/AlexBlayE/gourier"
	"github.com/AlexBlayE/gourier/example/server/domain"
)

func TestS1Controller(ctx gourier.Context) {
	elem := ctx.Get("tests1")
	t1 := elem.(domain.TestS1)

	fmt.Println("Ok test s1 -> ", t1)

	// ctx.Send([]byte("OK-S1"))
}
