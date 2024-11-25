package messaging

import (
	"github.com/Dert091499/Utilities/context"
	"github.com/Dert091499/Utilities/logs"
)

const (
	Handle DispatchType = "handle"
	Error  DispatchType = "error"
)

type (
	MiddlewareHandlerFunc func(ctx *context.Context, dto DispatchDTO) error
	MiddlewareFunc        func(MiddlewareHandlerFunc) MiddlewareHandlerFunc

	HandlerFunc      func(ctx *context.Context, msg Message) error
	ErrorHandlerFunc func(ctx *context.Context, msg Message, err error)

	Job interface {
		Process(ctx *context.Context, msg Message) error
		OnError(ctx *context.Context, msg Message, err error)
	}

	DispatchType string
	DispatchDTO  struct {
		Type      DispatchType
		Source    string
		RequestID string
		MsgType   string
		Msg       Message
		Log       logs.Logger
		Err       error
	}
	Dispatcher interface {
		AddHandler(handler HandlerFunc, errorHandler ErrorHandlerFunc, msgType ...string)
		Use(middlewareFunc ...MiddlewareFunc)
		Dispatch(dto DispatchDTO) error
	}
)

func applyMiddleware(h MiddlewareHandlerFunc, middleware ...MiddlewareFunc) MiddlewareHandlerFunc {
	for i := 0; i < len(middleware); i++ {
		h = middleware[i](h)
	}
	return h
}
