package response

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string   `json:"status"`
	Errors []string `json:"errors,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(errors ...string) Response {
	return Response{
		Status: StatusError,
		Errors: errors,
	}
}

func ValidationError(err validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range err {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return Error(errMsgs...)
}
