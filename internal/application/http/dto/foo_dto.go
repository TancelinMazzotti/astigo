package dto

type FooReadRequest struct {
	Id int `uri:"id" binding:"required,numeric"`
}

type FooReadResponse struct {
	Id    int    `json:"id"`
	Label string `json:"label"`
}

type FooCreateRequest struct {
	Label  string `json:"label"`
	Secret string `json:"secret"`
}

type FooUpdateRequest struct {
	Id     int    `json:"id"`
	Label  string `json:"label"`
	Secret string `json:"secret"`
}

type FooDeleteRequest struct {
	Id int `uri:"id" binding:"required,numeric"`
}
