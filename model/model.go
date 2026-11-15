package model

import "time"

// Role 表示账号类型。
type Role uint8

const (
	RoleUser   Role = iota // 发布者
	RoleRunner             // 跑者
	RoleAdmin              // 管理员
)

// UserStatus 表示账号当前状态。
type UserStatus uint8

const (
	UserStatusActive   UserStatus = iota // 正常
	UserStatusBanned                     // 已封禁
	UserStatusDisabled                   // 已停用
)

// OrderStatus 表示订单履约状态。
type OrderStatus uint8

const (
	OrderStatusPublished  OrderStatus = iota // 已发布，等待接单
	OrderStatusAccepted                      // 已接单
	OrderStatusDelivering                    // 配送中
	OrderStatusArrived                       // 已送达，等待确认
	OrderStatusConfirmed                     // 已确认完成
	OrderStatusCanceled                      // 已取消
)

// ActorType 表示订单操作者类型，system 用于超时自动确认等系统操作。
type ActorType string

const (
	ActorTypePublisher ActorType = "publisher"
	ActorTypeRunner    ActorType = "runner"
	ActorTypeAdmin     ActorType = "admin"
	ActorTypeSystem    ActorType = "system"
)

// ViolationStatus 表示举报与违规案件的处理状态。
type ViolationStatus uint8

const (
	ViolationStatusReviewing ViolationStatus = iota // 审核中
	ViolationStatusConfirmed                        // 已确认违规
	ViolationStatusRejected                         // 已驳回
	ViolationStatusClosed                           // 已关闭
)

// PenaltyType 表示确认违规后的处罚类型。
type PenaltyType uint8

const (
	PenaltyTypeWarning PenaltyType = iota // 警告
	PenaltyTypeSuspend                    // 暂停账号
	PenaltyTypeBan                        // 封禁账号
)

// PaymentMethod 表示线下支付方式。
type PaymentMethod uint8

const (
	PaymentMethodWeChat PaymentMethod = iota // 微信
	PaymentMethodAlipay                      // 支付宝
	PaymentMethodCash                        // 现金
	PaymentMethodOther                       // 其他
)

// PaymentStatus 表示订单付款确认状态。
type PaymentStatus uint8

const (
	PaymentStatusPending   PaymentStatus = iota // 待上报
	PaymentStatusReported                       // 用户已上报，待跑者确认
	PaymentStatusConfirmed                      // 跑者已确认收款
	PaymentStatusDisputed                       // 付款存在争议
	PaymentStatusRefunded                       // 已退款
	PaymentStatusCanceled                       // 已取消
)

// Schools 保存学校基础信息。
type Schools struct {
	SchoolID  uint64    `gorm:"column:school_id;primaryKey;autoIncrement"`
	Code      string    `gorm:"column:code;type:varchar(32);not null;uniqueIndex:uk_schools_code"`
	Name      string    `gorm:"column:name;type:varchar(100);not null"`
	Status    bool      `gorm:"column:status;not null;default:true;index:idx_schools_status"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime(3);not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime(3);not null"`
}

// Users 保存统一账号、角色和当前封禁状态。
type Users struct {
	ID             uint64     `gorm:"column:id;primaryKey;autoIncrement"`
	Username       string     `gorm:"column:username;type:varchar(50);not null;uniqueIndex:uk_users_username"`
	PasswordHash   string     `gorm:"column:password_hash;type:varchar(255);not null" json:"-"`
	Role           Role       `gorm:"column:role;type:tinyint unsigned;not null;default:0;check:chk_users_role,role BETWEEN 0 AND 2;index:idx_users_school_role,priority:2"`
	SchoolID       *uint64    `gorm:"column:school_id;index:idx_users_school_role,priority:1"`
	Phone          string     `gorm:"column:phone;type:varchar(20);not null;default:'';index:idx_users_phone"`
	Status         UserStatus `gorm:"column:status;type:tinyint unsigned;not null;default:0;check:chk_users_status,status BETWEEN 0 AND 2;index:idx_users_status"`
	BanReason      string     `gorm:"column:ban_reason;type:varchar(500);not null;default:''"`
	ViolationTimes uint64     `gorm:"column:violation_times;not null;default:0"`
	BannedAt       *time.Time `gorm:"column:banned_at;type:datetime(3)"`
	BannedUntil    *time.Time `gorm:"column:banned_until;type:datetime(3);index:idx_users_banned_until"`
	CreatedAt      time.Time  `gorm:"column:created_at;type:datetime(3);not null"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;type:datetime(3);not null"`
}

// RunnerProfiles 保存跑者统计缓存，事实来源为订单和评价表。
type RunnerProfiles struct {
	UserID              uint64    `gorm:"column:user_id;primaryKey"`
	CompletedOrderCount uint64    `gorm:"column:completed_order_count;not null;default:0"`
	PositiveCount       uint64    `gorm:"column:positive_count;not null;default:0"`
	NegativeCount       uint64    `gorm:"column:negative_count;not null;default:0"`
	CreatedAt           time.Time `gorm:"column:created_at;type:datetime(3);not null"`
	UpdatedAt           time.Time `gorm:"column:updated_at;type:datetime(3);not null"`
}

// UserProfiles 保存发布者统计缓存。
type UserProfiles struct {
	UserID              uint64    `gorm:"column:user_id;primaryKey"`
	PublishedOrderCount uint64    `gorm:"column:published_order_count;not null;default:0"`
	PositiveCount       uint64    `gorm:"column:positive_count;not null;default:0"`
	NegativeCount       uint64    `gorm:"column:negative_count;not null;default:0"`
	CreatedAt           time.Time `gorm:"column:created_at;type:datetime(3);not null"`
	UpdatedAt           time.Time `gorm:"column:updated_at;type:datetime(3);not null"`
}

// RunnerRatings 保存发布者对跑者的评价。
type RunnerRatings struct {
	ID         uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	OrderID    uint64    `gorm:"column:order_id;not null;uniqueIndex:uk_runner_ratings_order"`
	Content    string    `gorm:"column:content;type:varchar(500);not null;default:''"`
	UserID     uint64    `gorm:"column:user_id;not null;index:idx_runner_ratings_user"`
	RunnerID   uint64    `gorm:"column:runner_id;not null;index:idx_runner_ratings_runner"`
	IsPositive bool      `gorm:"column:is_positive;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime(3);not null"`
}

// UserRatings 保存跑者对发布者的评价。
type UserRatings struct {
	ID         uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	OrderID    uint64    `gorm:"column:order_id;not null;uniqueIndex:uk_user_ratings_order"`
	Content    string    `gorm:"column:content;type:varchar(500);not null;default:''"`
	UserID     uint64    `gorm:"column:user_id;not null;index:idx_user_ratings_user"`
	RunnerID   uint64    `gorm:"column:runner_id;not null;index:idx_user_ratings_runner"`
	IsPositive bool      `gorm:"column:is_positive;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime(3);not null"`
}

// OrderCategories 是按学校维护的订单类型字典。
type OrderCategories struct {
	ID          uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	SchoolID    uint64    `gorm:"column:school_id;not null;uniqueIndex:uk_order_categories_school_code,priority:1;index:idx_order_categories_list,priority:1"`
	Code        string    `gorm:"column:code;type:varchar(32);not null;uniqueIndex:uk_order_categories_school_code,priority:2"`
	Name        string    `gorm:"column:name;type:varchar(50);not null"`
	Description string    `gorm:"column:description;type:varchar(255);not null;default:''"`
	IconURL     string    `gorm:"column:icon_url;type:varchar(512);not null;default:''"`
	SortOrder   int       `gorm:"column:sort_order;not null;default:0;index:idx_order_categories_list,priority:3"`
	IsActive    bool      `gorm:"column:is_active;not null;default:true;index:idx_order_categories_list,priority:2"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime(3);not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime(3);not null"`
}

// Orders 保存订单发布、抢单、履约、确认和取消信息。
type Orders struct {
	ID                uint64      `gorm:"column:id;primaryKey;autoIncrement"`
	OrderNo           string      `gorm:"column:order_no;type:varchar(32);not null;uniqueIndex:uk_orders_order_no"`
	SchoolID          uint64      `gorm:"column:school_id;not null;index:idx_orders_school_status_created,priority:1"`
	CategoryID        uint64      `gorm:"column:category_id;not null;index:idx_orders_category"`
	Title             string      `gorm:"column:title;type:varchar(100);not null"`
	Description       string      `gorm:"column:description;type:text;not null"`
	PublisherID       uint64      `gorm:"column:publisher_id;not null;index:idx_orders_publisher_created,priority:1"`
	RunnerID          *uint64     `gorm:"column:runner_id;index:idx_orders_runner_status,priority:1"`
	RemunerationCents int64       `gorm:"column:remuneration_cents;not null;check:chk_orders_remuneration,remuneration_cents >= 0"`
	StartPlace        string      `gorm:"column:start_place;type:varchar(255);not null"`
	EndPlace          string      `gorm:"column:end_place;type:varchar(255);not null"`
	ContactName       string      `gorm:"column:contact_name;type:varchar(50);not null"`
	UserPhone         string      `gorm:"column:user_phone;type:varchar(20);not null"`
	Status            OrderStatus `gorm:"column:status;type:tinyint unsigned;not null;default:0;check:chk_orders_status,status BETWEEN 0 AND 5;index:idx_orders_school_status_created,priority:2;index:idx_orders_runner_status,priority:2"`
	Deadline          time.Time   `gorm:"column:deadline;type:datetime(3);not null;index:idx_orders_deadline"`
	AcceptTime        *time.Time  `gorm:"column:accept_time;type:datetime(3)"`
	RunnerFinishedAt  *time.Time  `gorm:"column:runner_finished_at;type:datetime(3)"`
	ConfirmedAt       *time.Time  `gorm:"column:confirmed_at;type:datetime(3)"`
	AutoConfirmAt     *time.Time  `gorm:"column:auto_confirm_at;type:datetime(3);index:idx_orders_auto_confirm"`
	CanceledByID      *uint64     `gorm:"column:canceled_by_id;index:idx_orders_canceled_by"`
	CanceledByRole    *Role       `gorm:"column:canceled_by_role;type:tinyint unsigned;check:chk_orders_canceled_role,canceled_by_role BETWEEN 0 AND 2"`
	CanceledAt        *time.Time  `gorm:"column:canceled_at;type:datetime(3)"`
	CancelReason      string      `gorm:"column:cancel_reason;type:varchar(500);not null;default:''"`
	CreatedAt         time.Time   `gorm:"column:created_at;type:datetime(3);not null;index:idx_orders_school_status_created,priority:3;index:idx_orders_publisher_created,priority:2"`
	UpdatedAt         time.Time   `gorm:"column:updated_at;type:datetime(3);not null"`
}

// OrderStatusLogs 记录订单状态变更和操作者，作为履约审计事实来源。
type OrderStatusLogs struct {
	ID           uint64       `gorm:"column:id;primaryKey;autoIncrement"`
	OrderID      uint64       `gorm:"column:order_id;not null;index:idx_order_status_logs_order_created,priority:1"`
	FromStatus   *OrderStatus `gorm:"column:from_status;type:tinyint unsigned;check:chk_order_status_logs_from,from_status BETWEEN 0 AND 5"`
	ToStatus     OrderStatus  `gorm:"column:to_status;type:tinyint unsigned;not null;check:chk_order_status_logs_to,to_status BETWEEN 0 AND 5"`
	OperatorType ActorType    `gorm:"column:operator_type;type:varchar(16);not null;check:chk_order_status_logs_operator,operator_type IN ('publisher','runner','admin','system')"`
	OperatorID   *uint64      `gorm:"column:operator_id;index:idx_order_status_logs_operator"`
	Reason       string       `gorm:"column:reason;type:varchar(500);not null;default:''"`
	CreatedAt    time.Time    `gorm:"column:created_at;type:datetime(3);not null;index:idx_order_status_logs_order_created,priority:2"`
}

// OrderPayments 保存每个订单的线下付款及确认状态。
type OrderPayments struct {
	ID            uint64         `gorm:"column:id;primaryKey;autoIncrement"`
	OrderID       uint64         `gorm:"column:order_id;not null;uniqueIndex:uk_order_payments_order"`
	AmountCents   int64          `gorm:"column:amount_cents;not null;check:chk_order_payments_amount,amount_cents >= 0"`
	Method        *PaymentMethod `gorm:"column:method;type:tinyint unsigned;check:chk_order_payments_method,method BETWEEN 0 AND 3"`
	Status        PaymentStatus  `gorm:"column:status;type:tinyint unsigned;not null;default:0;check:chk_order_payments_status,status BETWEEN 0 AND 5;index:idx_order_payments_status"`
	ReportedByID  *uint64        `gorm:"column:reported_by_id;index:idx_order_payments_reported_by"`
	ReportedAt    *time.Time     `gorm:"column:reported_at;type:datetime(3)"`
	ConfirmedByID *uint64        `gorm:"column:confirmed_by_id;index:idx_order_payments_confirmed_by"`
	ConfirmedAt   *time.Time     `gorm:"column:confirmed_at;type:datetime(3)"`
	Note          string         `gorm:"column:note;type:varchar(500);not null;default:''"`
	CreatedAt     time.Time      `gorm:"column:created_at;type:datetime(3);not null"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;type:datetime(3);not null"`
}

// OrderPaymentProofs 保存付款截图等付款凭证。
type OrderPaymentProofs struct {
	ID           uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	PaymentID    uint64    `gorm:"column:payment_id;not null;index:idx_order_payment_proofs_payment"`
	UploadedByID *uint64   `gorm:"column:uploaded_by_id;index:idx_order_payment_proofs_uploader"`
	FileURL      string    `gorm:"column:file_url;type:varchar(512);not null"`
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime(3);not null"`
}

// UserViolation 将举报线索和审核确认后的违规处罚保存在同一案件表中。
type UserViolation struct {
	ID             uint64          `gorm:"column:id;primaryKey;autoIncrement"`
	OrderID        *uint64         `gorm:"column:order_id;index:idx_user_violations_order"`
	ReporterID     *uint64         `gorm:"column:reporter_id;index:idx_user_violations_reporter"`
	TargetUserID   uint64          `gorm:"column:target_user_id;not null;index:idx_user_violations_target_status,priority:1"`
	ReportType     string          `gorm:"column:report_type;type:varchar(50);not null"`
	Description    string          `gorm:"column:description;type:text;not null"`
	Status         ViolationStatus `gorm:"column:status;type:tinyint unsigned;not null;default:0;check:chk_user_violations_status,status BETWEEN 0 AND 3;index:idx_user_violations_target_status,priority:2"`
	ViolationType  *string         `gorm:"column:violation_type;type:varchar(50)"`
	PenaltyType    *PenaltyType    `gorm:"column:penalty_type;type:tinyint unsigned;check:chk_user_violations_penalty,penalty_type BETWEEN 0 AND 2"`
	PenaltyStartAt *time.Time      `gorm:"column:penalty_start_at;type:datetime(3)"`
	PenaltyEndAt   *time.Time      `gorm:"column:penalty_end_at;type:datetime(3)"`
	Decision       *string         `gorm:"column:decision;type:varchar(1000)"`
	HandledBy      *uint64         `gorm:"column:handled_by;index:idx_user_violations_handled_by"`
	ResolvedAt     *time.Time      `gorm:"column:resolved_at;type:datetime(3)"`
	CreatedAt      time.Time       `gorm:"column:created_at;type:datetime(3);not null"`
	UpdatedAt      time.Time       `gorm:"column:updated_at;type:datetime(3);not null"`
}
