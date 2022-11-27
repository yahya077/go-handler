package gohandler

import (
	"github.com/gofiber/fiber/v2"
	"reflect"
)

type IBaseHandler interface {
	Handle(service interface{}, method string) error
	SetRequest(request IRequest) BaseHandler
	GetValueOfC() []reflect.Value
	CallServiceMethod(service interface{}, method string) error
	WithHandler(handler IHandler) IBaseHandler
}

type IHandler interface {
	LocalBinding(ctx *fiber.Ctx) error
}

type BaseHandler struct {
	Request IRequest
	C       *fiber.Ctx
	Handler IHandler
}

type Handler struct {
}

func (h BaseHandler) SetRequest(request IRequest) BaseHandler {
	h.Request = request
	return h
}

func New(ctx *fiber.Ctx) IBaseHandler {
	return BaseHandler{
		C: ctx,
	}
}

func (h BaseHandler) Handle(service interface{}, method string) error {
	if h.Request != nil {
		if validation := h.Request.Validation(h.C); validation != nil {
			return h.C.Status(validation.StatusCode).JSON(validation.Data)
		}
	}

	if e := h.Handler.LocalBinding(h.C); e != nil {
		return h.C.SendStatus(fiber.StatusNotFound)
	}

	return h.CallServiceMethod(service, method)
}

func (h BaseHandler) CallServiceMethod(service interface{}, method string) error {
	result := reflect.ValueOf(service).MethodByName(method).Call(h.GetValueOfC())
	if _, ok := result[0].Interface().(error); ok {
		return result[0].Interface().(error)
	}
	return nil
}

func (h BaseHandler) GetValueOfC() []reflect.Value {
	return []reflect.Value{reflect.ValueOf(h.C)}
}

func (h BaseHandler) WithHandler(handler IHandler) IBaseHandler {
	h.Handler = handler
	return h
}

func (h Handler) LocalBinding(ctx *fiber.Ctx) error {
	return nil
}
