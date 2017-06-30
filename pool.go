package thrift_transport_pool

import (
	"errors"
	"git.apache.org/thrift.git/lib/go/thrift"
	"sync"
	"time"
)

const DEFAULT_TIME_OUT_TO_GET_CLIENT = 5 * time.Second

var ErrPopTimeOut = errors.New("Timout to pop a client")

type Pool struct {
	TimeOut      time.Duration
	Size         int
	CreatedCount int
	Pool         chan thrift.TTransport
	HostPort     string
	Block        func(string) (thrift.TTransport, error)
	mu           sync.Mutex
}

func NewPool(size int, hostPort string, block func(hostPort string) (thrift.TTransport, error)) *Pool {
	return &Pool{
		TimeOut:      DEFAULT_TIME_OUT_TO_GET_CLIENT,
		Size:         size,
		CreatedCount: 0,
		Pool:         make(chan thrift.TTransport, size),
		HostPort:     hostPort,
		Block:        block,
		mu:           sync.Mutex{},
	}
}

func (this *Pool) Pop() (thrift.TTransport, error) {
	this.mu.Lock()
	defer this.mu.Unlock()

	if this.CreatedCount < this.Size {
		this.CreatedCount += 1
		return NewRetryedTransport(this.HostPort, this.Block)
	}

	select {
	case t := <-this.Pool:
		return t, nil
	case <-time.After(this.TimeOut):
		return nil, ErrPopTimeOut
	}
}

func (this *Pool) Push(t thrift.TTransport) {
	if len(this.Pool) < this.Size {
		this.Pool <- t
	}
}
