package thrift_transport_pool

type MockTTransportException struct {
	timeout bool
	typeID  int
	err     error
}

func (this *MockTTransportException) Timeout() bool {
	return this.timeout
}

func (this *MockTTransportException) TypeId() int {
	return this.typeID
}

func (this *MockTTransportException) Err() error {
	return this.err
}

func (this *MockTTransportException) Error() string {
	return "MockTTransportException"
}
