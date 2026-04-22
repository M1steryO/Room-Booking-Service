package app

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/M1steryO/Room-Booking-Service/pkg/logger"
	"net/http"
	"sync"
	"time"

	"github.com/M1steryO/Room-Booking-Service/internal/config"
	httpx "github.com/M1steryO/Room-Booking-Service/internal/delivery/http"
	"github.com/M1steryO/Room-Booking-Service/internal/delivery/http/middleware"
	"github.com/M1steryO/platform_common/pkg/closer"
)

const (
	defaultHTTPReadHeaderTimeout = 3 * time.Second
	defaultHTTPReadTimeout       = 5 * time.Second
	defaultHTTPWriteTimeout      = 10 * time.Second
	defaultHTTPIdleTimeout       = 30 * time.Second
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "local.env", "path to config file")
}

type App struct {
	serviceProvider *serviceProvider
	httpServer      *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	wg := sync.WaitGroup{}
	errCh := make(chan error, 1)

	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := a.runHTTPServer(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
			return err
		}

	case err := <-errCh:
		return err
	}

	wg.Wait()
	return nil
}

func (a *App) Handler() http.Handler {
	return a.httpServer.Handler
}

func (a *App) Seed(ctx context.Context) error {
	room, err := a.serviceProvider.RoomsUsecase(ctx).Create(
		ctx,
		"admin",
		"Focus Room",
		ptr("Small room for quick syncs"),
		ptrInt(4),
	)
	if err != nil {
		return err
	}

	_, err = a.serviceProvider.SchedulesUsecase(ctx).Create(
		ctx,
		"admin",
		room.ID,
		[]int{1, 2, 3, 4, 5},
		"09:00",
		"18:00",
	)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initLogger,
		a.initBootstrap,
		a.initHTTPServer,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	flag.Parse()

	return config.Load(configPath)
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	logger.Init(a.serviceProvider.LoggerConfig().Env())
	return nil
}

// Создание dummy пользователей
func (a *App) initBootstrap(ctx context.Context) error {
	if err := a.serviceProvider.AuthUsecase(ctx).EnsureDummyUsers(ctx); err != nil {
		return fmt.Errorf("failed to ensure dummy users: %w", err)
	}
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	authMW := middleware.AuthMiddleware(a.serviceProvider.JWTManager())

	handler := httpx.NewRouter(
		a.serviceProvider.AuthUsecase(ctx),
		a.serviceProvider.RoomsUsecase(ctx),
		a.serviceProvider.SchedulesUsecase(ctx),
		a.serviceProvider.SlotsUsecase(ctx),
		a.serviceProvider.BookingsUsecase(ctx),
		authMW,
	)

	a.httpServer = &http.Server{
		Addr:              a.serviceProvider.HTTPConfig().Address(),
		Handler:           handler,
		ReadHeaderTimeout: defaultHTTPReadHeaderTimeout,
		ReadTimeout:       defaultHTTPReadTimeout,
		WriteTimeout:      defaultHTTPWriteTimeout,
		IdleTimeout:       defaultHTTPIdleTimeout,
	}

	return nil
}

func (a *App) runHTTPServer() error {
	logger.Info("http server is running on", "addr", a.serviceProvider.HTTPConfig().Address())
	return a.httpServer.ListenAndServe()
}

func ptr(value string) *string { return &value }
func ptrInt(value int) *int    { return &value }
