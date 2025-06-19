package tool

func NewPointer[T any](t T) *T {
	return &t
}
func NewPointerNil[T any]() *T {
	return nil
}
