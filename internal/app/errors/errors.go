package errors

import (
	"net/http"
)

type ErrorType uint8

const (
	InternalError ErrorType = iota
	ConflictError
	NotFoundError
)

type Error struct {
	ErrorCode ErrorType `json:"code"`
	HttpError int       `json:"-"`
	Message   string    `json:"message"`
}

type Message struct {
	Message string `json:"message"`
}

var CustomErrors = map[ErrorType]*Error{
	InternalError: {
		ErrorCode: InternalError,
		HttpError: http.StatusInternalServerError,
		Message:   "somthing wrong",
	},
	ConflictError: {
		ErrorCode: ConflictError,
		HttpError: http.StatusConflict,
		Message:   "user is registered",
	},
	NotFoundError: {
		ErrorCode: NotFoundError,
		HttpError: http.StatusNotFound,
		Message:   "Not Found",
	},
}


func UnexpectedInternal(err error) *Error {
	unexpErr := CustomErrors[InternalError]
	unexpErr.Message = err.Error()

	return unexpErr
}

func NotFoundBody(str string) *Error {
	nfErr := CustomErrors[NotFoundError]
	nfErr.Message = str

	return nfErr
}

func ConflictErrorBody(str string) *Error {
	nfErr := CustomErrors[ConflictError]
	nfErr.Message = str

	return nfErr
}
