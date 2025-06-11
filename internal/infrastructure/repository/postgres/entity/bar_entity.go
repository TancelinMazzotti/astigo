package entity

type Bar struct {
	BarId int    `db:"bar_id"`
	Label string `db:"label"`
	FooId int    `db:"foo_id"`
}
