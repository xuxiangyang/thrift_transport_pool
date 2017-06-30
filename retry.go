package thrift_transport_pool

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"time"
)

const (
	MAX_TRY_TIMES  = 5
	RETRY_DURATION = time.Duration(10) * time.Millisecond
)

type RetryedTransport struct {
	Transport thrift.TTransport
	Block     func(hostPort string) (thrift.TTransport, error)
	HostPort  string
	buffer    []byte
}

func NewRetryedTransport(hostPort string, block func(hostPort string) (thrift.TTransport, error)) (*RetryedTransport, error) {
	transport, err := block(hostPort)
	if err != nil {
		return nil, err
	}

	rt := &RetryedTransport{
		Transport: transport,
		Block:     block,
		HostPort:  hostPort,
		buffer:    make([]byte, 0),
	}

	return rt, rt.Open()
}

func (this *RetryedTransport) Read(p []byte) (n int, err error) {
	for i := 0; i < MAX_TRY_TIMES; i++ {
		n, err = this.Transport.Read(p)
		if !IsNeedRetryError(err) {
			return n, err
		}
		time.Sleep(RETRY_DURATION)
		this.Reconnect()
		this.Transport.Write(this.buffer)
		this.Transport.Flush()
	}
	return n, err
}

func (this *RetryedTransport) Write(p []byte) (n int, err error) {
	this.buffer = p
	return this.Transport.Write(p)
}

func (this *RetryedTransport) Close() error {
	return this.Transport.Close()
}

func (this *RetryedTransport) Flush() error {
	return this.Transport.Flush()
}

func (this *RetryedTransport) RemainingBytes() uint64 {
	return this.Transport.RemainingBytes()
}

func (this *RetryedTransport) Open() (err error) {
	if this.IsOpen() {
		return nil
	}

	for i := 0; i < MAX_TRY_TIMES; i++ {
		err = this.Transport.Open()
		if !IsNeedRetryError(err) {
			return err
		}
		this.Reconnect()
		time.Sleep(RETRY_DURATION)
	}
	return err
}

func (this *RetryedTransport) IsOpen() bool {
	return this.Transport.IsOpen()
}

func (this *RetryedTransport) Reconnect() (err error) {
	this.Transport.Close()
	this.Transport, err = this.Block(this.HostPort)
	if err != nil {
		return err
	}
	return this.Transport.Open()
}

func IsNeedRetryError(err error) bool {
	if err == nil {
		return false
	}

	var ok bool
	_, ok = err.(thrift.TTransportException)
	if ok {
		return true
	}

	return false
}
