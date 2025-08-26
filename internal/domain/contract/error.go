package contract

import "fmt"

var (
	ErrorNotFound         *ErrNotFound
	ErrorAlreadyExists    *ErrAlreadyExists
	ErrorNoAffectedData   *ErrNoAffectedData
	ErrorInvalidReference *ErrInvalidReference
)

type ErrNotFound struct {
	Resource string
	Field    string
	Value    string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("%s with %s '%s' not found", e.Resource, e.Field, e.Value)
}

func NewErrNotFound(resource, field, value string) error {
	return &ErrNotFound{Resource: resource, Field: field, Value: value}
}

type ErrAlreadyExists struct {
	Resource string
	Field    string
	Value    string
}

func (e *ErrAlreadyExists) Error() string {
	return fmt.Sprintf("%s with %s '%s' already exists", e.Resource, e.Field, e.Value)
}

func NewErrAlreadyExists(resource, field, value string) error {
	return &ErrAlreadyExists{Resource: resource, Field: field, Value: value}
}

type ErrNoAffectedData struct {
	Resource string
	ID       string
}

func (e *ErrNoAffectedData) Error() string {
	return fmt.Sprintf("%s with id '%s' is not affected with request", e.Resource, e.ID)
}

func NewErrNoAffectedData(resource, id string) error {
	return &ErrNoAffectedData{Resource: resource, ID: id}
}

type ErrInvalidReference struct {
	Resource string
	Field    string
	Value    string
}

func (e *ErrInvalidReference) Error() string {
	return fmt.Sprintf("invalid reference for %s with %s '%s'", e.Resource, e.Field, e.Value)
}

func NewErrInvalidReference(resource, field, value string) error {
	return &ErrInvalidReference{Resource: resource, Field: field, Value: value}
}
