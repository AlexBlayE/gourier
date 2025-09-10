package gourier

import "errors"

type handlerRunner struct {
	ctx *Context
}

func (ct *handlerRunner) RunHandlers(handlers ...HandleFunc) error {
	for _, handler := range handlers {
		if ct.ctx.abortFlag {
			return errors.New("aborted handler")
		}

		handler(ct.ctx)
	}

	return nil
}
