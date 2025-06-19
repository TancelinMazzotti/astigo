package entity

import (
	"fmt"
	"github.com/google/uuid"
)

type FooKey struct {
	Id uuid.UUID
}

func (f FooKey) GetKey() string {
	return fmt.Sprintf("foo:%s", f.Id)
}
