package request

type GetById struct {
	Id uint `json:"id" form:"id" uri:"id" binding:"required"`
}

type PageInfo struct {
	Page     int `json:"page" form:"page" binding:"required,gt=0"`         // 页码
	PageSize int `json:"pageSize" form:"pageSize" binding:"required,gt=0"` // 每页大小
}
