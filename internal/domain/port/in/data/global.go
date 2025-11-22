package data

type PaginationOffset struct {
	Offset int
	Limit  int
}

type Optional[T any] struct {
	Value T
	Set   bool
}
