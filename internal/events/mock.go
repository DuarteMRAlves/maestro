package events

// MockPubSub implements the PubSub interface but does no underlying processing.
type MockPubSub struct {
}

func (m *MockPubSub) Subscribe() *Subscription {
	return nil
}

func (m *MockPubSub) Unsubscribe(token SubscriptionToken) error {
	return nil
}

func (m *MockPubSub) Publish(s string) {
}

func (m *MockPubSub) Close() {
}
