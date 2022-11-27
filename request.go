package handler

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

const Parameters = "parameters"

type IRequest interface {
	Validation(ctx *fiber.Ctx) *ValidationError
	BeforeValidate(ctx *fiber.Ctx) error
	Validate() ([]*RequestError, error)
	AfterValidate(ctx *fiber.Ctx) error
	GetSchema() interface{}
	BindRequest(ctx *fiber.Ctx)
}

type BaseRequest struct {
	RequestSchema interface{}
}

func (receiver BaseRequest) resolve(ctx *fiber.Ctx) error {
	return ctx.BodyParser(&receiver.RequestSchema)
}

func (receiver BaseRequest) validate(s interface{}, validate *validator.Validate) ([]*RequestError, error) {
	var errors []*RequestError
	err := validate.Struct(s)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return errors, err
		}
		for _, err := range err.(validator.ValidationErrors) {
			var element RequestError
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			element.FieldValue = err.Value()
			errors = append(errors, &element)
		}
	}
	return errors, err
}

func (receiver BaseRequest) BeforeValidate(ctx *fiber.Ctx) error {
	return receiver.resolve(ctx)
}

func (receiver BaseRequest) Validate() ([]*RequestError, error) {
	var validate = validator.New()
	return receiver.validate(receiver.RequestSchema, validate)
}

func (receiver BaseRequest) AfterValidate(ctx *fiber.Ctx) error {
	receiver.BindRequest(ctx)
	return nil
}

func (receiver BaseRequest) BindRequest(ctx *fiber.Ctx) {
	ctx.Locals(Parameters, receiver.GetSchema())
}

func (receiver BaseRequest) GetSchema() interface{} {
	return receiver.RequestSchema
}

func (receiver BaseRequest) Validation(ctx *fiber.Ctx) *ValidationError {
	if eB := receiver.BeforeValidate(ctx); eB != nil {
		return BuildException(fiber.StatusBadRequest, eB)
	}
	if eData, eV := receiver.Validate(); eV != nil {
		return UnprocessableEntityException(eData)
	}
	if eA := receiver.AfterValidate(ctx); eA != nil {
		return BuildException(fiber.StatusInternalServerError, eA)
	}
	return nil
}

type ValidationError struct {
	StatusCode int
	Data       []*RequestError
	Err        error
}

func BuildException(code int, err error, data ...*RequestError) *ValidationError {
	return &ValidationError{
		StatusCode: code,
		Err:        err,
		Data:       data,
	}
}

func UnprocessableEntityException(errorList []*RequestError) *ValidationError {
	return &ValidationError{
		StatusCode: fiber.StatusUnprocessableEntity,
		Err:        errors.New("UnprocessableEntity"),
		Data:       errorList,
	}
}
