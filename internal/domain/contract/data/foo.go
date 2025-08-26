package data

import (
	"astigo/internal/domain/model"

	"github.com/google/uuid"
)

var (
	_ IFooUpdateMerger = (*FooUpdateInput)(nil)
	_ IFooUpdateMerger = (*FooPatchInput)(nil)
)

type IFooUpdateMerger interface {
	GetID() uuid.UUID
	Merge(foo *model.Foo) error
}

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

func (f *FooUpdateInput) GetID() uuid.UUID {
	return f.Id
}

func (f *FooUpdateInput) Merge(foo *model.Foo) error {
	foo.Label = f.Label
	foo.Secret = f.Secret
	foo.Value = f.Value
	foo.Weight = f.Weight

	return nil
}

type FooPatchInput struct {
	Id     uuid.UUID
	Label  Optional[string]
	Secret Optional[string]
	Value  Optional[int]
	Weight Optional[float32]
}

func (f *FooPatchInput) GetID() uuid.UUID {
	return f.Id
}

func (f *FooPatchInput) Merge(foo *model.Foo) error {
	if f.Label.Set {
		foo.Label = f.Label.Value
	}

	if f.Secret.Set {
		foo.Secret = f.Secret.Value
	}

	if f.Value.Set {
		foo.Value = f.Value.Value
	}

	if f.Weight.Set {
		foo.Weight = f.Weight.Value
	}

	return nil
}

type FooDeleteInput struct {
	Id uuid.UUID
}
