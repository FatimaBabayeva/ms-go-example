package ctmerror

import (
	"errors"
	"github.com/go-pg/pg"
	"net/http"
)

type MessageError struct {
	errorCode string
	err       error
	httpCode  int
}

// Error() func indicates that MessageError implements error interface
func (e *MessageError) Error() string {
	return e.errorCode
}

func (e *MessageError) HttpCode() int {
	return e.httpCode
}

func NewMessageError(repoError error) *MessageError {
	var msgError MessageError

	if errors.Is(repoError, pg.ErrNoRows) {
		msgError = MessageError{
			errorCode: "error.go-example.message-not-found",
			err:       repoError,
			httpCode:  http.StatusNotFound,
		}
	} else {
		msgError = MessageError{
			errorCode: "error.go-example.unexpected-error",
			err:       repoError,
			httpCode:  http.StatusInternalServerError,
		}
	}
	return &msgError
}
