package promise

import (
	"fmt"
	"testing"
)

func TestNew_resolve(t *testing.T) {
	promise := New(func(resolve func(int), _ func(error)) {
		resolve(5)
		resolve(1)
	})
	data, err := promise.Await()
	if err != nil {
		t.Error(err)
	}
	if data != 5 {
		t.Error("wrong data")
	}
}

func TestNew_reject(t *testing.T) {
	happyError := fmt.Errorf("happy error")
	promise := New(func(_ func(int), reject func(error)) {
		reject(happyError)
		reject(fmt.Errorf("unhappy error"))
	})
	data, err := promise.Await()
	if err == nil {
		t.Error("error expected")
	}
	if data != 0 {
		t.Error("data should be 0")
	}
	if err.Error() != happyError.Error() {
		t.Error("wrong error")
	}
}

func TestNew_panic(t *testing.T) {
	promise := New(func(_ func(int), _ func(error)) {
		panic("test panic")
	})
	data, err := promise.Await()
	if err == nil {
		t.Error("error expected")
	}
	if data != 0 {
		t.Error("data should be 0")
	}
}

func TestNew_panicWithError(t *testing.T) {
	panicError := fmt.Errorf("test panic with error")
	promise := New(func(_ func(int), _ func(error)) {
		panic(panicError)
	})
	data, err := promise.Await()
	if err == nil {
		t.Error("error expected")
	}
	if data != 0 {
		t.Error("data should be 0")
	}
	if err != panicError {
		t.Error("wrong error")
	}
}

func TestNew_awaitTwice(t *testing.T) {
	promise := New(func(resolve func(int), _ func(error)) {
		resolve(5)
	})
	data, err := promise.Await()
	data2, err2 := promise.Await()
	if err != nil {
		t.Error(err)
	}
	if data != 5 {
		t.Error("wrong data")
	}
	if data2 != data {
		t.Error("second await should return same value as first")
	}
	if err2 != err {
		t.Error("second await should return same error (nil) as first")
	}
}
