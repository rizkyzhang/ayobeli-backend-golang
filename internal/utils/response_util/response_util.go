package response_util

import (
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ValidationError struct {
	Field string `json:"field"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Response struct {
	Code             int               `json:"code"`
	Status           string            `json:"status"`
	Data             interface{}       `json:"data,omitempty"`
	Error            string            `json:"error,omitempty"`
	BindingError     string            `json:"binding_error,omitempty"`
	ValidationErrors []ValidationError `json:"validation_errors,omitempty"`
}

func FromOK() *Response {
	return &Response{
		Status: http.StatusText(http.StatusOK),
		Code:   http.StatusOK,
	}
}

func FromCreated() *Response {
	return &Response{
		Status: http.StatusText(http.StatusCreated),
		Code:   http.StatusCreated,
	}
}

func FromData(data interface{}) *Response {
	return &Response{
		Status: http.StatusText(http.StatusOK),
		Code:   http.StatusOK,
		Data:   data,
	}
}

func FromError(err error) *Response {
	if strings.Contains(strings.ToLower(err.Error()), "not found") {
		return &Response{
			Status: http.StatusText(http.StatusNotFound),
			Code:   http.StatusNotFound,
			Error:  err.Error(),
		}
	}

	return &Response{
		Status: http.StatusText(http.StatusInternalServerError),
		Code:   http.StatusInternalServerError,
		Error:  err.Error(),
	}
}

func FromBindingError(err error) *Response {
	return &Response{
		Status:       http.StatusText(http.StatusBadRequest),
		Code:         http.StatusBadRequest,
		BindingError: err.Error(),
	}
}

func FromValidationError(err ValidationError) *Response {
	return &Response{
		Status: http.StatusText(http.StatusBadRequest),
		Code:   http.StatusBadRequest,
		ValidationErrors: []ValidationError{
			err,
		},
	}
}

func FromValidationErrors(_errs validator.ValidationErrors) *Response {
	var errs []ValidationError

	for _, err := range _errs {
		errs = append(errs, ValidationError{
			Field: strings.ToLower(err.StructField()),
			Name:  err.Tag(),
			Value: err.Param(),
		})
	}

	return &Response{
		Status:           http.StatusText(http.StatusBadRequest),
		Code:             http.StatusBadRequest,
		ValidationErrors: errs,
	}
}

func FromBadRequestError(err error) *Response {
	return &Response{
		Status: http.StatusText(http.StatusBadRequest),
		Code:   http.StatusBadRequest,
		Error:  err.Error(),
	}
}

func FromForbiddenError(err error) *Response {
	return &Response{
		Status: http.StatusText(http.StatusForbidden),
		Code:   http.StatusForbidden,
		Error:  err.Error(),
	}
}

func FromNotFoundError(err error) *Response {
	if err != nil {
		return &Response{
			Status: http.StatusText(http.StatusNotFound),
			Code:   http.StatusNotFound,
			Error:  err.Error(),
		}
	}

	return &Response{
		Status: http.StatusText(http.StatusNotFound),
		Code:   http.StatusNotFound,
	}
}

func FromInternalServerError() *Response {
	return &Response{
		Status: http.StatusText(http.StatusInternalServerError),
		Code:   http.StatusNotFound,
	}
}

func (r *Response) WithCode(code int) *Response {
	return &Response{
		Status: http.StatusText(code),
		Code:   code,
	}
}

func (r *Response) WithEcho(c echo.Context) error {
	return c.JSON(r.Code, r)
}
