package gohandler

type RequestError struct {
	FailedField string
	Tag         string
	Value       string
	FieldValue  interface{}
}
