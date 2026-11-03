package httpx

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"changeme/internal/config"
)

const shutdownTimeout = 30 * time.Second

// Server 是可启动与关闭的 HTTP 服务。
type Server struct {
	httpServer *http.Server
	env        string
}

// NewServer 按配置与根处理器构造服务。
func NewServer(cfg config.Config, handler http.Handler) *Server {
	return &Server{
		env: cfg.Env,
		httpServer: &http.Server{
			Addr:              cfg.Addr,
			Handler:           handler,
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       15 * time.Minute,
			WriteTimeout:      15 * time.Minute,
			IdleTimeout:       60 * time.Second,
		},
	}
}

// Start 在后台监听，并通过通道返回非预期退出错误。
func (s *Server) Start() <-chan error {
	errCh := make(chan error, 1)
	go func() {
		slog.InfoContext(context.Background(), "http listening", "addr", s.httpServer.Addr, "env", s.env)
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()
	return errCh
}

// Shutdown 在 shutdownTimeout 内优雅关闭 HTTP；超时则记录警告并强制关闭剩余连接，
// 仅返回强制关闭本身的错误。
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	err := s.httpServer.Shutdown(ctx)
	if err == nil {
		return nil
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		return err
	}
	slog.WarnContext(ctx, "http shutdown timed out, closing remaining connections", "timeout", shutdownTimeout)
	return s.httpServer.Close()
}
