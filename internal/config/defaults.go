package config

import (
	"github.com/spf13/viper"
)

// setDefaults 设置配置项默认参数
func setDefaults(v *viper.Viper) {
	v.SetDefault("app.name", "LeotureWeb")
	v.SetDefault("app.version", "0.0.0")
	v.SetDefault("app.env", "dev")

	v.SetDefault("server.host", "127.0.0.1")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.mode", "release")
	v.SetDefault("server.enabled_cors", false)

	v.SetDefault("cors.allow_origins", []string{"*"})
	v.SetDefault("cors.allow_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	v.SetDefault("cors.allow_headers", []string{
		"Authorization",
		"Content-Type",
		"X-Request-ID",
	})
	v.SetDefault("cors.expose_headers", []string{
		"X-Requested-With",
		"Content-Type",
	})
	v.SetDefault("cors.allow_credentials", true)
	v.SetDefault("cors.max_age", "10m")

	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")
	v.SetDefault("log.add_source", false)
	v.SetDefault("log.file_path", "leoture.log")
	v.SetDefault("log.max_size", 10)
	v.SetDefault("log.max_backups", 10)
	v.SetDefault("log.max_age", 3)
	v.SetDefault("log.compress", false)

	v.SetDefault("scheduler.enabled", false)
	v.SetDefault("scheduler.timezone", "Asia/Shanghai")

	v.SetDefault("database.host", "127.0.0.1")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "leoture")
	v.SetDefault("database.password", "")
	v.SetDefault("database.dbname", "leoture_web")
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.max_conns", 10)
	v.SetDefault("database.min_conns", 5)
	v.SetDefault("database.max_conn_lifetime", "30m")
	v.SetDefault("database.max_conn_idle_time", "5m")
	v.SetDefault("database.health_check_period", "1m")
	v.SetDefault("database.connect_timeout", "3s")
	v.SetDefault("database.statement_timeout", "5s")
	v.SetDefault("database.lock_timeout", "2s")
	v.SetDefault("database.idle_in_transaction_session_timeout", "30s")

	v.SetDefault("redis.enabled", false)
	v.SetDefault("redis.addr", "127.0.0.1:6379")
	v.SetDefault("redis.password", "leoture_cache")
	v.SetDefault("redis.database", 0)
	v.SetDefault("redis.pool_size", 10)
	v.SetDefault("redis.min_idle_conns", 5)
	v.SetDefault("redis.conn_max_lifetime", "1h")
	v.SetDefault("redis.conn_max_idle_time", "5m")
	v.SetDefault("redis.dial_timeout", "5s")
	v.SetDefault("redis.read_timeout", "3s")
	v.SetDefault("redis.write_timeout", "3s")
	v.SetDefault("redis.pool_timeout", "4s")

	v.SetDefault("jwt.issuer", "leoture")
	v.SetDefault("jwt.audience", []string{"leoture-goweb"})
	v.SetDefault("jwt.algorithm", "HS256")
	v.SetDefault("jwt.secret", "leoture-web-default-secret-key-please-change-in-production")
	v.SetDefault("jwt.private_key", "")
	v.SetDefault("jwt.public_key", "")
	v.SetDefault("jwt.access_token_expire", "2h")
	v.SetDefault("jwt.refresh_token_expire", "168h")

	v.SetDefault("casbin.model_path", "config/rbac.conf")
	v.SetDefault("table_name", "casbin_rule")
	v.SetDefault("cache_ttl", "300s")
	v.SetDefault("enabled_redis_watcher", false)
	v.SetDefault("watcher_channel", "casbin_watcher")

	v.SetDefault("swagger.enabled", false)
	v.SetDefault("swagger.title", "LeotureWeb API")
	v.SetDefault("swagger.desc", "基于 Gin 的 Web 项目 API 文档")
	v.SetDefault("swagger.path", "/swagger")
}
