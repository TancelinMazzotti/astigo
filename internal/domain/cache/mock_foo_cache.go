package cache

import "github.com/stretchr/testify/mock"

var (
	_ IFooCahe = (*MockFooCache)(nil)
)

type MockFooCache struct {
	mock.Mock
}
