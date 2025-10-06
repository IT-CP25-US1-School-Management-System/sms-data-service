package errs

import (
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AppError struct {
	Message    string
	HTTPStatus int
	GRPCCode   int
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%v", e.Err)
	}
	return e.Message
}

func NewAppError(message string, httpStatus int, grpcCode int, err error) *AppError {
	return &AppError{
		Message:    message,
		HTTPStatus: httpStatus,
		GRPCCode:   grpcCode,
		Err:        err,
	}
}

func ToGRPCError(err error) error {
	if appErr, ok := err.(*AppError); ok {
		return status.Error(codes.Code(appErr.GRPCCode), appErr.Message)
	}
	return status.Error(codes.Unknown, err.Error())
}

var (
	ErrNotFound = func(err error) *AppError {
		return NewAppError(err.Error(), http.StatusNotFound, int(codes.NotFound), err)
	}
	ErrUnauthorized = func(err error) *AppError {
		return NewAppError(err.Error(), http.StatusUnauthorized, int(codes.Unauthenticated), err)
	}
	ErrBadRequest = func(err error) *AppError {
		return NewAppError(err.Error(), http.StatusBadRequest, int(codes.InvalidArgument), err)
	}
	ErrForbidden = func(err error) *AppError {
		return NewAppError(err.Error(), http.StatusForbidden, int(codes.PermissionDenied), err)
	}
	ErrConflict = func(err error) *AppError {
		return NewAppError(err.Error(), http.StatusConflict, int(codes.AlreadyExists), err)
	}
	ErrNoContent = func() *AppError {
		return NewAppError("", http.StatusNoContent, int(codes.NotFound), nil)
	}
	ErrInternalServer = func(err error) *AppError {
		return NewAppError(err.Error(), http.StatusInternalServerError, int(codes.Internal), err)
	}
)

// Helper functions to create errors with specific messages
func NewNotFoundError(message string) *AppError {
	return NewAppError(message, http.StatusNotFound, int(codes.NotFound), nil)
}

func NewUnauthorizedError(message string) *AppError {
	return NewAppError(message, http.StatusUnauthorized, int(codes.Unauthenticated), nil)
}

func NewBadRequestError(message string) *AppError {
	return NewAppError(message, http.StatusBadRequest, int(codes.InvalidArgument), nil)
}

func NewInternalServerError(message string) *AppError {
	return NewAppError(message, http.StatusInternalServerError, int(codes.Internal), nil)
}
