package entity

import (
	"fmt"
)

type FooKey struct {
	Id int
}

func (f FooKey) GetKey() string {
	return fmt.Sprintf("foo:%d", f.Id)
}
