package Http

import (
	"Campus-errand-platform/model"
)

// | 参数 | 类型 | 必填 | 说明 |
// |---|---|---:|---|
// | `username` | string | 是 | 3-50 个字符，全平台唯一 |
// | `password` | string | 是 | 8-72 个字符，服务端只保存哈希 |
// | `role` | string | 是 | `publisher` 或 `runner` |
// | `school_id` | integer | 是 | 已启用的学校 ID |
// | `phone` | string | 是 | 手机号，用于订单联系 |
type RegistRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	SchoolID uint64 `json:"school_id"`
	Phone    string `json:"phone"`
}

// | 参数 | 类型 | 必填 | 说明 |
// |---|---|---:|---|
// | `username` | string | 是 | 用户名 |
// | `password` | string | 是 | 明文密码，仅用于校验 |
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ```json
//
//	{
//	  "access_token": "<jwt>",
//	  "token_type": "Bearer",
//	  "expires_in": 604800,
//	  "user": {}
//	}
//
// ```
type RegistResponse struct {
	AccessToken string      `json:"access_token"`
	TokenType   string      `json:"token_type"`
	ExpiresIn   uint64      `json:"expires_in"`
	User        model.Users `json:"user"`
}

type LoginResponse struct {
	AccessToken string      `json:"access_token"`
	TokenType   string      `json:"token_type"`
	ExpiresIn   uint64      `json:"expires_in"`
	User        model.Users `json:"user"`
}
