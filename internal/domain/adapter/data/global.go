package data

type Optional[T any] struct {
	Value T
	Set   bool
}
