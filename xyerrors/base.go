package xyerrors

import "fmt"

type XyError struct {
	msg   string
	code  int
	errno int
}

func (se XyError) StatusCode() int {
	return se.code
}

func (se XyError) Errno() int {
	return se.errno
}

func (se XyError) Error() string {
	return se.msg
}

func (se XyError) New(msg string, args ...interface{}) XyError {
	return XyError{
		msg:   fmt.Sprintf(msg, args...),
		errno: se.errno,
		code:  se.code,
	}
}

func NewError(code int, errno int) XyError {
	return XyError{
		errno: errno,
		code:  code,
	}
}
