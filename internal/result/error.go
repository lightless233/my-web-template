package result

import (
	"fmt"
	"runtime"

	"my-web-template/internal/constant"
)

type AppError interface {
	error
	ToAppResult() *AppResult
}

type appError struct {
	Code       constant.ResultCode
	Message    string
	ErrorStack string
}

func (e *appError) ToAppResult() *AppResult {
	return &AppResult{
		Code:       e.Code,
		Message:    e.Message,
		Data:       nil,
		ErrorStack: e.ErrorStack,
	}
}

func (e *appError) Error() string {
	return fmt.Sprintf("[APP_ERROR|%s|%d] Message=%s\nStack=%s", constant.GetResultCodeName(e.Code), e.Code, e.Message, e.ErrorStack)
}

func NewAppError(code constant.ResultCode, message string, withStack ...bool) AppError {
	stackString := ""
	if len(withStack) > 0 && withStack[0] {
		stack := make([]byte, 1024*1024)
		n := runtime.Stack(stack, true)
		stackString = string(stack[:n])
	}
	return &appError{
		Code:       code,
		Message:    message,
		ErrorStack: stackString,
	}
}

func NewAppErrorFromError(code constant.ResultCode, err error, withStack ...bool) AppError {
	stackString := ""
	if len(withStack) > 0 && withStack[0] {
		stack := make([]byte, 1024*1024)
		n := runtime.Stack(stack, true)
		stackString = string(stack[:n])
	}
	return &appError{
		Code:       code,
		Message:    err.Error(),
		ErrorStack: stackString,
	}
}
