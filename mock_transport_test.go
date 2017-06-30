package thrift_transport_pool

type MockTransport struct {
	Error    error
	IsOpened bool
	MockOpen func() error
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
	if this.MockOpen != nil {
		err := this.MockOpen()
		if err == nil {
			this.IsOpened = true
			return nil
		} else {
			return err
		}
	} else {
		this.IsOpened = true
		return nil
	}

}

func (this *MockTransport) IsOpen() bool {
	return this.IsOpened
}
