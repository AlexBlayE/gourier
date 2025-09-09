package gourier

type handlerRunner struct {
	ctx *Context
}

func (ct *handlerRunner) RunHandlers(handlers ...HandleFunc) {
	for _, handler := range handlers {
		// TODO: implementar mecanisme abort
		handler(ct.ctx)
	}
}
