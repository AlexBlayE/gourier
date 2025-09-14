package handlerrunner

import (
	"errors"

	"github.com/AlexBlayE/gourier"
)

type HandlerRunner struct {
	Ctx gourier.Context
}

func (ct *HandlerRunner) RunHandlers(handlers ...gourier.HandleFunc) error {
	for _, handler := range handlers {
		if ct.Ctx.GetAbortFlag() {
			return errors.New("aborted handler")
		}

		handler(ct.Ctx)
	}

	return nil
}
