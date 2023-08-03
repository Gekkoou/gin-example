package request

type Login struct {
	Name     string `json:"name" form:"name" binding:"required,alphanum"`         // 用户名
	Password string `json:"password" form:"password" binding:"required,alphanum"` // 密码
}

type Register struct {
	Name     string `json:"name" form:"name" binding:"required,alphanum"`         // 用户名
	Password string `json:"password" form:"password" binding:"required,alphanum"` // 密码
	Phone    string `json:"phone" form:"phone" binding:"omitempty,number"`        // 手机号
}

type ChangePassword struct {
	Id          uint64 `json:"-"`                                                          // 从 JWT 中提取 user id，避免越权
	Password    string `json:"password" form:"password" binding:"required,alphanum"`       // 密码
	NewPassword string `json:"newPassword" form:"newPassword" binding:"required,alphanum"` // 新密码
}

type UpdateUser struct {
	Id    uint64 `json:"-"`                                             // 从 JWT 中提取 user id，避免越权
	Name  string `json:"name" form:"name" binding:"required,alphanum"`  // 用户名
	Phone string `json:"phone" form:"phone" binding:"omitempty,number"` // 手机号
}

type DeleteUser struct {
	Id uint64 `json:"id" form:"id" binding:"required"`
}
