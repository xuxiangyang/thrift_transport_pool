package thrift_transport_pool

type MockTransport struct {
	OpenError error
	Error     error
	IsOpened  bool
}

func (this *MockTransport) Read(p []byte) (int, error) {
	return 0, this.Error
}

func (this *MockTransport) Write(p []byte) (int, error) {
	return 0, nil
}

func (this *MockTransport) Close() error {
	return nil
}

func (this *MockTransport) Flush() error {
	return nil
}

func (this *MockTransport) RemainingBytes() uint64 {
	return 0
}

func (this *MockTransport) Open() error {
	this.IsOpened = true
	return this.OpenError
}

func (this *MockTransport) IsOpen() bool {
	return this.IsOpened
}
