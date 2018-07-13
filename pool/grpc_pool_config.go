package pool

import (
	"time"
)

// GRPCOptionFunc GRPCOptionFunc
type GRPCOptionFunc func(pool *GRPCConfig)

type GRPCConfig struct {
	*ClientPool
	connectTimeout time.Duration
	readTimeout    time.Duration
}

// GenGRPCClientPool grpc链接池
func GenGRPCClientPool(options ...GRPCOptionFunc) *GRPCConfig {
	p := &GRPCConfig{
		ClientPool: &ClientPool{Wait: true},
	}

	for i := range options {
		options[i](p)
	}

	if p.MaxIdle == 0 {
		p.MaxIdle = 256
	}
	if p.connectTimeout == 0 {
		p.connectTimeout = time.Millisecond * 100
	}
	if p.readTimeout == 0 {
		p.readTimeout = time.Millisecond * 100
	}

	if p.IdleTimeout == 0 {
		p.IdleTimeout = time.Minute * 30
	}

	return p
}

// WithGRPCMaxIdle max idl
func WithGRPCMaxIdle(n uint32) GRPCOptionFunc {
	return func(p *GRPCConfig) {
		p.MaxIdle = int(n)
	}
}

// WithGRPCConnectTimeout conn timeout
func WithGRPCConnectTimeout(t time.Duration) GRPCOptionFunc {
	return func(p *GRPCConfig) {
		p.connectTimeout = t
	}
}

// WithGRPCReadTimeout read timeout
func WithGRPCReadTimeout(t time.Duration) GRPCOptionFunc {
	return func(p *GRPCConfig) {
		p.readTimeout = t
	}
}

// WithGRPCIdleTimeout idl timeout
func WithGRPCIdleTimeout(t time.Duration) GRPCOptionFunc {
	return func(p *GRPCConfig) {
		p.IdleTimeout = t
	}
}

// WithGRPCMaxActive set max active conn in pool
func WithGRPCMaxActive(n int) GRPCOptionFunc {
	return func(p *GRPCConfig) {
		p.MaxActive = n
	}
}

// WithGRPCMaxLiveTime is conn max lifetime
func WithGRPCMaxLiveTime(t time.Duration) GRPCOptionFunc {
	return func(p *GRPCConfig) {
		p.MaxLiveTime = t
	}
}
