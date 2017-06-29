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
	}

	return rt, rt.Open()
}

func (this *RetryedTransport) Read(p []byte) (n int, err error) {
	for i := 0; i < MAX_TRY_TIMES; i++ {
		n, err = this.Transport.Read(p)
		if !IsNeedRetryError(err) {
			return n, err
		}
		this.Reconnect()
		time.Sleep(RETRY_DURATION)
	}
	return n, err
}

func (this *RetryedTransport) Write(p []byte) (n int, err error) {
	for i := 0; i < MAX_TRY_TIMES; i++ {
		n, err = this.Transport.Write(p)
		if !IsNeedRetryError(err) {
			return n, err
		}
		this.Reconnect()
		time.Sleep(RETRY_DURATION)
	}
	return n, err
}

func (this *RetryedTransport) Close() error {
	return this.Transport.Close()
}

func (this *RetryedTransport) Flush() (err error) {
	for i := 0; i < MAX_TRY_TIMES; i++ {
		err = this.Transport.Flush()
		if !IsNeedRetryError(err) {
			return err
		}
		this.Reconnect()
		time.Sleep(RETRY_DURATION)
	}
	return err
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
	this.Transport, err = this.Block(this.HostPort)
	return err
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
