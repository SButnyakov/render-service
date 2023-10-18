package response

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

type Response struct {
	Error string `json:"error,omitempty"`
}

func Error(msg string) Response {
	return Response{
		Error: msg,
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "email":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid email", err.Field()))
		case "min", "max":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is out of its range", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return Response{
		Error: strings.Join(errMsgs, ", "),
	}
}
