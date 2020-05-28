package nils

import (
	"net"
	"reflect"
	"testing"
)

type TSlice []string

type IErrNil interface {
	IErrNil() IErrNil
	Error() string
}

type ErrNil struct {
	msg string
}

func (t *ErrNil) Error() string {
	return t.msg
}
func (t *ErrNil) ErrNil() *ErrNil {
	return nil
}
func (t *ErrNil) IErrNil() IErrNil {
	return nil
}

// 当我们用一个空指针类型的变量(如，var t *ErrNil)调用此方法时，该方法是会执行的，只有在执行该空指针变量的解指针操作(t.msg)时，才会 panic。
func (t *ErrNil) PrintMsg() string {
	if t == nil {
		return "<nil>"
	}
	return t.msg
}

// 由于接受者是 ErrNil 而不是 *ErrNil，使用指针访问该函数时，Golang 内部会在调用时自动解指针，故使用空指针类型的变量(如，var t *ErrNil)调用此方法时会 panic。
func (t ErrNil) PrintMsgV2() string {
	return t.msg
}

// FAQ: https://golang.org/doc/faq#nil_error
//     简单说，interface 被两个元素 value 和 type 所表示。只有在 value 和 type 同时为 nil 的时候，判断 interface == nil 才会为 true。
func Test_Nil(t *testing.T) {
	var err error

	tmp := ErrNil{}
	err = tmp.ErrNil()
	if err == nil {
		t.Logf("(ErrNil == nil) ok: %v, err: %v", err == nil, err)
	} else {
		t.Errorf("(ErrNil == nil) err: %v, err: %v, type(err): %v", err == nil, err, reflect.TypeOf(err).Kind())
	}

	err = tmp.IErrNil()
	if err == nil {
		t.Logf("(IErrNil == nil) ok: %v, err: %v", err == nil, err)
	} else {
		t.Errorf("(IErrNil == nil) err: %v, err: %v, type(err): %v", err == nil, err, reflect.TypeOf(err).Kind())
	}

	ip := net.ParseIP("111.1.111")
	if ip == nil {
		t.Logf("(ip == nil) ok: %v, err: %v, type(ip): %v", ip == nil, ip, reflect.TypeOf(ip).Kind())
	} else {
		t.Errorf("(ip == nil) err: %v, err: %v, type(ip): %v", ip == nil, ip, reflect.TypeOf(ip).Kind())
	}
}

// slice 的 nil指针 问题
func Test_SliceNil(t *testing.T) {
	var a TSlice = nil
	var b []string = nil
	var c ErrNil
	t.Errorf("(var a TSlice = nil) isNil:%v type: %v", a == nil, reflect.TypeOf(a).Kind())
	t.Errorf("(var b []string = nil) isNil:%v type: %v", b == nil, reflect.TypeOf(b).Kind())
	// mismatched types ErrNil and nil
	// t.Errorf("(var c ErrNil) isNil:%v type: %v", c==nil, reflect.TypeOf(c).Kind())
	t.Errorf("(var c ErrNil) type: %v", reflect.TypeOf(c).Kind())
}
