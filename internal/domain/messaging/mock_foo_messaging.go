package messaging

import "github.com/stretchr/testify/mock"

var (
	_ IFooMessaging = (*MockFooMessaging)(nil)
)

type MockFooMessaging struct {
	mock.Mock
}
