package handler

type RequestError struct {
	FailedField string
	Tag         string
	Value       string
	FieldValue  interface{}
}
