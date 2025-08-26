package message

import (
	"time"

	"github.com/TancelinMazzotti/astigo/internal/domain/model"

	"github.com/google/uuid"
)

// FooMessage represents a data transfer object for Foo, used for messaging or serialization purposes.
type FooMessage struct {
	Id        uuid.UUID  `json:"id"`
	Label     string     `json:"label"`
	Value     int        `json:"value"`
	Weight    float32    `json:"weight"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// NewFooMessage transforms a model.Foo instance into a corresponding FooMessage instance for external usage or serialization.
func NewFooMessage(foo *model.Foo) *FooMessage {
	return &FooMessage{
		Id:        foo.Id,
		Label:     foo.Label,
		Value:     foo.Value,
		Weight:    foo.Weight,
		CreatedAt: foo.CreatedAt,
		UpdatedAt: foo.UpdatedAt,
	}
}
