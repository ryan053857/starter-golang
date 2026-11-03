// Package app 装配并运行进程。
package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"changeme/internal/config"
	"changeme/internal/httpx"
)

// App 是已装配的运行实例。
type App struct {
	srv *httpx.Server
}

// Open 装配实例；调用方负责 Close。
func Open(_ context.Context, cfg config.Config) (*App, error) {
	return &App{
		srv: httpx.NewServer(cfg, NewHTTPHandler()),
	}, nil
}

// Close 关闭 HTTP，返回第一个错误。
func (a *App) Close() error {
	if a.srv == nil {
		return nil
	}
	err := a.srv.Shutdown()
	a.srv = nil
	return err
}

// Run 运行直至退出；SIGHUP 会热重载。
func Run(ctx context.Context, configFile string) error {
	cfg, err := config.Load(configFile)
	if err != nil {
		return err
	}
	a, err := Open(ctx, cfg)
	if err != nil {
		return err
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	defer signal.Stop(sigCh)

	for {
		runCtx, cancel := context.WithCancel(ctx)
		errCh, wait := a.start(runCtx)
		slog.InfoContext(ctx, "running")

		select {
		case <-ctx.Done():
			cancel()
			wait()
			if closeErr := a.Close(); closeErr != nil {
				return closeErr
			}
			return ctx.Err()
		case runErr := <-errCh:
			cancel()
			wait()
			_ = a.Close()
			return runErr
		case sig := <-sigCh:
			if sig == syscall.SIGHUP {
				slog.InfoContext(ctx, "reloading")
				cancel()
				wait()
				if err := a.Close(); err != nil {
					slog.WarnContext(ctx, "close before reload", "error", err)
				}
				cfg, err = config.Load(configFile)
				if err != nil {
					return err
				}
				a, err = Open(ctx, cfg)
				if err != nil {
					return err
				}
				continue
			}
			slog.InfoContext(ctx, "shutting down")
			cancel()
			wait()
			return a.Close()
		}
	}
}

// start 拉起 HTTP，返回错误通道与等待停妥的函数。
func (a *App) start(ctx context.Context) (<-chan error, func()) {
	errCh := make(chan error, 1)
	var wg sync.WaitGroup
	if a.srv != nil {
		httpErrCh := a.srv.Start()
		wg.Go(func() {
			select {
			case err, ok := <-httpErrCh:
				if ok && err != nil {
					errCh <- err
				}
			case <-ctx.Done():
			}
		})
	}
	return errCh, wg.Wait
}
