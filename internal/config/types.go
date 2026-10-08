package config

import "time"

// 项目环境模式
const (
	EnvDev     = "dev"
	EnvTest    = "test"
	EnvStaging = "staging"
	EnvProd    = "prod"
)

// Config 项目配置根结构
type Config struct {
	App       *App       `mapstructure:"app" yaml:"app"`
	Server    *Server    `mapstructure:"server" yaml:"server"`
	Cors      *Cors      `mapstructure:"cors" yaml:"cors"`
	Log       *Log       `mapstructure:"log" yaml:"log"`
	Scheduler *Scheduler `mapstructure:"scheduler" yaml:"scheduler"`
	Database  *Database  `mapstructure:"database" yaml:"database"`
	Redis     *Redis     `mapstructure:"redis" yaml:"redis"`
	JWT       *JWT       `mapstructure:"jwt" yaml:"jwt"`
	Casbin    *Casbin    `mapstructure:"casbin" yaml:"casbin"`
	Swagger   *Swagger   `mapstructure:"swagger" yaml:"swagger"`
}

// App 项目配置
type App struct {
	// 项目名称
	Name string `mapstructure:"name" yaml:"name"`
	// 项目版本
	Version string `mapstructure:"version" yaml:"version"`
	// 环境模式 dev | test | staging | prod
	Env string `mapstructure:"env" yaml:"env"`
}

// Server HTTP服务配置
type Server struct {
	// 服务监听地址
	Host string `mapstructure:"host" yaml:"host"`
	// 服务监听端口
	Port int `mapstructure:"port" yaml:"port"`
	// Gin 运行模式: debug | release | test
	Mode string `mapstructure:"mode" yaml:"mode"`
	// 是否启用跨域中间件
	EnabledCors bool `mapstructure:"enabled_cors" yaml:"enabled_cors"`
}

// Cors HTTP服务跨域配置
type Cors struct {
	// 允许访问的前端域名, 生产环境务必收紧
	AllowOrigins []string `mapstructure:"allow_origins" yaml:"allow_origins"`
	// 允许的 HTTP 方法
	AllowMethods []string `mapstructure:"allow_methods" yaml:"allow_methods"`
	// 允许的请求头
	AllowHeaders []string `mapstructure:"allow_headers" yaml:"allow_headers"`
	// 暴露给客户端的响应头
	ExposeHeaders []string `mapstructure:"expose_headers" yaml:"expose_headers"`
	// 是否允许携带凭证 (Cookie | Authorization)
	AllowCredentials bool `mapstructure:"allow_credentials" yaml:"allow_credentials"`
	// 预检请求缓存时间
	MaxAge time.Duration `mapstructure:"max_age" yaml:"max_age"`
}

// Log 项目日志配置 (slog)
type Log struct {
	// 日志级别: debug | info | warn | error
	Level string `mapstructure:"level" yaml:"level"`
	// 日志格式: json | text
	Format string `mapstructure:"format" yaml:"format"`
	// 是否记录源码位置（文件名 + 行号）
	AddSource bool `mapstructure:"add_source" yaml:"add_source"`
	// 日志文件路径，空值表示输出到 stdout
	FilePath string `mapstructure:"file_path" yaml:"file_path"`
	// 单个日志文件最大大小（MB）
	MaxSize int `mapstructure:"max_size" yaml:"max_size"`
	// 保留的旧日志文件数量
	MaxBackups int `mapstructure:"max_backups" yaml:"max_backups"`
	// 日志保留天数
	MaxAge int `mapstructure:"max_age" yaml:"max_age"`
	// 是否压缩旧日志文件
	Compress bool `mapstructure:"compress" yaml:"compress"`
}

// Scheduler 定时任务配置 (robfig/cron)
type Scheduler struct {
	// 是否启用定时任务组件
	Enabled bool `mapstructure:"enabled" yaml:"enabled"`
	// 定时任务时区
	Timezone string `mapstructure:"timezone" yaml:"timezone"`
	// 定时任务列表
	Jobs map[string]CronJob `mapstructure:"jobs" yaml:"jobs"`
}

// CronJob 定时任务
type CronJob struct {
	// 任务是否启用
	Enabled bool `mapstructure:"enabled" yaml:"enabled"`
	// Cron 表达式（秒 分 时 日 月 周）
	Spec string `mapstructure:"spec" yaml:"spec"`
	// 任务描述
	Desc string `mapstructure:"desc" yaml:"desc"`
}

// Database 数据库配置 (pgx)
type Database struct {
	// 数据库主机地址
	Host string `mapstructure:"host" yaml:"host"`
	// 数据库端口
	Port int `mapstructure:"port" yaml:"port"`
	// 数据库用户名
	User string `mapstructure:"user" yaml:"user"`
	// 数据库密码
	Password string `mapstructure:"password" yaml:"password"`
	// 数据库名称
	DBName string `mapstructure:"dbname" yaml:"dbname"`
	// SSL 连接模式: disable | require | verify-full
	SSLMode string `mapstructure:"ssl_mode" yaml:"ssl_mode"`
	// 连接池最大连接数（建议 CPU * 2 ~ CPU * 4）
	MaxConns int32 `mapstructure:"max_conns" yaml:"max_conns"`
	// 连接池最小空闲连接数 max_conns/2
	MinConns int32 `mapstructure:"min_conns" yaml:"min_conns"`
	// 连接最大存活时间
	MaxConnLifetime time.Duration `mapstructure:"max_conn_lifetime" yaml:"max_conn_lifetime"`
	// 空闲连接最大存活时间
	MaxConnIdleTime time.Duration `mapstructure:"max_conn_idle_time" yaml:"max_conn_idle_time"`
	// 健康检查周期
	HealthCheckPeriod time.Duration `mapstructure:"health_check_period" yaml:"health_check_period"`
	// 建立连接超时时间
	ConnectTimeout time.Duration `mapstructure:"connect_timeout" yaml:"connect_timeout"`
	// SQL 语句执行超时时间
	StatementTimeout time.Duration `mapstructure:"statement_timeout" yaml:"statement_timeout"`
	// 锁等待超时时间
	LockTimeout time.Duration `mapstructure:"lock_timeout" yaml:"lock_timeout"`
	// 事务内空闲超时时间
	IdleInTransactionSessionTimeout time.Duration `mapstructure:"idle_in_transaction_session_timeout" yaml:"idle_in_transaction_session_timeout"`
}

// Redis 缓存配置 (go-redis)
type Redis struct {
	// 是否启用 Redis
	Enabled bool `mapstructure:"enabled" yaml:"enabled"`
	// 地址（host:port）
	Addr string `mapstructure:"addr" yaml:"addr"`
	// 密码
	Password string `mapstructure:"password" yaml:"password"`
	// 数据库编号
	Database int `mapstructure:"database" yaml:"database"`
	// 连接池最大连接数
	PoolSize int `mapstructure:"pool_size" yaml:"pool_size"`
	// 最小空闲连接数
	MinIdleConns int `mapstructure:"min_idle_conns" yaml:"min_idle_conns"`
	// 空闲连接最大存活时间
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time" yaml:"conn_max_idle_time"`
	// 连接最大存活时间
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" yaml:"conn_max_lifetime"`
	// 建立连接超时
	DialTimeout time.Duration `mapstructure:"dial_timeout" yaml:"dial_timeout"`
	// 读超时
	ReadTimeout time.Duration `mapstructure:"read_timeout" yaml:"read_timeout"`
	// 写超时
	WriteTimeout time.Duration `mapstructure:"write_timeout" yaml:"write_timeout"`
	// 连接池获取连接超时
	PoolTimeout time.Duration `mapstructure:"pool_timeout" yaml:"pool_timeout"`
}

// JWT 权限配置
type JWT struct {
	// Token 签发者标识
	Issuer string `mapstructure:"issuer" yaml:"issuer"`
	// Token 接收者标识
	Audience []string `mapstructure:"audience" yaml:"audience"`
	// 签名算法，仅支持: HS256 | ES256 | EdDSA
	Algorithm string `mapstructure:"algorithm" yaml:"algorithm"`
	// HS256 密钥, 至少32字节
	Secret string `mapstructure:"secret" yaml:"secret"`
	// ES256 | EdDSA 私钥PEM(PKCS#8格式)
	PrivateKey string `mapstructure:"private_key" yaml:"private_key"`
	// ES256 | EdDSA 公钥 PEM(PKIX格式)
	PublicKey string `mapstructure:"public_key" yaml:"public_key"`
	// AccessToken 有效期
	AccessTokenExpire time.Duration `mapstructure:"access_token_expire" yaml:"access_token_expire"`
	// RefreshToken 有效期
	RefreshTokenExpire time.Duration `mapstructure:"refresh_token_expire" yaml:"refresh_token_expire"`
}

// Casbin RBAC 权限配置
type Casbin struct {
	// Casbin 规则表名
	ModelPath string `mapstructure:"model_path" yaml:"model_path"`
	// Casbin 规则表名
	TableName string `mapstructure:"table_name" yaml:"table_name"`
	// 权限策略缓存时间
	CacheTTL time.Duration `mapstructure:"cache_ttl" yaml:"cache_ttl"`
	// 是否启用 Redis Watcher（多实例同步）
	EnabledRedisWatcher bool `mapstructure:"enabled_redis_watcher" yaml:"enable_redis_watcher"`
	// Redis Watcher 订阅频道名称
	WatcherChannel string `mapstructure:"watcher_channel" yaml:"watcher_channel"`
}

// Swagger 接口文档配置
type Swagger struct {
	// 是否启用 Swagger 文档
	Enabled bool `mapstructure:"enabled" yaml:"enabled"`
	// API 文档标题
	Title string `mapstructure:"title" yaml:"title"`
	// API 文档描述
	Desc string `mapstructure:"desc" yaml:"desc"`
	// Swagger UI 访问路径
	Path string `mapstructure:"path" yaml:"path"`
}
