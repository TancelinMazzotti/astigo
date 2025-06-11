package handler

type PaginationInput struct {
	Offset int `form:"offset" binding:"numeric"`
	Limit  int `form:"limit" binding:"numeric"`
}
