// Package app 按角色装配并运行进程。
package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"

	"github.com/ryan053857/starter-golang/internal/config"
	"github.com/ryan053857/starter-golang/internal/httpx"
)

// Role 是进程启用的能力位。
type Role uint8

const (
	// RoleHTTP 启动 HTTP。
	RoleHTTP Role = 1 << iota
	// RoleWorker 启动消息消费者。
	RoleWorker
	// RoleAll 启用全部角色。
	RoleAll = RoleHTTP | RoleWorker
)

// App 是已装配的运行实例。
type App struct {
	roles  Role
	pubSub *gochannel.GoChannel
	srv    *httpx.Server
}

// Open 按角色装配实例；调用方负责 Close。
func Open(_ context.Context, cfg config.Config, roles Role) (*App, error) {
	if roles == 0 {
		return nil, fmt.Errorf("app: no roles")
	}
	a := &App{roles: roles}
	if roles&(RoleHTTP|RoleWorker) != 0 {
		// gochannel 仅进程内有效：发布与订阅必须共用同一实例。
		a.pubSub = gochannel.NewGoChannel(gochannel.Config{}, watermill.NewSlogLogger(slog.Default()))
	}
	if roles&RoleHTTP != 0 {
		a.srv = httpx.NewServer(cfg, NewHTTPHandler(a.pubSub))
	}
	return a, nil
}

// Close 依次关闭 HTTP 与 pub/sub，返回第一个错误。
func (a *App) Close() error {
	var errs []error
	if a.srv != nil {
		errs = append(errs, a.srv.Shutdown())
		a.srv = nil
	}
	if a.pubSub != nil {
		errs = append(errs, a.pubSub.Close())
		a.pubSub = nil
	}
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

// Run 运行已装配角色直到退出；含 RoleHTTP 时 SIGHUP 会热重载。
func Run(ctx context.Context, configFile string, roles Role) error {
	cfg, err := config.Load(configFile)
	if err != nil {
		return err
	}
	a, err := Open(ctx, cfg, roles)
	if err != nil {
		return err
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	defer signal.Stop(sigCh)

	for {
		runCtx, cancel := context.WithCancel(ctx)
		errCh, wait := a.start(runCtx)
		slog.InfoContext(ctx, "running", "roles", rolesString(roles))

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
			if sig == syscall.SIGHUP && roles&RoleHTTP != 0 {
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
				a, err = Open(ctx, cfg, roles)
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

// start 拉起后台组件，返回错误通道与等待停妥的函数。
func (a *App) start(ctx context.Context) (<-chan error, func()) {
	errCh := make(chan error, 2)
	var wg sync.WaitGroup
	// 先订阅再监听 HTTP，避免 gochannel 在无订阅者时丢弃消息。
	if a.roles&RoleWorker != 0 && a.pubSub != nil {
		messages, err := subscribeDemo(ctx, a.pubSub)
		if err != nil {
			errCh <- err
			return errCh, wg.Wait
		}
		wg.Go(func() {
			if err := consumeDemo(ctx, messages); err != nil && ctx.Err() == nil {
				errCh <- err
			}
		})
	}
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

func rolesString(roles Role) string {
	var parts []string
	if roles&RoleHTTP != 0 {
		parts = append(parts, "http")
	}
	if roles&RoleWorker != 0 {
		parts = append(parts, "worker")
	}
	if len(parts) == 0 {
		return "none"
	}
	return strings.Join(parts, "+")
}
