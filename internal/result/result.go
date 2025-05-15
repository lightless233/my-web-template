package result

import (
	"runtime"

	"my-web-template/internal/constant"
)

type AppResult struct {
	Code       constant.ResultCode `json:"code"`
	Message    string              `json:"message"`
	Data       any                 `json:"data"`
	ErrorStack string              `json:"error_stack,omitempty"`
}

func NewAppResult(code constant.ResultCode, message string, data any, withStack ...bool) *AppResult {
	if len(withStack) > 0 && withStack[0] {
		stack := make([]byte, 1024*1024)
		n := runtime.Stack(stack, true)
		return &AppResult{
			Code:       code,
			Message:    message,
			Data:       data,
			ErrorStack: string(stack[:n]),
		}
	}
	return &AppResult{
		Code:       code,
		Message:    message,
		Data:       data,
		ErrorStack: "",
	}
}

func NewSuccessResult(data any) *AppResult {
	return NewAppResult(constant.CodeSuccess, "ok", data)
}

func NewErrorResult(code constant.ResultCode, message string, withStack ...bool) *AppResult {
	return NewAppResult(code, message, nil, withStack...)
}
