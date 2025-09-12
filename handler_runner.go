package gourier

import "errors"

type handlerRunner struct { // TODO: utilitzar interface que implementi RunHandler
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
