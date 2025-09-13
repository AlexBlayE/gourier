package middlewares

import (
	"fmt"

	"github.com/AlexBlayE/gourier"
	"github.com/AlexBlayE/gourier/example/server/domain"
)

func DeserializeTestS1Middleware(ctx *gourier.Context) {
	payload := ctx.GetPayload()

	t1, err := domain.DeserializeTestS1(payload)
	if err != nil {
		fmt.Println("Deserialization error")
		ctx.Abort(nil)
	}

	ctx.Set("tests1", t1)
}
