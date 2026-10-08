package config

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Cfg 全局配置对象（只读访问）
var (
	cfg  = atomic.Value{}
	once sync.Once
)

// SetupViper 初始化Viper配置
func SetupViper() (*Config, error) {
	var setupErr error
	var c *Config
	once.Do(func() {
		v := viper.New()

		// 配置文件类型
		v.SetConfigType("yaml")

		// 设置配置文件路径
		v.SetConfigFile("./config/config.yaml")

		// 校验并设置默认值 (优先级：最低)
		setDefaults(v)

		// 读取配置 (优先级：中)
		if err := v.ReadInConfig(); err != nil {
			setupErr = fmt.Errorf("viper: config read failed: %w", err)
			return
		}

		// 读取环境变量 (优先级：高)
		// 加载 .env（开发环境用，生产环境用系统环境变量）
		err := godotenv.Load(".env")
		if err != nil {
			slog.Warn("godotenv: not found .env file", slog.String("error", err.Error()))
			// 环境变量前缀
			v.SetEnvPrefix("LEOTURE")
		}
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		v.AutomaticEnv()

		// 加载配置数据
		if err := loadConfig(v); err != nil {
			setupErr = fmt.Errorf("viper: config unmarshal or verification failed: %w", err)
			return
		}
		c = cfg.Load().(*Config)

		if c.App.Env == EnvDev {
			// 开启热加载
			v.WatchConfig()
			v.OnConfigChange(func(e fsnotify.Event) {
				if e.Op&fsnotify.Write != 0 {
					slog.Warn("viper: config file changed", slog.String("path", e.Name))
					if err := loadConfig(v); err != nil {
						slog.Error("viper: failed to reload config", slog.String("error", err.Error()))
						return
					}
					slog.Warn("viper: config hot reloading is complete")
				}
			})
		}
	})
	if setupErr != nil {
		return nil, setupErr
	}
	return c, nil
}

// loadConfig 加载配置
func loadConfig(v *viper.Viper) error {
	c := &Config{}
	if err := v.Unmarshal(c); err != nil {
		return err
	}

	if err := verify(*c); err != nil {
		return err
	}

	// 原子替换，保证线上请求安全
	cfg.Store(c)

	return nil
}

// verify 配置校验
func verify(c Config) error {
	// 设置gin模式
	if c.Server.Mode == "" {
		switch c.App.Env {
		case EnvDev:
			c.Server.Mode = gin.DebugMode
		case EnvTest:
			c.Server.Mode = gin.TestMode
		case EnvStaging:
			c.Server.Mode = gin.ReleaseMode
		case EnvProd:
			c.Server.Mode = gin.ReleaseMode
		default:
			c.Server.Mode = gin.DebugMode
		}
	}

	return nil
}
