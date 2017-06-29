package thrift_transport_pool

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"testing"
	"time"
)

func Test_Pop(t *testing.T) {
	pool := NewPool(1, "", func(hostPort string) (thrift.TTransport, error) {
		return &MockTransport{}, nil
	})

	_, err := pool.Pop()
	if err != nil {
		t.Fatal("should pop")
	}
	if pool.createdCount != 1 {
		t.Fatal("createdCount should eq 1")
	}
}

func Test_Pop_WithTimeout(t *testing.T) {
	pool := NewPool(1, "", func(hostPort string) (thrift.TTransport, error) {
		return &MockTransport{}, nil
	})
	pool.TimeOut = time.Millisecond

	pool.Pop()
	_, err := pool.Pop()
	if err != ErrPopTimeOut {
		t.Fatal("should return ErrTimeOut")
	}
}

func Test_Push(t *testing.T) {
	pool := NewPool(1, "", func(hostPort string) (thrift.TTransport, error) {
		return &MockTransport{}, nil
	})

	tr, err := pool.Pop()
	if err != nil {
		t.Fatal("should pop")
	}
	pool.Push(tr)
	if len(pool.pool) != 1 {
		t.Fatal("should push")
	}
}
