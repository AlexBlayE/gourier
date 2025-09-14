package middlewares

import (
	"fmt"

	"github.com/AlexBlayE/gourier"
	"github.com/AlexBlayE/gourier/example/server/domain"
)

func DeserializeTestS2Middleware(ctx gourier.Context) {
	payload := ctx.GetPayload()

	t2, err := domain.DeserializeTestS2(payload)
	if err != nil {
		fmt.Println("Deserialization error")
		ctx.Abort(nil)
	}

	ctx.Set("tests2", t2)
}
