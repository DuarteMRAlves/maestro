package events

import "github.com/DuarteMRAlves/maestro/internal/api"

// MockPubSub implements the PubSub interface but does no underlying processing.
type MockPubSub struct {
}

func (m *MockPubSub) Subscribe() *api.Subscription {
	return nil
}

func (m *MockPubSub) Unsubscribe(token api.SubscriptionToken) error {
	return nil
}

func (m *MockPubSub) Publish(s string) {
}

func (m *MockPubSub) Close() {
}
