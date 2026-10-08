package utils

import (
	"fmt"
	"log/slog"
	"net"
	"time"
)

// TCPProbeConfig 定义 TCP 探测配置
type TCPProbeConfig struct {
	Addr     string        // host:port
	Timeout  time.Duration // 单次 Dial 超时
	Retries  int           // 重试次数
	Interval time.Duration // 重试间隔
}

// DefaultTCPProbeConfig 返回默认配置
func DefaultTCPProbeConfig(addr string) TCPProbeConfig {
	return TCPProbeConfig{
		Addr:     addr,
		Timeout:  2 * time.Second,
		Retries:  3,
		Interval: 500 * time.Millisecond,
	}
}

// TCPProbe 执行 TCP 探测（带重试）
func TCPProbe(cfg TCPProbeConfig) error {
	var lastErr error

	for i := 0; i < cfg.Retries; i++ {
		conn, err := net.DialTimeout("tcp", cfg.Addr, cfg.Timeout)
		if err == nil {
			if err := conn.Close(); err != nil {
				return fmt.Errorf(
					"tcp close failed: %w", err)
			}
			return nil
		}

		lastErr = err
		slog.Warn("TCP probe failed",
			slog.String("addr", cfg.Addr),
			slog.Duration("timeout", cfg.Timeout),
			slog.Duration("interval", cfg.Interval),
			"addr", cfg.Addr,
			"attempt", i+1,
			"err", err,
		)

		if i < cfg.Retries-1 {
			time.Sleep(cfg.Interval)
		}
	}

	return fmt.Errorf("tcp probe failed after %d attempts: %w", cfg.Retries, lastErr)
}
