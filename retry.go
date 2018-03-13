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
	buffer    [][]byte
	oldBuffer [][]byte
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
		buffer:    make([][]byte, 0),
		oldBuffer: make([][]byte, 0),
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
		for _, data := range this.oldBuffer {
			this.Transport.Write(data)
		}
		this.Transport.Flush()
		time.Sleep(time.Duration(i) * RETRY_DURATION)
	}
	return n, err
}

func (this *RetryedTransport) Write(p []byte) (n int, err error) {
	tmp := make([]byte, len(p))
	copy(tmp, p)
	this.buffer = append(this.buffer, tmp)
	for i := 0; i < MAX_TRY_TIMES; i++ {
		n, err = this.Transport.Write(p)
		if !IsNeedRetryError(err) {
			return n, err
		}
		this.Reconnect()
		for _, data := range this.buffer {
			this.Transport.Write(data)
		}
		time.Sleep(time.Duration(i) * RETRY_DURATION)
	}
	return n, err
}

func (this *RetryedTransport) Close() error {
	return this.Transport.Close()
}

func (this *RetryedTransport) Flush() (err error) {
	this.oldBuffer = this.buffer
	this.buffer = make([][]byte, 0)
	for i := 0; i < MAX_TRY_TIMES; i++ {
		err = this.Transport.Flush()
		if !IsNeedRetryError(err) {
			return err
		}
		this.Reconnect()
		for _, data := range this.oldBuffer {
			this.Transport.Write(data)
		}
		time.Sleep(time.Duration(i) * RETRY_DURATION)
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
		time.Sleep(time.Duration(i) * RETRY_DURATION)
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
	return this.Open()
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

	_, ok = err.(thrift.TProtocolException)
	if ok {
		return true
	}

	return false
}
