# LeotureWeb项目技术选型决策说明

---

## viper v1.21.0作为配置管理工具
- Viper 是目前 Go Web 生态里唯一能做到“一套 API 通吃文件、ENV、Flag、远程配置”的方案。
- 使用 fsnotify v1.10.1 作为热更新监听器
- 优点：
  - 生态对齐：Gin / Cobra / Swag / Casbin 社区示例几乎全用 Viper，复制即用。
  - 运维友好：config.yaml+ APP_xxx环境变量覆盖，天然适配 Docker / K8s。
  - 演进平滑：后期需要 etcd / Consul / Apollo，只需换 Provider，不动上层代码。
  - 开发效率高：Unmarshal(&cfg)一行完成配置组装，心智负担低。
- 缺点：
  - 隐式失败：字段缺失不报错，靠默认值兜底，容易在运行时才暴露问题。
  - 反射开销：启动时慢几毫秒，对 Web 服务几乎无感，但高频读取不合适。
  - 依赖偏重：引入 10+ 间接依赖，不适合极致精简的微服务或边缘组件。

## slog 作为日志管理工具
- log/slog是 Go 标准库原生结构化日志方案，零三方依赖、API 稳定、可与 context/ gin/ 第三方 handler 无缝集成。
- 优点：
    - 标准库（Go 1.21+）无第三方依赖，避免日志组件版本分裂，长期可维护性高。
    - 结构化日志天然支持，内置 InfoContext(ctx, "msg", "user_id", uid)，With(...)可轻松实现链路追踪字段（trace_id/request_id）。
    - Level / Handler 可扩展，支持 slog.LevelDebug、slog.NewJSONHandler，也可自定义 slog.Handler（如输出到 file + rotate / 上报 ELK / Loki）。
    - Context 友好，slog.With("trace_id", traceID)+ slog.LoggerFromContext非常适合 Gin 中间件注入。
    - 生态兼容，主流框架、中间件已开始适配 slog，未来替换或增强成本低。
- 缺点：
    - 功能相对“基础”，不像 zap/zerolog那样内置日志切割（rotate）、异步 buffer、采样等，需要配合 lumberjack或自行封装。
    - 性能略低于 zap（极端场景），在超高频日志场景下吞吐量/分配略逊于 zap，但绝大多数 Web 项目无感知差异。

## gin v1.12.0作为web框架
- Gin 基于 net/http 封装、采用 Radix 树高性能路由的轻量微框架，API 极简、中间件生态成熟，最适合构建RESTful API / 微服务。
- 优点：
    - 性能顶尖，Radix 树路由 + sync.Pool 复用 Context，QPS 高、内存占用低。
    - 开发效率高，内置路由分组、JSON/Form 绑定校验、链式中间件、优雅 panic 恢复。
    - 生态最成熟，JWT/CORS/Swagger/限流等中间件丰富，兼容 http.Handler，社区资料丰富。
    - 零魔法，轻量微内核不捆绑 ORM/Session，技术栈可自主决策。
- 缺点：
    - 非全栈，无内置 ORM、Session、模板引擎、WebSocket，需自行集成第三方库。
    - 不够"约束"，无内置 DI 和强约定，大项目若不分好 internal/层次容易变乱。
    - gin.Context持有 *http.Request，部分场景（如脱离 Gin 调 repository）需注意避免把 Context 穿透过深。

## gin-swagger v1.6.1作为api文档框架
- gin-swagger（swaggo）用 handler 函数上方的注释直接生成 OpenAPI 文档并内嵌 Swagger UI，零业务侵入、与 Gin 集成仅需几行代码，能天然保证文档随接口同步更新。
- 优点：
    - 注释即文档，文档与代码同文件同 PR 维护，不易脱节。
    - 生成 docs.go通过 embed编入二进制，部署无需额外静态资源。
    - Gin 原生适配，ginSwagger.WrapHandler一行注册路由即可访问交互式 UI。
    - 生态成熟、学习成本低，适合快速迭代的中后台 CRUD 项目。
- 缺点：
    - 主要生成 OpenAPI 2.0（Swagger 2.0），对 OpenAPI 3.x 高级特性（oneOf/anyOf、回调等）支持有限。
    - 复杂嵌套参数/泛型结构体注解冗长，注释格式错误只能在 swag init时报错，IDE 智能提示弱。
    - 修改接口后需手动重跑 swag init，遗忘会导致文档过时（可用 go generate+ CI 校验缓解）。
    - 不支持 Go 泛型结构体字段级描述，泛型 Model 需降级为具体类型。

## pgx v5.9.2作为数据库驱动框架
- pgx 直接实现 PostgreSQL 原生二进制协议并完整支持 PG 专有特性（COPY、LISTEN/NOTIFY、jsonb/数组等原生类型、Batch 批量、pgxpool 连接池），相比 database/sql + lib/pq性能更高、功能更贴合 PostgreSQL。
- 优点：
    - 高性能，直连 PG 扩展协议，默认二进制格式传输 + 自动预编译语句缓存，比 database/sql文本协议少一层抽象，QPS 更高、内存分配更少。
    - PG 专有特性一等公民，原生支持 COPY FROM（高速批量导入）、LISTEN/NOTIFY（实时消息）、Batch（单次 RTT 多查询）、pgtype完整映射 uuid/jsonb/数组/inet 等 PG 特有类型。
    - 内置 pgxpool，带最小/最大连接数、空闲超时、健康检查、AfterConnectHook、Stat()监控，比 sql.DB对 PG 更精细。
    - 兼容 database/sql，可通过 stdlib.OpenDB()适配需 *sql.DB的 ORM（如 GORM），也可随时取回原生 pgx.Conn调 PG 专属能力。
    - 可观测性，原生支持 QueryTracer接口，方便接入 OpenTelemetry Tracing。
- 缺点：
    - PG 锁定（Vendor Lock-in），API 和类型是 pgx 专有的，无法无缝切换 MySQL/Oracle 等数据库。
    - 学习曲线略高，v5 原生模式放弃 database/sql兼容接口（pgxpool.Pool不是 *sql.DB），Scan/CollectRows 写法不同于 sql.Rows。
    - 不绑定全功能 ORM，pgx 是驱动层，无自动迁移、关联查询、软删除等 ORM 魔法，复杂 CRUD 需配合 sqlc/手写 SQL，或上层再加 GORM/Ent。

## go-redis v9.20.1作为Redis缓存驱动框架
- Redis 是一款高性能、支持多种数据结构的内存 KV 存储，非常适合作为 Web 应用的分布式缓存、会话存储及临时状态中心。
- 优点：
    - 极高性能，基于内存操作，读写 QPS 可达 10w+，延迟通常 <1ms。
    - 丰富数据结构，String / Hash / List / Set / SortedSet / Bitmap / HyperLogLog / Stream。
    - 高可用方案成熟，支持主从复制、Sentinel、Cluster，社区和运维生态完善。
    - 原子操作 & TTL，原生支持过期、INCR、Lua 脚本，适合计数、限流、锁、Session。
    - 广泛采用，Go + go-redis 生态成熟，与 Gin / 微服务集成成本低。
- 缺点：
    - 数据易失，内存存储，宕机或误操作可能丢数据（可部分通过 AOF/RDB 缓解，但不是强一致持久化）。
    - 内存成本较高，大数据量缓存需合理设计 Key 过期与淘汰策略（LRU/LFU）。
    - 冷热数据边界，不适合作为唯一持久化存储（关系型数据仍应存 PG）。
    - 集群复杂度，Redis Cluster 需关注槽位、重定向及客户端兼容。

## robfig/cron v3.0.1作为本地单机定时任务组件
- robfig/cron/v3 是 Go 生态最成熟、零外部依赖的单机 Cron 表达式调度器，API 简洁且支持秒级精度/时区/Panic 恢复，完美契合单实例后台定时作业（日志清理、缓存刷新、统计汇总等）的需求。
- 优点：
    - 标准 & 扩展 Cron 表达式（@every/@daily、可选秒字段）。
    - 时区感知，WithChain(Recover/DelayIfStillRunning)防 panic 崩溃和任务重入。
    - 并发安全，零外部依赖、轻量、社区验证充分。
- 缺点：
    - 纯内存调度，进程重启后未触发任务丢失，无持久化、无分布式锁（多实例部署会重复执行）。
    - 默认串行调度需注意任务耗时阻塞后续触发。
    - 不支持可视化管理 UI 或任务优先级。

## jwt v5.3.1作为身份认证工具
- JWT 在 Stateless（无会话）的 Web / 微服务架构中，通过自包含的 Token 在客户端与服务端之间安全传递已认证的用户身份与声明，避免服务端集中存储登录态。
- 优点：
    - 无状态（Stateless），服务端不需要 Session / Redis 记录登录状态，易于水平扩展、适合 Gin + 多实例部署。
    - 自包含（Self-contained），Payload 可携带 userID / role / tenant等 Claims，减少每次请求查库。
    - 跨域 / 跨服务友好，适合 SPA + API、微服务、BFF 场景，可与 Gateway / Casbin 配合做鉴权。
    - 标准化 & 生态成熟，Go 的 golang-jwt/jwt/v5稳定，前端（LocalStorage / HttpOnly Cookie）易集成。        
- 缺点：
    - 无法主动撤销（Revoke），Token 过期前默认一直有效，通常需借助：短 AccessToken + 长 RefreshToken。
    - Redis 黑名单 / 用户 version 校验，Payload 可被解码（非加密）。
    - 仅 Base64Url 编码，不能存敏感信息（密码、权限明细应走 Casbin）。
    - 体积可控性要求，Claims 过大时会放大每个请求的 Header 尺寸。
    - 过期时间需权衡，太短影响体验，太长增加安全风险。

## casbin v3.10.0 作为rbac业务处理框架
- Casbin 将 RBAC/ABAC 等访问控制模型抽象为配置文件（PERM 元模型），权限逻辑与业务代码彻底解耦，支持策略内存匹配高性能判定及多存储适配器，是 Go 生态中成熟通用的授权事实标准。
- noho-digital/casbin-pgx-adapter v1.2.1 作为casbin-pgx适配器
- apache/casbin-redis-watcher v2.8.0 作为redis-watcher观察者
- 优点：
    - 支持 ACL / RBAC / RBAC with Domain（多租户）/ ABAC 混合模型，通过修改 .conf即可切换无需改代码。
    - 策略加载内存匹配，微秒级 Enforce()判断，内置 keyMatch/regexMatch适配 RESTful 路径。
    - 丰富 Adapter（PG / MySQL / Redis / etcd）和 Watcher 机制，支持动态增删策略及分布式节点同步。
    - 跨语言实现一致，角色继承、资源角色、超级用户均原生支持。
- 缺点：
    - 只做 Authorization（授权），不负责 Authentication（登录认证），需配合 JWT/Session 使用。
    - PERM 模型 .conf语法和 Matcher 编写有一定学习曲线，字段顺序错会导致静默失败。
    - 大规模策略全量加载有内存占用，分布式需自行配 Watcher 保证策略一致性；复杂 ABAC 自定义函数调试不够直观。
