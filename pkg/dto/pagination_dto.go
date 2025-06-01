package dto

type PaginationRequestDto struct {
	Offset int `form:"offset" binding:"required,numeric"`
	Limit  int `form:"limit" binding:"required,numeric"`
}
