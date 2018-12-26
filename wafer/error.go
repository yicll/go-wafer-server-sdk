package wafer

import "fmt"

type WaferError struct {
	Code    int
	Message string
}

func (w WaferError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", w.Code, w.Message)
}

func (w WaferError) GetCode() int {
	return w.Code
}

func (w WaferError) GetMessage() string {
	return w.Message
}
