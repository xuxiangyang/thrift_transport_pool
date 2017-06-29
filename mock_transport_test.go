package thrift_transport_pool

type MockTransport struct {
	Error    error
	IsOpened bool
}

func (this *MockTransport) Read(p []byte) (int, error) {
	return 0, this.Error
}

func (this *MockTransport) Write(p []byte) (int, error) {
	return 0, this.Error
}

func (this *MockTransport) Close() error {
	return nil
}

func (this *MockTransport) Flush() error {
	return this.Error
}

func (this *MockTransport) RemainingBytes() uint64 {
	return 0
}

func (this *MockTransport) Open() error {
	return this.Error
}

func (this *MockTransport) IsOpen() bool {
	return this.IsOpened
}
