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
	size         int
	createdCount int
	pool         chan thrift.TTransport
	hostPort     string
	block        func(string) (thrift.TTransport, error)
	mu           sync.Mutex
}

func NewPool(size int, hostPort string, block func(hostPort string) (thrift.TTransport, error)) *Pool {
	return &Pool{
		TimeOut:      DEFAULT_TIME_OUT_TO_GET_CLIENT,
		size:         size,
		createdCount: 0,
		pool:         make(chan thrift.TTransport, size),
		hostPort:     hostPort,
		block:        block,
		mu:           sync.Mutex{},
	}
}

func (this *Pool) Pop() (thrift.TTransport, error) {
	this.mu.Lock()
	defer this.mu.Unlock()

	if this.createdCount < this.size {
		this.createdCount += 1
		return NewRetryedTransport(this.hostPort, this.block)
	}

	select {
	case t := <-this.pool:
		return t, nil
	case <-time.After(this.TimeOut):
		return nil, ErrPopTimeOut
	}
}

func (this *Pool) Push(t thrift.TTransport) {
	if len(this.pool) < this.size {
		this.pool <- t
	}
}
