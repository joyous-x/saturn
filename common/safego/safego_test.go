package safego

import (
	"testing"
	"time"
)

func ExampleGo(t *testing.T) {
	Go(func() {
		a, b := 1, 0
		t.Logf("a / b = %v \n", a/b)
	})
}

func ExampleGo2(t *testing.T) {
	Go(func() {
		a, b := 1, 0
		t.Logf("a / b = %v \n", a/b)
	}, func(err interface{}) {
		t.Logf("=> my panic handler: %s \n", err)
	})
}

func ExampleGoWith(t *testing.T) {
	GoWith(func(args ...interface{}) {
		a := args[0].(int)
		b := args[1].(int)
		t.Logf("a / b = %v \n", a/b)
	}, func(err interface{}) {
		t.Logf("=> my panic handler: %s \n", err)
	})(1, 0)
}

func Test_Go_Panic(t *testing.T) {
	Go(func() {
		a, b := 1, 0
		t.Logf("a / b = %v \n", a/b)
	})

	time.Sleep(time.Second)
}

func Test_Go_Panic2(t *testing.T) {
	Go(func() {
		a, b := 1, 0
		t.Logf("a / b = %v \n", a/b)
	}, func(err interface{}) {
		t.Logf("=> my panic handler: %s \n", err)
	})

	time.Sleep(time.Second)
}

func Test_GoWith_Panic(t *testing.T) {
	GoWith(func(args ...interface{}) {
		a := args[0].(int)
		b := args[1].(int)
		t.Logf("a / b = %v \n", a/b)
	}, func(err interface{}) {
		t.Logf("=> my panic handler: %s \n", err)
	})(1, 0)

	time.Sleep(time.Second)
}
