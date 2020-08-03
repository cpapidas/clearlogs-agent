package mock

// MockTCPClient mock implementation
type MockTCPClient struct {
	ConnectFunc func(address string) error
	SendFunc    func(message string) error
}

// Connect mock implementation
func (m MockTCPClient) Connect(address string) error {
	if m.ConnectFunc(address) != nil {
		return m.ConnectFunc(address)
	}
	return nil
}

// Send mock implementation
func (m MockTCPClient) Send(message string) error {
	if m.SendFunc != nil {
		return m.SendFunc(message)
	}
	return nil
}
