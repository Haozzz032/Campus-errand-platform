package Http

import "Campus-errand-platform/model"

// | 参数 | 类型 | 必填 | 说明 |
// |---|---|---:|---|
// | `school_id` | integer | 是 | 所属学校 |
// | `code` | string | 是 | 学校内唯一的机器编码 |
// | `name` | string | 是 | 显示名称 |
// | `description` | string | 否 | 分类说明 |
// | `icon_url` | string | 否 | 图标地址 |
// | `sort_order` | integer | 否 | 排序值，默认 0 |
// | `is_active` | boolean | 否 | 默认启用 |
type AddOrderCategoriesRequest struct {
	SchoolID    uint64 `json:"school_id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IconUrl     string `json:"icon_url"`
	SortOrder   int    `json:"sort_order"`
	IsActivite  bool   `json:"is_active"`
}

type AddOrderCategoriesResponse struct {
	Data model.OrderCategories `json:"data"`
}

// | 参数 | 类型 | 必填 | 说明 |
// |---|---|---:|---|
// | `name` | string | 否 | 显示名称 |
// | `description` | string | 否 | 分类说明 |
// | `icon_url` | string | 否 | 图标地址 |
// | `sort_order` | integer | 否 | 排序值 |
// | `is_active` | boolean | 否 | 是否启用 |
type ModOrderCategoriesRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IconUrl     string `json:"icon_url"`
	SortOrder   int    `json:"sort_order"`
	IsActivite  bool   `json:"is_active"`
}

type ModOrderCategoriesResponse struct {
	Data model.OrderCategories `json:"data"`
}
