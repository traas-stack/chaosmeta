package errors

import (
	"net/http"
)

func OK() Error {
	return NewError(http.StatusOK, http.StatusText(http.StatusOK), 0)
}

func ErrServer() Error {
	return NewError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), 2)
}

func ErrParam() Error {
	return NewError(http.StatusBadRequest, "the user is already disabled", 2)
}

func ErrSignParam() Error {
	return NewError(http.StatusForbidden, http.StatusText(http.StatusForbidden), 2)
}

func ErrUnauthorized() Error {
	return NewError(http.StatusUnauthorized, "Unauthorized: Invalid username or password.", 2)
}

func ErrNotFound() Error {
	return NewError(http.StatusNotFound, http.StatusText(http.StatusNotFound), 2)
}
