package dto

type BarReadDto struct {
	Id    int    `json:"id"`
	Label string `json:"label"`
}

type BarCreateDto struct {
	Label  string `json:"label"`
	Secret string `json:"secret"`
}
