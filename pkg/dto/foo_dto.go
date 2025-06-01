package dto

type FooReadDto struct {
	Id    int    `json:"id"`
	Label string `json:"label"`
	Bars  []int  `json:"bars"`
}

type FooCreateDto struct {
	Label  string `json:"label"`
	Secret string `json:"secret"`
}
