package dto

type ListRequest struct {
	Offset int `form:"offset,default=0" binding:"numeric"`
	Limit  int `form:"limit,default=10" binding:"numeric"`
}

type SortOrder struct {
	Field     string `json:"field" binding:"required"`
	Dir       string `json:"dir" binding:"oneof=asc desc"`
	Collation string `json:"collation,omitempty"`
}
