package thrift_transport_pool

import (
	"errors"
	"git.apache.org/thrift.git/lib/go/thrift"
	"testing"
)

func Test_IsNeedRetryError_with_nil(t *testing.T) {
	if IsNeedRetryError(nil) {
		t.Fatal("Nil should not retry")
	}
}

func Test_IsNeedRetryError_with_TTransportException(t *testing.T) {
	if !IsNeedRetryError(&MockTTransportException{}) {
		t.Fatal("TTransportException should not retry")
	}
}

func Test_IsNeedRetryError_with_other_erros(t *testing.T) {
	if IsNeedRetryError(errors.New("sjk")) {
		t.Fatal("Nil should not retry")
	}
}

func Test_Read_with_retry(t *testing.T) {
	rt := newRetryedTransport(t)
	firstTransport := rt.Transport
	_, err := rt.Read([]byte{})
	if err != nil {
		t.Fatal("should not fail")
	}
	if rt.Transport == firstTransport {
		t.Fatal("Transport should change")
	}
}

func Test_Wirte_with_retry(t *testing.T) {
	rt := newRetryedTransport(t)
	firstTransport := rt.Transport
	_, err := rt.Write([]byte{})
	if err != nil {
		t.Fatal("should not fail")
	}
	if rt.Transport == firstTransport {
		t.Fatal("Transport should change")
	}
}

func Test_Open_with_retry(t *testing.T) {
	rt := newRetryedTransport(t)
	firstTransport := rt.Transport
	err := rt.Open()
	if err != nil {
		t.Fatal("should not fail")
	}

	if rt.Transport == firstTransport {
		t.Fatal("Transport should change")
	}
}

func Test_Flush_with_retry(t *testing.T) {
	rt := newRetryedTransport(t)
	firstTransport := rt.Transport
	err := rt.Flush()
	if err != nil {
		t.Fatal("should not fail")
	}

	if rt.Transport == firstTransport {
		t.Fatal("Transport should change")
	}
}

func newRetryedTransport(t *testing.T) *RetryedTransport {
	try_times := 0
	f := func(hostPort string) (thrift.TTransport, error) {
		t.Log("sjk")
		defer func() { try_times += 1 }()
		if try_times == 0 {
			return &MockTransport{Error: &MockTTransportException{}}, nil
		} else {
			return &MockTransport{}, nil
		}

	}
	rt, err := NewRetryedTransport("hostport", f)
	if err != nil {
		t.Fatal(err)
		return nil
	}
	if rt.Transport == nil {
		t.Fatal("Transport should not be nil after New")
		return nil
	}
	return rt
}
