package entity

import (
	"fmt"
	"time"

	"github.com/TancelinMazzotti/astigo/internal/domain/model"

	"github.com/google/uuid"
)

// FooKey represents a unique identifier for a Foo entity using a UUID.
type FooKey struct {
	Id uuid.UUID
}

// GetKey generates a unique string key for a FooKey instance by combining the prefix "foo:" with the Id field.
func (f FooKey) GetKey() string {
	return fmt.Sprintf("foo:%s", f.Id)
}

// FooEntity represents an entity with unique ID, descriptive label, secret, numerical values, and timestamps.
type FooEntity struct {
	Id        uuid.UUID  `json:"id" redis:"id,omitempty"`
	Label     string     `json:"label" redis:"label,omitempty"`
	Secret    string     `json:"secret" redis:"secret,omitempty"`
	Value     int        `json:"value" redis:"value,omitempty"`
	Weight    float32    `json:"weight" redis:"weight,omitempty"`
	CreatedAt time.Time  `json:"createdAt" redis:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt" redis:"updated_at,omitempty"`
}

// ToModel converts the FooEntity instance into a model.Foo object.
func (f *FooEntity) ToModel() *model.Foo {
	return &model.Foo{
		Id:        f.Id,
		Label:     f.Label,
		Secret:    f.Secret,
		Value:     f.Value,
		Weight:    f.Weight,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}

// NewFooEntity creates a new instance of FooEntity from the provided model.Foo object.
func NewFooEntity(foo *model.Foo) *FooEntity {
	return &FooEntity{
		Id:        foo.Id,
		Label:     foo.Label,
		Secret:    foo.Secret,
		Value:     foo.Value,
		Weight:    foo.Weight,
		CreatedAt: foo.CreatedAt,
		UpdatedAt: foo.UpdatedAt,
	}
}
