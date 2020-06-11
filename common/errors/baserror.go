package errors

import "fmt"

const (
	clientErr = 4
	serverErr = 5
	logicErr  = 6
)

// 这句话表示 BaseError 类型（非指针）实现 error 接口
var _ error = BaseError{}

// IBaseError ...
type IBaseError interface {
	Code() int
	Msg() string
	Err() error
}

// BaseError ...
type BaseError struct {
	Code int
	Msg  string
	// Err error // embedded error
}

// NewError ...
func NewError(code int, msg string) BaseError {
	return BaseError{Code: code, Msg: msg}
}

// Error interface for type error
func (e BaseError) Error() string {
	return fmt.Sprintf("(%d)%s", e.Code, e.Msg)
}

// Err get error
func (e BaseError) Err() error {
	if e.Code == 0 {
		return nil
	}
	return fmt.Errorf("%v", e.Msg)
}

// SetErr ...
func (e BaseError) SetErr(code int, err error) BaseError {
	e.Code = code
	e.Msg = func() string {
		if err == nil {
			return ""
		} else {
			return err.Error()
		}
	}()
	return e
}

// Equals ...
func (e BaseError) Equals(b BaseError) bool {
	return e.Code == b.Code
}
