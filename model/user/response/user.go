package response

type GetUser struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`  // 用户名
	Phone string `json:"phone"` // 手机号
}

type GetUserList struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`  // 用户名
	Phone string `json:"phone"` // 手机号
}

type Login struct {
	User      GetUser `json:"userInfo"`
	Token     string  `json:"token"`
	ExpiresAt int64   `json:"expires_at"`
}
