package repository

import "fmt"

var (
	ErrorNotFound       *ErrNotFound
	ErrorNoAffectedData *ErrNoAffectedData
)

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

type ErrNoAffectedData struct {
	Resource string
	ID       string
}

func (e *ErrNoAffectedData) Error() string {
	return fmt.Sprintf("%s with id '%s' is not affected with request", e.Resource, e.ID)
}
