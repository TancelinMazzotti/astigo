package data

import (
	"astigo/internal/domain/model"
	"github.com/google/uuid"
)

type FooReadListInput struct {
	Offset int
	Limit  int
}

type FooReadInput struct {
	Id uuid.UUID
}

type FooCreateInput struct {
	Label  string
	Secret string
	Value  int
	Weight float32
}

type FooUpdateInput struct {
	Id     uuid.UUID
	Label  string
	Secret string
	Value  int
	Weight float32
}

func (f FooUpdateInput) Merge(foo *model.Foo) error {
	foo.Label = f.Label
	foo.Secret = f.Secret
	foo.Value = f.Value
	foo.Weight = f.Weight

	return nil
}

type FooDeleteInput struct {
	Id uuid.UUID
}
