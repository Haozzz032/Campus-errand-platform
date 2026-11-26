package Http

import "Campus-errand-platform/model"

// | 参数 | 类型 | 必填 | 说明 |
// |---|---|---:|---|
// | `school_id` | integer | 是 | 学校 ID |
// | `active_only` | boolean | 否 | 默认 `true`，仅返回启用分类 |
type GetOrdercategoriesRequest struct {
	SchoolID   uint64 `json:"scholl_id"`
	ActiveOnly bool   `json:"active_only"`
}

type GetOrdercategoriesResponse struct {
	Data []model.OrderCategories `json:"data"`
}
