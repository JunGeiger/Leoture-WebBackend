package errors

import "net/http"

// ResponseCode 业务响应码
//
// 编码规则:
//   - 1xxx: 成功 (Success)
//   - 4xxx: 客户端错误 (Client Error)  —— 请求本身有问题
//   - 5xxx: 服务端错误 (Server Error)  —— 服务端处理异常
//
// 设计原则:
//   - 每个错误码有唯一且明确的使用场景，不存在语义重叠
//   - 4xxx/5xxx 与 HTTP 状态码保持直觉对齐，便于前端和运维理解
//   - 各段预留空位，按业务增长按需扩展，不预定义无用码
type ResponseCode int

const (
	// ============================
	// 1xxx - 成功
	// ============================

	// CodeSuccess 通用操作成功
	// HTTP 200 | 场景: 查询、通用操作等无特殊语义的成功响应
	CodeSuccess ResponseCode = 1000

	// CodeCreated 资源创建成功
	// HTTP 201 | 场景: 创建用户、角色、菜单、权限等实体后返回
	CodeCreated ResponseCode = 1001

	// CodeUpdated 资源更新成功
	// HTTP 200 | 场景: 修改用户信息、更新角色权限、编辑菜单等
	CodeUpdated ResponseCode = 1002

	// CodeDeleted 资源删除成功
	// HTTP 200 | 场景: 删除用户、移除角色、清除权限绑定等
	CodeDeleted ResponseCode = 1003
)

const (
	// ============================
	// 4xxx - 客户端错误
	// ============================

	// CodeInvalidParam 请求参数校验失败
	// HTTP 400 | 场景: JSON 字段类型错误、必填字段缺失、值超出范围、格式不符等
	// 示例: 创建用户时 username 为空、password 长度不足、email 格式非法
	CodeInvalidParam ResponseCode = 4000

	// CodeUnauthenticated 未认证（身份凭证缺失或无效）
	// HTTP 401 | 场景: 未携带 Token、Token 格式非法、Token 签名校验失败、Token 已过期
	// 注意: 与 CodeTokenExpired 的区别 —— 本码用于 Token 本身无效，CodeTokenExpired 用于 Token 有效但已过期
	CodeUnauthenticated ResponseCode = 4001

	// CodeTokenExpired 令牌过期（需刷新）
	// HTTP 401 | 场景: JWT Access Token 过期，前端应使用 Refresh Token 换取新 Token
	// 前端处理: 收到此码后自动调用 /auth/refresh 接口，成功则重试原请求
	CodeTokenExpired ResponseCode = 4002

	// CodeForbidden 无权限访问
	// HTTP 403 | 场景: 用户已认证但 Casbin 策略判定无操作权限
	// 示例: 普通用户尝试删除角色、访问未授权菜单、执行越权操作
	CodeForbidden ResponseCode = 4003

	// CodeNotFound 资源不存在
	// HTTP 404 | 场景: 按 ID 查询用户/角色/菜单/权限时记录不存在，或路由对应的资源为空
	// 示例: GET /api/v1/users/999 返回 "用户不存在"
	CodeNotFound ResponseCode = 4004

	// CodeAlreadyExists 资源已存在（唯一性冲突）
	// HTTP 409 | 场景: 创建时违反唯一约束
	// 示例: 注册用户名已存在、角色名称重复、菜单路径冲突
	CodeAlreadyExists ResponseCode = 4009

	// CodeTooManyRequests 请求过于频繁（限流）
	// HTTP 429 | 场景: 登录失败次数过多触发锁定、接口调用频率超限
	// 前端处理: 提示用户稍后重试，可根据 Retry-After 响应头计算等待时间
	CodeTooManyRequests ResponseCode = 4029
)

const (
	// ============================
	// 5xxx - 服务端错误
	// ============================

	// CodeInternal 服务器内部错误
	// HTTP 500 | 场景: 未归类的服务端异常（兜底错误码）
	// 注意: 生产环境不应将内部错误细节暴露给客户端，详细错误仅写入日志
	CodeInternal ResponseCode = 5000

	// CodeServiceUnavailable 服务暂不可用
	// HTTP 503 | 场景: 系统维护、服务降级、依赖组件未就绪
	CodeServiceUnavailable ResponseCode = 5003

	// CodeTimeout 请求超时
	// HTTP 504 | 场景: 数据库查询超时、外部 API 调用超时、长时间运算超时
	CodeTimeout ResponseCode = 5004

	// CodeDatabaseError 数据库操作异常
	// HTTP 500 | 场景: 连接池耗尽、SQL 执行异常、事务提交失败、迁移失败
	// 注意: 与 CodeInternal 区别 —— 本码明确标识为数据库层故障，便于运维快速定位
	CodeDatabaseError ResponseCode = 5005

	// CodeCacheError 缓存操作异常
	// HTTP 500 | 场景: Redis 连接失败、序列化/反序列化异常、缓存雪崩触发降级
	// 注意: 缓存故障时业务应降级到数据库直查，而非直接返回错误
	CodeCacheError ResponseCode = 5006
)

// HTTPStatus 将业务码映射为 HTTP 状态码
func (c ResponseCode) HTTPStatus() int {
	switch c {
	case CodeSuccess, CodeUpdated, CodeDeleted:
		return http.StatusOK
	case CodeCreated:
		return http.StatusCreated
	case CodeInvalidParam:
		return http.StatusBadRequest
	case CodeUnauthenticated, CodeTokenExpired:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	case CodeAlreadyExists:
		return http.StatusConflict
	case CodeTooManyRequests:
		return http.StatusTooManyRequests
	case CodeServiceUnavailable:
		return http.StatusServiceUnavailable
	case CodeTimeout:
		return http.StatusGatewayTimeout
	default:
		return http.StatusInternalServerError
	}
}
