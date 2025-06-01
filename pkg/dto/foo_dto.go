package dto

type FooRequestReadDto struct {
	Id int `uri:"id" binding:"required,numeric"`
}

type FooResponseReadDto struct {
	Id    int    `json:"id"`
	Label string `json:"label"`
	Bars  []int  `json:"bars"`
}

type FooRequestCreateDto struct {
	Label  string `json:"label"`
	Secret string `json:"secret"`
}

type FooRequestUpdateDto struct {
	Id     int    `json:"id"`
	Label  string `json:"label"`
	Bars   []int  `json:"bars"`
	Secret string `json:"secret"`
}

type FooRequestDeleteDto struct {
	Id int `uri:"id" binding:"required,numeric"`
}
