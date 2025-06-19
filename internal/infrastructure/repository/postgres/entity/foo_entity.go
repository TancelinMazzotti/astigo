package entity

type Foo struct {
	FooId  int    `db:"foo_id"`
	Label  string `db:"label"`
	Secret string `db:"secret"`
}
