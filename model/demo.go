package model

import "time"

type Role uint64

const (
	User Role = iota
	Runner
	Controler
)

type OrderStatus uint64

const (
	Published Role = iota
	Accepted
	Delivering
	Arrived
	Confirmed
)

// - `schools`：学校代码、名称、启用状态；所有业务数据通过 `school_id` 校验学校边界。
type Schools struct {
	SchoolID uint64
	Name     string
	Status   bool
}

// - `users`：用户名、密码哈希、角色、学校、手机号、账号状态、封禁原因、封禁时间、创建/更新时间。管理员可将 `school_id` 设为空以管理全部学校。
type Users struct {
	ID           uint64
	Username     string
	PasswordHash string
	Role         Role
	SchoolID     *uint64
	Phone        string
	Status       uint64
	BanReason    string
	BanTime      *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// - `runner_profiles`：与跑者用户一对一，保存完成订单数、好评、差评；订单和评分记录仍是统计事实来源。
type RunnerProfiles struct {
	UserID        uint64
	FinishOrders  uint64
	PositiveCount uint64
	NegativeCount uint64
}

// - `order_categories`：按学校配置订单类型、名称、排序和启用状态，避免把“代取快递”等类型写死为枚举。
type OrderCategories struct {
	ID          uint64
	SchoolID    uint64
	Code        string
	Name        string
	Description string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// - `orders`：业务订单号、学校、类型、标题、描述、发布者、接单跑者、固定报酬、取件地点、送达地点、联系人、联系电话、发布时间、截止时间、当前状态、接单/完成/取消时间、取消原因。`runner_id` 可为空。
type Orders struct {
	ID            uint64
	SchoolID      uint64
	CategoriesID  uint64
	Title         string
	Discribe      string
	UserID        uint64
	RunnerID      *uint64
	Remuneration  uint64
	StartPlace    string
	EndPlace      string
	UserPhone     string
	PublishedTime time.Time
	Deadline      time.Time
	Status        OrderStatus
	AcceptTime    time.Time
	FinishTime    time.Time
	CancelTime    time.Time
	CancelReason  string
}

// - `order_status_logs`：记录订单状态变更前后值、操作者类型、操作者 ID、原因和时间，供履约追踪和纠纷审计。
type OrderStatusLogs struct {
	ID           uint64
	OldOrder     Orders
	NewOrder     Orders
	UserRole     uint64
	UserID       uint64
	UpdateReason string
	UpdateTime   time.Time
}

// 交易与治理模型
// - `order_payments`：每个订单一条付款记录，保存金额、支付方式、付款状态、用户上报时间、跑者确认时间；支持待上报、待确认、已确认、争议、退款和取消。
type OrderPayments struct {
}

// - `order_payment_proofs`：保存付款截图等凭证，一笔付款允许多个附件。
// - `order_ratings`：发布者在订单完成后对跑者打 1–5 分；每个订单只能评分一次，不允许文字评价。
// - `disputes`：记录订单纠纷的举报人、被举报人、类型、说明、处理状态、管理员结论和处理时间。
// - `dispute_evidence`：保存纠纷双方提交的图片或其他证据地址。
// - `admin_audit_logs`：记录管理员封禁账号、处理纠纷、强制取消订单等操作及变更前后数据。
// - `notifications`：保存订单状态、付款确认、纠纷处理等站内通知和已读状态。
