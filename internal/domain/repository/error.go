package repository

import "fmt"

var ErrorNotFound *ErrNotFound

type ErrNotFound struct {
	Resource string
	ID       string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("%s with id '%s' not found", e.Resource, e.ID)
}

func NewNotFound(resource, id string) error {
	return &ErrNotFound{Resource: resource, ID: id}
}
