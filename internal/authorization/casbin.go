package authorization

import (
	"LeotureWeb/internal/config"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	rediswatcher "github.com/casbin/redis-watcher/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxadapter "github.com/noho-digital/casbin-pgx-adapter"
	"github.com/redis/go-redis/v9"
)

// Casbin权限同行状态
const (
	EftAllow = "allow"
	EftDeny  = "deny"
)

// CasbinService 封装 Casbin 权限引擎及其依赖组件
//
// 说明：
// - Enforcer 始终为 *casbin.SyncedCachedEnforcer
// - Redis 仅作为 Watcher（Pub/Sub），不存储 Policy
// - Adapter 使用 ContextFilteredAdapter，支持按条件加载策略
type CasbinService struct {
	Enforcer *casbin.SyncedCachedEnforcer
	Adapter  *pgxadapter.PgxAdapter
}

var casbinService *CasbinService

func SetupCasbin(cfg config.Casbin, db *pgxpool.Pool, redisCli *redis.Client) (*CasbinService, error) {
	casbinService = &CasbinService{}

	// 加载 RBAC 模型
	m, err := model.NewModelFromFile(cfg.ModelPath)
	if err != nil {
		return nil, fmt.Errorf("casbin: model configuration file %s failed to load: %w", cfg.ModelPath, err)
	}

	// 初始化 pgx adapter
	adapter, err := pgxadapter.NewAdapterWithPool(db,
		pgxadapter.WithDatabaseName(db.Config().ConnConfig.Database),
		pgxadapter.WithTableName(cfg.TableName))
	if err != nil {
		return nil, fmt.Errorf("casbin: pgx-adapter failed to load: %w", err)
	}
	casbinService.Adapter = adapter

	// 根据 Redis 是否存在选择初始化方式
	if redisCli != nil && cfg.EnabledRedisWatcher {
		enforcer, err := initEnforcerWithRedisWatcher(m, adapter, cfg.CacheTTL, redisCli, cfg.WatcherChannel)
		if err != nil {
			return nil, fmt.Errorf("casbin: enforcer with redis-watcher failed to load: %w", err)
		}
		casbinService.Enforcer = enforcer
	} else {
		enforcer, err := initEnforcerWithLocalCache(m, adapter, cfg.CacheTTL)
		if err != nil {
			return nil, fmt.Errorf("casbin: enforcer with local-cache failed to load: %w", err)
		}
		casbinService.Enforcer = enforcer
	}

	if err := casbinService.Enforcer.LoadPolicy(); err != nil {
		return nil, err
	}
	slog.Info("casbin: Load all policies, initialization succeeded")
	return casbinService, nil
}

// initEnforcerWithRedisWatcher 初始化带 Redis Watcher 的 Enforcer
//
// 说明：
// - 使用 SyncedCachedEnforcer（本地内存缓存）
// - Redis 仅作为 Pub/Sub 通知通道
// - 禁用 AutoLoadPolicy，完全由 Watcher 驱动
func initEnforcerWithRedisWatcher(m model.Model, adapter *pgxadapter.PgxAdapter, ttl time.Duration,
	redisCli *redis.Client, channel string) (*casbin.SyncedCachedEnforcer, error) {

	enforcer, err := casbin.NewSyncedCachedEnforcer(m, adapter)
	if err != nil {
		return nil, fmt.Errorf("casbin: synced cached enforcer failed to load: %w", err)
	}

	// 禁用自动加载，策略刷新由 Watcher 回调触发
	enforcer.StopAutoLoadPolicy()
	enforcer.EnableAutoNotifyDispatcher(false)
	// 本地缓存 TTL
	enforcer.SetExpireTime(ttl)

	// 创建 Redis Watcher
	watcher, err := rediswatcher.NewWatcher("", rediswatcher.WatcherOptions{
		SubClient: redisCli,
		PubClient: redisCli,
		Channel:   channel,
	})
	if err != nil {
		return nil, fmt.Errorf("casbin: redis-watcher failed to load: %w", err)
	}

	// 注册更新回调
	if err := watcher.SetUpdateCallback(newDefaultUpdateCallback(enforcer)); err != nil {
		return nil, fmt.Errorf("casbin: watcher callback setup failed: %w", err)
	}
	if err := enforcer.SetWatcher(watcher); err != nil {
		return nil, fmt.Errorf("casbin: bind watcher to enforcer failed: %w", err)
	}

	slog.Info("casbin: Redis-watcher enforcer initialized", slog.String("channel", channel), slog.Duration("cache_ttl", ttl))
	return enforcer, nil
}

// initEnforcerWithLocalCache 初始化仅使用本地缓存的 Enforcer
//
// 适用场景：
// - 单实例部署
// - 无需跨进程策略同步
func initEnforcerWithLocalCache(m model.Model, adapter *pgxadapter.PgxAdapter, ttl time.Duration) (*casbin.SyncedCachedEnforcer, error) {
	enforcer, err := casbin.NewSyncedCachedEnforcer(m, adapter)
	if err != nil {
		return nil, fmt.Errorf("casbin: synced cached enforcer failed to load: %w", err)
	}

	// 从数据库定时全量拉取
	enforcer.StartAutoLoadPolicy(ttl)
	// 缓存过期后，下次 Enforce() 时惰性刷新
	enforcer.SetExpireTime(ttl)
	enforcer.EnableAutoNotifyWatcher(false)
	enforcer.EnableAutoNotifyDispatcher(false)
	slog.Info("casbin: Local-cache enforcer initialized", slog.Duration("cache_ttl", ttl))
	return enforcer, nil
}

/*
rediswatcher.Update                        "Update"
rediswatcher.UpdateForAddPolicy            "UpdateForAddPolicy"
rediswatcher.UpdateForRemovePolicy         "UpdateForRemovePolicy"
rediswatcher.UpdateForRemoveFilteredPolicy "UpdateForRemoveFilteredPolicy"
rediswatcher.UpdateForSavePolicy           "UpdateForSavePolicy"
rediswatcher.UpdateForAddPolicies          "UpdateForAddPolicies"
rediswatcher.UpdateForRemovePolicies       "UpdateForRemovePolicies"
rediswatcher.UpdateForUpdatePolicy         "UpdateForUpdatePolicy"
rediswatcher.UpdateForUpdatePolicies       "UpdateForUpdatePolicies"
*/
//Self* 系列方法是 纯内存操作，它们不会触发 Watcher 通知，也不会经过 Adapter 持久化。这在回调中是正确的（避免死循环），但前提是消息来源必须是可信的且已持久化的。
func newDefaultUpdateCallback(e casbin.IEnforcer) func(string) {
	return func(msg string) {
		msgStruct := &rediswatcher.MSG{}

		err := msgStruct.UnmarshalBinary([]byte(msg))
		if err != nil {
			slog.Error("casbin: Redis message parsing failed ", slog.Any("error", err))
			return
		}

		var res bool
		switch msgStruct.Method {
		case rediswatcher.Update, rediswatcher.UpdateForSavePolicy:
			err = e.LoadPolicy()
			res = true
		case rediswatcher.UpdateForAddPolicy:
			res, err = e.SelfAddPolicy(msgStruct.Sec, msgStruct.Ptype, msgStruct.NewRule)
		case rediswatcher.UpdateForAddPolicies:
			res, err = e.SelfAddPolicies(msgStruct.Sec, msgStruct.Ptype, msgStruct.NewRules)
		case rediswatcher.UpdateForRemovePolicy:
			res, err = e.SelfRemovePolicy(msgStruct.Sec, msgStruct.Ptype, msgStruct.NewRule)
		case rediswatcher.UpdateForRemoveFilteredPolicy:
			res, err = e.SelfRemoveFilteredPolicy(msgStruct.Sec, msgStruct.Ptype, msgStruct.FieldIndex, msgStruct.FieldValues...)
		case rediswatcher.UpdateForRemovePolicies:
			res, err = e.SelfRemovePolicies(msgStruct.Sec, msgStruct.Ptype, msgStruct.NewRules)
		case rediswatcher.UpdateForUpdatePolicy:
			res, err = e.SelfUpdatePolicy(msgStruct.Sec, msgStruct.Ptype, msgStruct.OldRule, msgStruct.NewRule)
		case rediswatcher.UpdateForUpdatePolicies:
			res, err = e.SelfUpdatePolicies(msgStruct.Sec, msgStruct.Ptype, msgStruct.OldRules, msgStruct.NewRules)
		default:
			err = errors.New("unknown update type")
		}
		if err != nil {
			slog.Error("casbin: Policies failed to load", slog.Any("error", err))
		}
		if !res {
			slog.Warn("casbin: Policy update returned false by watch",
				slog.Any("method", msgStruct.Method),
				slog.String("sec", msgStruct.Sec),
				slog.String("ptype", msgStruct.Ptype))
		}
	}
}
