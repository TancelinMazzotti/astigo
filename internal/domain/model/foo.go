package model

import (
	"github.com/google/uuid"
	"sync"
	"time"
)

type Foo struct {
	Id     uuid.UUID `validate:"required"`
	Label  string    `validate:"required,min=3,max=100"`
	Secret string    `validate:"required,min=3,max=100"`
	Value  int       `validate:"required,gte=0,lte=1000"`
	Weight float32   `validate:"required,gte=0"`

	CreatedAt time.Time  `validate:"omitempty"`
	UpdatedAt *time.Time `validate:"omitempty"`

	Bars []*Bar `validate:"dive"`
}

var fooPool = &sync.Pool{
	New: func() interface{} {
		return &Foo{}
	},
}

func GetFoo() *Foo {
	return fooPool.Get().(*Foo)
}

func PutFoo(foo *Foo) {
	foo.Id = uuid.Nil
	foo.Label = ""
	foo.Secret = ""
	foo.Value = 0
	foo.Weight = 0

	foo.CreatedAt = time.Time{}
	foo.UpdatedAt = nil

	foo.Bars = nil

	fooPool.Put(foo)

}
