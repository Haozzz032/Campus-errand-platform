# 校园跑腿平台 API 接口文档

> 版本：v1.0  
> 服务地址：`http://localhost:8080`  
> 数据格式：`application/json; charset=utf-8`  
> 金额单位：分（例如 `350` 表示 3.50 元）  
> 时间格式：RFC3339，例如 `2026-11-15T20:00:00+08:00`

## 1. 当前实现说明

当前项目的 `main.go` 已建立 Gin 服务、JWT 认证中间件，以及 `/api/user` 认证路由组，但尚未注册业务路由和 handler。本文档是根据当前 `model/model.go` 与 `model/init.sql` 制定的 API 契约，路由实现时应以本文档为准。

JWT 约定：

```http
Authorization: Bearer <access_token>
```

JWT 默认有效期为 7 天。封禁、停用账号不能登录，也不能调用需要登录的接口。

## 2. 通用响应格式

成功：

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

失败：

```json
{
  "code": 40001,
  "message": "参数错误",
  "data": null
}
```

分页接口的 `data` 格式：

```json
{
  "items": [],
  "total": 0,
  "page": 1,
  "page_size": 20
}
```

通用分页参数：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---:|---:|---|
| `page` | integer | 否 | 1 | 页码，从 1 开始 |
| `page_size` | integer | 否 | 20 | 每页数量，1-100 |

## 3. 认证与角色

| 角色值 | 中文含义 | 权限概述 |
|---|---|---|
| `publisher` | 发布者 | 发布订单、确认完成、评价跑者、上报付款 |
| `runner` | 跑者 | 浏览可接订单、抢单、履约、确认收款、评价发布者 |
| `admin` | 管理员 | 学校、分类、用户、订单、举报案件管理 |

接口表中的“登录”列：

- `否`：公开接口，不需要 JWT。
- `是`：必须携带有效 JWT。
- `管理员`：必须携带管理员 JWT。

## 4. 状态枚举

### 4.1 账号状态

| API 值 | 数据库值 | 含义 |
|---|---:|---|
| `active` | 0 | 正常 |
| `banned` | 1 | 已封禁 |
| `disabled` | 2 | 已停用 |

### 4.2 订单状态

| API 值 | 数据库值 | 含义 |
|---|---:|---|
| `published` | 0 | 已发布，等待接单 |
| `accepted` | 1 | 已接单 |
| `delivering` | 2 | 配送中 |
| `arrived` | 3 | 已送达，等待发布者确认 |
| `confirmed` | 4 | 已确认完成 |
| `canceled` | 5 | 已取消 |

### 4.3 付款状态

| API 值 | 数据库值 | 含义 |
|---|---:|---|
| `pending` | 0 | 待上报 |
| `reported` | 1 | 用户已上报，待跑者确认 |
| `confirmed` | 2 | 跑者已确认收款 |
| `disputed` | 3 | 付款存在争议 |
| `refunded` | 4 | 已退款 |
| `canceled` | 5 | 已取消 |

### 4.4 举报/违规案件状态

| API 值 | 数据库值 | 含义 |
|---|---:|---|
| `reviewing` | 0 | 审核中 |
| `confirmed` | 1 | 已确认违规 |
| `rejected` | 2 | 已驳回 |
| `closed` | 3 | 已关闭 |

## 5. 共享返回对象

### 5.1 School

```json
{
  "school_id": 1,
  "code": "CAMPUS-A",
  "name": "示例大学",
  "status": true,
  "created_at": "2026-11-15T10:00:00+08:00"
}
```

### 5.2 UserSummary

```json
{
  "id": 10,
  "username": "student01",
  "role": "publisher",
  "school": {"school_id": 1, "name": "示例大学"},
  "status": "active",
  "violation_times": 0,
  "created_at": "2026-11-15T10:00:00+08:00"
}
```

### 5.3 UserSelf

在 `UserSummary` 基础上增加手机号和封禁信息：

```json
{
  "id": 10,
  "username": "student01",
  "role": "publisher",
  "school_id": 1,
  "phone": "13800000000",
  "status": "active",
  "ban_reason": "",
  "banned_until": null,
  "violation_times": 0,
  "created_at": "2026-11-15T10:00:00+08:00"
}
```

### 5.4 RunnerProfile

```json
{
  "user": {
    "id": 20,
    "username": "runner01",
    "role": "runner",
    "school": {"school_id": 1, "name": "示例大学"},
    "status": "active"
  },
  "completed_order_count": 12,
  "positive_count": 10,
  "negative_count": 1,
  "recent_ratings": []
}
```

手机号、密码哈希不通过公开跑者资料返回。

### 5.5 OrderCategory

```json
{
  "id": 3,
  "school_id": 1,
  "code": "express_pickup",
  "name": "快递代取",
  "description": "从驿站取件并送达",
  "icon_url": "",
  "sort_order": 10,
  "is_active": true
}
```

### 5.6 OrderSummary

```json
{
  "order_no": "CE202611150001",
  "school_id": 1,
  "category": {"id": 3, "name": "快递代取"},
  "title": "取一个快递",
  "publisher": {"id": 10, "username": "student01"},
  "runner": null,
  "remuneration_cents": 350,
  "start_place": "东区菜鸟驿站",
  "end_place": "6 号宿舍楼",
  "status": "published",
  "deadline": "2026-11-15T20:00:00+08:00",
  "created_at": "2026-11-15T16:00:00+08:00"
}
```

### 5.7 OrderDetail

在 `OrderSummary` 基础上增加：

```json
{
  "description": "取件码将在接单后提供",
  "contact_name": "张同学",
  "user_phone": "13800000000",
  "accept_time": null,
  "runner_finished_at": null,
  "auto_confirm_at": null,
  "confirmed_at": null,
  "canceled_by_id": null,
  "canceled_by_role": null,
  "canceled_at": null,
  "cancel_reason": "",
  "allowed_actions": ["accept"],
  "payment": null
}
```

`allowed_actions` 由当前用户角色和订单状态计算，可能包含：

`accept`、`start_delivery`、`arrive`、`confirm`、`cancel`、`report`、`report_payment`、`confirm_payment`、`rate`、`view_logs`。

隐私规则：

- `user_phone` 仅对订单发布者、已接单跑者和管理员返回。
- 未接单跑者查看可接订单时不返回 `user_phone` 和 `contact_name`。
- `password_hash` 永远不返回。

### 5.8 StatusLog

```json
{
  "id": 1,
  "from_status": "published",
  "to_status": "accepted",
  "operator_type": "runner",
  "operator_id": 20,
  "reason": "",
  "created_at": "2026-11-15T16:05:00+08:00"
}
```

### 5.9 Payment

```json
{
  "order_no": "CE202611150001",
  "amount_cents": 350,
  "method": "wechat",
  "status": "reported",
  "reported_by_id": 10,
  "reported_at": "2026-11-15T18:00:00+08:00",
  "confirmed_by_id": null,
  "confirmed_at": null,
  "note": "已微信转账",
  // "proofs": [
  //   {"id": 1, "file_url": "/uploads/payments/1.jpg", "uploaded_by_id": 10}
  // ],
  "created_at": "2026-11-15T18:00:00+08:00",
  "updated_at": "2026-11-15T18:00:00+08:00"
}
```

### 5.10 Rating

```json
{
  "id": 1,
  "order_no": "CE202611150001",
  "reviewer_id": 10,
  "target_user_id": 20,
  "is_positive": true,
  "content": "送达很快",
  "created_at": "2026-11-15T19:00:00+08:00"
}
```

### 5.11 ReportCase

```json
{
  "id": 5,
  "order_no": "CE202611150001",
  "reporter_id": 10,
  "target_user_id": 20,
  "report_type": "payment_dispute",
  "description": "已完成付款但对方表示未收到",
  "status": "reviewing",
  "violation_type": null,
  "penalty_type": null,
  "penalty_start_at": null,
  "penalty_end_at": null,
  "decision": null,
  "handled_by": null,
  "resolved_at": null,
  "created_at": "2026-11-15T19:10:00+08:00"
}
```

## 6. 接口总览

| # | 方法 | 路径 | 作用 | 登录 |
|---:|---|---|---|---|
| 1 | POST | `/api/auth/register` | 注册发布者或跑者 | 否 |
| 2 | POST | `/api/auth/login` | 登录并获取 JWT | 否 |
| 3 | GET | `/api/schools` | 查询启用学校 | 否 |
| 4 | GET | `/api/order-categories` | 查询学校订单分类 | 否 |
| 5 | GET | `/api/runners/{user_id}` | 查看跑者公开资料与统计 | 是 |
| 6 | GET | `/api/user/profile` | 获取自己的账号资料 | 是 |
| 7 | PUT | `/api/user/profile` | 修改自己的用户名或手机号 | 是 |
| 8 | PUT | `/api/user/password` | 修改密码 | 是 |
| 9 | POST | `/api/orders` | 发布订单 | 是（发布者） |
| 10 | GET | `/api/orders/available` | 查询本校可接订单 | 是（跑者） |
| 11 | GET | `/api/orders/mine/published` | 查询自己发布的订单 | 是（发布者） |
| 12 | GET | `/api/orders/mine/accepted` | 查询自己接取的订单 | 是（跑者） |
| 13 | GET | `/api/orders/{order_no}` | 获取订单详情 | 是 |
| 14 | POST | `/api/orders/{order_no}/accept` | 抢单 | 是（跑者） |
| 15 | POST | `/api/orders/{order_no}/start-delivery` | 开始配送 | 是（接单跑者） |
| 16 | POST | `/api/orders/{order_no}/arrive` | 标记已送达 | 是（接单跑者） |
| 17 | POST | `/api/orders/{order_no}/confirm` | 确认订单完成 | 是（发布者） |
| 18 | POST | `/api/orders/{order_no}/cancel` | 取消订单 | 是（订单双方或管理员） |
| 19 | GET | `/api/orders/{order_no}/status-logs` | 查询订单状态日志 | 是（订单双方或管理员） |
| 20 | POST | `/api/orders/{order_no}/rating` | 评价订单对方 | 是（订单双方） |
| 21 | GET | `/api/orders/{order_no}/payment` | 获取付款记录 | 是（订单双方或管理员） |
| 22 | POST | `/api/orders/{order_no}/payment/report` | 上报线下付款及凭证 | 是（发布者） |
| 23 | POST | `/api/orders/{order_no}/payment/confirm` | 跑者确认收款 | 是（接单跑者） |
| 24 | POST | `/api/orders/{order_no}/payment/dispute` | 发起付款争议 | 是（订单双方） |
| 25 | POST | `/api/reports` | 举报用户或订单违规行为 | 是 |
| 26 | GET | `/api/reports/mine` | 查询自己提交的举报 | 是 |
| 27 | GET | `/api/reports/{id}` | 查询举报案件 | 是（举报人或管理员） |
| 28 | GET | `/api/admin/users` | 查询用户列表 | 管理员 |
| 29 | GET | `/api/admin/users/{id}` | 查询用户管理详情 | 管理员 |
| 30 | PATCH | `/api/admin/users/{id}/status` | 封禁、解封或停用账号 | 管理员 |
| 31 | GET | `/api/admin/orders` | 查询全部订单 | 管理员 |
| 32 | POST | `/api/admin/orders/{order_no}/cancel` | 管理员强制取消订单 | 管理员 |
| 33 | GET | `/api/admin/reports` | 查询全部举报案件 | 管理员 |
| 34 | PATCH | `/api/admin/reports/{id}/review` | 审核举报并处罚用户 | 管理员 |
| 35 | POST | `/api/admin/schools` | 新增学校 | 管理员 |
| 36 | PUT | `/api/admin/schools/{id}` | 修改学校信息或状态 | 管理员 |
| 37 | POST | `/api/admin/order-categories` | 新增订单分类 | 管理员 |
| 38 | PUT | `/api/admin/order-categories/{id}` | 修改订单分类或状态 | 管理员 |
---

## 7. 详细接口说明

### 1. 注册账号

`POST /api/auth/register`

- 功能：注册发布者或跑者账号，成功后直接返回 JWT。
- 登录：否。
- 角色：仅允许注册 `publisher` 或 `runner`；管理员不能通过公开接口注册。

请求体：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---:|---|
| `username` | string | 是 | 3-50 个字符，全平台唯一 |
| `password` | string | 是 | 8-72 个字符，服务端只保存哈希 |
| `role` | string | 是 | `publisher` 或 `runner` |
| `school_id` | integer | 是 | 已启用的学校 ID |
| `phone` | string | 是 | 手机号，用于订单联系 |

返回 `data`：

```json
{
  "access_token": "<jwt>",
  "token_type": "Bearer",
  "expires_in": 604800,
  "user": {}
}
```

`user` 为 `UserSelf`。重复用户名返回 `40901`。

### 2. 登录

`POST /api/auth/login`

- 功能：用户名和密码登录，签发 JWT。
- 登录：否。
- 角色：不限。

请求体：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---:|---|
| `username` | string | 是 | 用户名 |
| `password` | string | 是 | 明文密码，仅用于校验 |

返回 `data` 与注册接口相同。

账号被封禁或停用时返回 `40302`；用户名或密码错误统一返回 `40102`，不暴露账号是否存在。

### 3. 查询学校

`GET /api/schools`

- 功能：查询平台启用的学校，用于注册和筛选。
- 登录：否。

查询参数：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---:|---|
| `keyword` | string | 否 | 学校名称或代码模糊搜索 |
| `page` | integer | 否 | 页码 |
| `page_size` | integer | 否 | 每页数量 |

返回 `data`：`School` 的分页对象。

### 4. 查询订单分类

`GET /api/order-categories`

- 功能：查询指定学校可用的订单分类。
- 登录：否。

查询参数：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---:|---|
| `school_id` | integer | 是 | 学校 ID |
| `active_only` | boolean | 否 | 默认 `true`，仅返回启用分类 |

返回 `data`：`OrderCategory` 数组，按 `sort_order` 升序。

### 5. 查询跑者公开资料

`GET /api/runners/{user_id}`

- 功能：查看跑者统计、最近评价和公开账号信息。
- 登录：否。

路径参数：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---:|---|
| `user_id` | integer | 是 | 跑者用户 ID |

返回 `data`：`RunnerProfile`。目标用户必须是 `runner`；不返回手机号、封禁原因和密码信息。

### 6. 获取自己的资料

`GET /api/user/profile`

- 功能：获取当前登录用户资料和角色统计。
- 登录：是。
- 角色：不限。

返回 `data`：`UserSelf`。跑者或发布者可同时返回对应 profile 统计。

### 7. 修改自己的资料

`PUT /api/user/profile`

- 功能：修改用户名或手机号。
- 登录：是。
- 角色：不限。

请求体：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---:|---|
| `username` | string | 否 | 修改后仍需全平台唯一 |
| `phone` | string | 否 | 新手机号 |

返回 `data`：修改后的 `UserSelf`。角色、学校、账号状态不能通过该接口修改。

### 8. 修改密码

`PUT /api/user/password`

- 功能：验证旧密码后修改登录密码。
- 登录：是。
- 角色：不限。

请求体：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---:|---|
| `old_password` | string | 是 | 当前密码 |
| `new_password` | string | 是 | 8-72 个字符 |

成功返回 `data: null`。旧密码错误返回 `40002`。修改成功后建议客户端清除旧 JWT 并重新登录。

### 9. 发布订单

`POST /api/orders`

- 功能：发布者创建固定报酬的跑腿订单。
- 登录：是。
- 角色：`publisher`。

请求体：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---:|---|
| `category_id` | integer | 是 | 与发布者同校且启用的分类 |
| `title` | string | 是 | 1-100 个字符 |
| `description` | string | 是 | 订单详细说明 |
| `remuneration_cents` | integer | 是 | 固定报酬，单位分，不能小于 0 |
| `start_place` | string | 是 | 取件地点 |
| `end_place` | string | 是 | 送达地点 |
| `contact_name` | string | 是 | 订单联系人 |
| `user_phone` | string | 是 | 本次订单联系电话 |
| `deadline` | string | 是 | RFC3339 时间，必须晚于当前时间 |

返回 `data`：新建的 `OrderDetail`，状态为 `published`，HTTP 状态码建议为 `201`。

服务端生成 `order_no`，并根据当前用户确定 `school_id` 和 `publisher_id`。

### 10. 查询可接订单

`GET /api/orders/available`

- 功能：跑者查询本校处于 `published` 状态的可接订单。
- 登录：是。
- 角色：`runner`。

查询参数：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---:|---|
| `category_id` | integer | 否 | 分类筛选 |
| `keyword` | string | 否 | 标题和描述搜索 |
| `sort` | string | 否 | `newest`、`reward_desc`、`deadline_asc` |
| `page` | integer | 否 | 页码 |
| `page_size` | integer | 否 | 每页数量 |

返回 `data`：`OrderSummary` 分页对象。不得返回联系电话和联系人。

### 11. 查询自己发布的订单

`GET /api/orders/mine/published`

- 功能：发布者查询自己发布的订单。
- 登录：是。
- 角色：`publisher`。

查询参数：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---:|---|
| `status` | string | 否 | 按订单状态筛选 |
| `page` | integer | 否 | 页码 |
| `page_size` | integer | 否 | 每页数量 |

返回 `data`：`OrderSummary` 分页对象，按创建时间倒序。

### 12. 查询自己接取的订单

`GET /api/orders/mine/accepted`

- 功能：跑者查询自己接取过的订单。
- 登录：是。
- 角色：`runner`。

查询参数与“查询自己发布的订单”相同。

返回 `data`：`OrderSummary` 分页对象，按接单时间或创建时间倒序。

### 13. 获取订单详情

`GET /api/orders/{order_no}`

- 功能：获取订单详情、当前用户可执行操作和付款摘要。
- 登录：是。
- 角色：不限；未参与订单的普通用户只能看到脱敏详情。

路径参数：`order_no` 为业务订单号。

返回 `data`：`OrderDetail`。

详情权限：

- 发布者和接单跑者可以看到联系电话。
- 管理员可以看到完整信息。
- 其他登录用户可以用该接口查看仍可供接单的订单摘要，但不返回联系人、电话和内部日志。

### 14. 抢单

`POST /api/orders/{order_no}/accept`

- 功能：跑者以先到先得方式接单。
- 登录：是。
- 角色：`runner`。

请求体：无。

返回 `data`：更新后的 `OrderDetail`，状态为 `accepted`。

业务规则：

- 跑者账号必须是 `active`。
- 跑者学校必须与订单学校一致。
- 仅 `published` 且 `runner_id` 为空的订单可接。
- 必须使用条件更新或行锁保证并发下只有一个跑者成功。
- 失败返回 `40902`（订单已被接走）或 `42202`（学校不匹配）。

### 15. 开始配送

`POST /api/orders/{order_no}/start-delivery`

- 功能：接单跑者将订单从 `accepted` 更新为 `delivering`。
- 登录：是。
- 角色：接单跑者。

返回 `data`：更新后的 `OrderDetail`。只有订单当前跑者可以调用，否则返回 `40301`。

### 16. 标记已送达

`POST /api/orders/{order_no}/arrive`

- 功能：跑者标记商品或物品已送达，进入等待发布者确认状态。
- 登录：是。
- 角色：接单跑者。

返回 `data`：更新后的 `OrderDetail`，状态为 `arrived`。

服务端同时写入：

- `runner_finished_at`
- `auto_confirm_at`，默认设为当前时间加 24 小时

### 17. 确认订单完成

`POST /api/orders/{order_no}/confirm`

- 功能：发布者确认收到物品并完成订单；系统也会在超时后自动确认。
- 登录：是。
- 角色：订单发布者。

返回 `data`：更新后的 `OrderDetail`，状态为 `confirmed`。

只有 `arrived` 状态可以确认。系统自动确认以 `operator_type=system` 写入状态日志。

### 18. 取消订单

`POST /api/orders/{order_no}/cancel`

- 功能：取消尚未完成的订单并记录取消人、角色和原因。
- 登录：是。
- 角色：订单发布者、接单跑者或管理员。

请求体：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---:|---|
| `reason` | string | 是 | 1-500 个字符 |

返回 `data`：更新后的 `OrderDetail`，状态为 `canceled`。

规则：

- `published` 状态由发布者取消。
- `accepted`、`delivering`、`arrived` 状态由发布者或接单跑者取消。
- `confirmed` 状态不能取消，只能举报。
- 管理员可以强制取消未确认订单。
- 每次取消都必须写 `order_status_logs`。

### 19. 查询订单状态日志

`GET /api/orders/{order_no}/status-logs`

- 功能：查看订单状态变更历史。
- 登录：是。
- 角色：发布者、接单跑者或管理员。

返回 `data`：`StatusLog` 数组，按创建时间升序。

### 20. 评价订单对方

`POST /api/orders/{order_no}/rating`

- 功能：订单完成后，发布者评价跑者，或跑者评价发布者。
- 登录：是。
- 角色：订单双方。

请求体：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---:|---|
| `is_positive` | boolean | 是 | `true` 为好评，`false` 为差评 |
| `content` | string | 否 | 最多 500 个字符 |

返回 `data`：`Rating`。

规则：

- 只有 `confirmed` 订单可以评价。
- 发布者和跑者各只能评价一次。
- 同一评价方向重复提交返回 `40903`。
- 评价后事务更新对应用户 profile 的好评/差评计数。
### 21. 获取付款记录

`GET /api/orders/{order_no}/payment`

- 功能：查询订单的线下付款状态和凭证。
- 登录：是。
- 角色：订单发布者、接单跑者或管理员。

返回 `data`：`Payment`。未创建付款记录时返回 `data: null`。

### 22. 上报线下付款

`POST /api/orders/{order_no}/payment/report`

- 功能：发布者上报已完成线下付款，并提交付款凭证。
- 登录：是。
- 角色：订单发布者。

请求体：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---:|---|
| `amount_cents` | integer | 是 | 实付金额，单位分，必须等于订单报酬 |
| `method` | string | 是 | `wechat`、`alipay`、`cash`、`other` |
| `note` | string | 否 | 最多 500 个字符 |
| `proof_urls` | string[] | 是 | 至少 1 个、最多 5 个付款凭证地址 |

返回 `data`：`Payment`，状态为 `reported`。

规则：

- 订单必须已经接单，且尚未取消。
- 仅发布者可以上报。
- 相同订单重复上报时，仅允许在 `pending` 或 `reported` 状态补充或替换凭证。
- 付款凭证由独立文件上传服务生成 URL，业务接口只保存 URL。

### 23. 跑者确认收款

`POST /api/orders/{order_no}/payment/confirm`

- 功能：接单跑者确认已经收到线下付款。
- 登录：是。
- 角色：接单跑者。

请求体：无。

返回 `data`：`Payment`，状态为 `confirmed`，同时写入确认人和确认时间。

只有 `reported` 状态可以确认；不能确认自己上报的付款。

### 24. 发起付款争议

`POST /api/orders/{order_no}/payment/dispute`

- 功能：订单双方对付款是否完成、金额或凭证有争议时发起审核案件。
- 登录：是。
- 角色：订单双方。

请求体：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---:|---|
| `reason` | string | 是 | 1-500 个字符 |
| `proof_urls` | string[] | 否 | 额外说明材料地址 |

返回 `data`：

```json
{
  "payment": {},
  "report": {}
}
```

付款状态更新为 `disputed`，同时创建状态为 `reviewing` 的 `ReportCase`，目标用户自动设为订单对方。

### 25. 提交举报

`POST /api/reports`

- 功能：举报用户违规行为，可关联订单。
- 登录：是。
- 角色：不限。

请求体：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---:|---|
| `target_user_id` | integer | 是 | 被举报用户 ID |
| `order_no` | string | 否 | 关联订单号 |
| `report_type` | string | 是 | 例如 `payment_dispute`、`malicious_cancel`、`fraud`、`abuse`、`other` |
| `description` | string | 是 | 1-5000 个字符 |

返回 `data`：`ReportCase`，初始状态为 `reviewing`。

规则：

- 不能举报自己。
- 关联订单时，举报人和被举报人必须是订单双方。
- 不能重复提交完全相同的审核中案件。
- 管理员的审核结果才代表违规事实，`reviewing` 不等于已认定违规。

### 26. 查询自己提交的举报

`GET /api/reports/mine`

- 功能：查看当前用户提交过的举报案件。
- 登录：是。

查询参数：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---:|---|
| `status` | string | 否 | 案件状态 |
| `page` | integer | 否 | 页码 |
| `page_size` | integer | 否 | 每页数量 |

返回 `data`：`ReportCase` 分页对象。

### 27. 查询举报案件详情

`GET /api/reports/{id}`

- 功能：举报人查看自己的案件，或管理员查看任意案件。
- 登录：是。
- 角色：举报人或管理员。

返回 `data`：`ReportCase`。其他普通用户访问返回 `40301`。
### 28. 管理员查询用户

`GET /api/admin/users`

- 功能：分页查询用户账号和状态。
- 登录：管理员。

查询参数：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---:|---|
| `keyword` | string | 否 | 用户名或手机号 |
| `role` | string | 否 | `publisher`、`runner`、`admin` |
| `status` | string | 否 | `active`、`banned`、`disabled` |
| `school_id` | integer | 否 | 学校筛选 |
| `page` | integer | 否 | 页码 |
| `page_size` | integer | 否 | 每页数量 |

返回 `data`：`UserSummary` 分页对象。手机号仅管理员可见。

### 29. 管理员查询用户详情

`GET /api/admin/users/{id}`

- 功能：查看用户账号、统计数据、最近订单和违规案件。
- 登录：管理员。

返回 `data`：

```json
{
  "user": {},
  "profile": {},
  "recent_orders": [],
  "recent_reports": []
}
```

### 30. 修改用户账号状态

`PATCH /api/admin/users/{id}/status`

- 功能：封禁、解封或停用账号。
- 登录：管理员。

请求体：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---:|---|
| `status` | string | 是 | `active`、`banned`、`disabled` |
| `reason` | string | 条件必填 | 封禁或停用时必填 |
| `banned_until` | string/null | 否 | 封禁到期时间；`null` 表示永久封禁 |

返回 `data`：更新后的 `UserSelf`。

规则：

- 管理员不能封禁自己。
- 设置 `active` 时清空封禁原因、封禁时间和封禁截止时间。
- 处罚应同时写入 `user_violations` 或关联已有案件，不能只有用户状态变化。

### 31. 管理员查询全部订单

`GET /api/admin/orders`

- 功能：按学校、状态、发布者、跑者和时间查询订单。
- 登录：管理员。

查询参数：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---:|---|
| `school_id` | integer | 否 | 学校 |
| `status` | string | 否 | 订单状态 |
| `publisher_id` | integer | 否 | 发布者 |
| `runner_id` | integer | 否 | 跑者 |
| `created_from` | string | 否 | 开始时间 |
| `created_to` | string | 否 | 结束时间 |
| `page` | integer | 否 | 页码 |
| `page_size` | integer | 否 | 每页数量 |

返回 `data`：`OrderSummary` 分页对象。管理员可以继续调用订单详情接口获取完整信息。

### 32. 管理员强制取消订单

`POST /api/admin/orders/{order_no}/cancel`

- 功能：管理员处理无法正常履约的订单。
- 登录：管理员。

请求体：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---:|---|
| `reason` | string | 是 | 管理员取消原因 |

返回 `data`：更新后的 `OrderDetail`，状态为 `canceled`。

已完成订单不能取消，只能通过举报案件处理。

### 33. 管理员查询举报案件

`GET /api/admin/reports`

- 功能：查询全部举报和付款争议案件。
- 登录：管理员。

查询参数：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---:|---|
| `status` | string | 否 | 案件状态 |
| `report_type` | string | 否 | 举报类型 |
| `target_user_id` | integer | 否 | 被举报人 |
| `order_no` | string | 否 | 关联订单 |
| `page` | integer | 否 | 页码 |
| `page_size` | integer | 否 | 每页数量 |

返回 `data`：`ReportCase` 分页对象。

### 34. 审核举报案件

`PATCH /api/admin/reports/{id}/review`

- 功能：管理员确认或驳回举报，并给出处罚结论。
- 登录：管理员。

请求体：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---:|---|
| `status` | string | 是 | `confirmed`、`rejected` 或 `closed` |
| `violation_type` | string | 条件必填 | `confirmed` 时必填 |
| `penalty_type` | string | 条件必填 | `warning`、`suspend`、`ban` |
| `penalty_start_at` | string | 否 | 处罚开始时间，默认当前时间 |
| `penalty_end_at` | string/null | 否 | 处罚结束时间；永久封禁传 `null` |
| `decision` | string | 是 | 管理员处理说明 |

返回 `data`：更新后的 `ReportCase`。

事务要求：

- 仅 `reviewing` 状态的案件可以审核。
- `confirmed` 时增加目标用户的 `violation_times`。
- `suspend` 或 `ban` 同时更新 `users.status`、`ban_reason`、`banned_at`、`banned_until`。
- `rejected` 不增加违规次数，也不改变用户状态。
- 所有操作必须可追溯到处理管理员。
### 35. 新增学校

`POST /api/admin/schools`

- 功能：创建学校配置。
- 登录：管理员。

请求体：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---:|---|
| `code` | string | 是 | 唯一学校代码，最多 32 字符 |
| `name` | string | 是 | 学校名称，最多 100 字符 |
| `status` | boolean | 否 | 默认启用 |

返回 `data`：`School`。

### 36. 修改学校

`PUT /api/admin/schools/{id}`

- 功能：修改学校名称或启用状态。
- 登录：管理员。

请求体：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---:|---|
| `name` | string | 否 | 学校名称 |
| `status` | boolean | 否 | 是否启用 |

返回 `data`：修改后的 `School`。学校不物理删除，停用后不能注册新用户或发布新订单。

### 37. 新增订单分类

`POST /api/admin/order-categories`

- 功能：为指定学校新增订单分类。
- 登录：管理员。

请求体：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---:|---|
| `school_id` | integer | 是 | 所属学校 |
| `code` | string | 是 | 学校内唯一的机器编码 |
| `name` | string | 是 | 显示名称 |
| `description` | string | 否 | 分类说明 |
| `icon_url` | string | 否 | 图标地址 |
| `sort_order` | integer | 否 | 排序值，默认 0 |
| `is_active` | boolean | 否 | 默认启用 |

返回 `data`：`OrderCategory`。同一学校内 `code` 不能重复。

### 38. 修改订单分类

`PUT /api/admin/order-categories/{id}`

- 功能：修改分类显示信息或启用状态。
- 登录：管理员。

请求体：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---:|---|
| `name` | string | 否 | 显示名称 |
| `description` | string | 否 | 分类说明 |
| `icon_url` | string | 否 | 图标地址 |
| `sort_order` | integer | 否 | 排序值 |
| `is_active` | boolean | 否 | 是否启用 |

返回 `data`：修改后的 `OrderCategory`。分类不物理删除，已有订单继续保留 `category_id`。

---

## 8. 订单状态机

| 操作 | 操作者 | 原状态 | 新状态 |
|---|---|---|---|
| 发布订单 | 发布者 | - | `published` |
| 抢单 | 跑者 | `published` | `accepted` |
| 开始配送 | 接单跑者 | `accepted` | `delivering` |
| 标记送达 | 接单跑者 | `delivering` | `arrived` |
| 发布者确认 | 发布者 | `arrived` | `confirmed` |
| 超时自动确认 | 系统 | `arrived` | `confirmed` |
| 发布者取消 | 发布者 | `published` | `canceled` |
| 双方取消 | 发布者/跑者 | `accepted`、`delivering`、`arrived` | `canceled` |
| 管理员强制取消 | 管理员 | `published`、`accepted`、`delivering`、`arrived` | `canceled` |

每次状态变更必须：

1. 使用条件更新或事务锁校验当前状态。
2. 在同一个事务中更新 `orders`。
3. 写入 `order_status_logs`。
4. 对跑者完成订单计数和评分缓存进行可重算更新。

## 9. 通用错误码

| code | HTTP 状态 | 含义 |
|---:|---:|---|
| `0` | 200/201 | 成功 |
| `40001` | 400 | 参数格式或字段校验失败 |
| `40002` | 400 | 旧密码错误或业务输入错误 |
| `40101` | 401 | 未携带 JWT |
| `40102` | 401 | 用户名或密码错误 |
| `40103` | 401 | JWT 无效或已过期 |
| `40301` | 403 | 当前角色无权限执行操作 |
| `40302` | 403 | 账号被封禁或停用 |
| `40401` | 404 | 学校、用户、订单、付款或举报不存在 |
| `40901` | 409 | 用户名或学校代码等唯一字段冲突 |
| `40902` | 409 | 订单已被其他跑者接走 |
| `40903` | 409 | 重复评价或重复举报 |
| `42201` | 422 | 当前订单状态不允许该操作 |
| `42202` | 422 | 用户学校与订单学校不匹配 |
| `42203` | 422 | 付款金额、付款状态或凭证不合法 |
| `42901` | 429 | 请求过于频繁 |
| `50001` | 500 | 服务器内部错误 |

## 10. 实现路由时的注意事项

建议路由分组：

```text
/api/auth               公开登录注册
/api/schools            公开学校查询
/api/order-categories   公开分类查询
/api/runners            公开跑者资料
/api/user               JWT 登录用户
/api/orders             JWT 登录用户
/api/reports            JWT 登录用户
/api/admin              JWT + 管理员角色
```

当前只有 `middleware.Auth` 和 `middleware.OptionalAuth`，还没有管理员角色中间件。实现时应新增 `RequireRole(model.RoleAdmin)` 等权限校验，不能仅依赖前端隐藏按钮。

统一实现要求：

- 所有接口只返回 DTO，不直接返回 GORM 模型。
- `password_hash` 永远不出现于任何响应。
- 手机号仅向订单双方和管理员展示。
- 所有写操作必须校验账号状态、角色、学校归属和资源归属。
- 金额始终使用整数分，不使用浮点数。
- 时间统一使用 RFC3339 输出。
- 分页上限固定为 100，避免一次请求拉取全部数据。
- 订单抢单、状态流转、评价和审核必须在事务中完成。