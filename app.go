package gaap

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/mxpaul/cancler"
	"github.com/mxpaul/gaap/pkg/loggy"
	"github.com/mxpaul/gaap/pkg/tooling"
)

type Configurable[ConfigType any] interface {
	Configure(
		cfg ConfigType,
		logger *slog.Logger,
		registry *prometheus.Registry,
	) error
}

type Spawnable interface {
	Spawn(canc *cancler.Cancler)
}

type Application[ConfigType any] struct {
	Opt      CommandLineOptions
	Config   Config[ConfigType]
	App      Spawnable
	Log      *slog.Logger
	Registry *prometheus.Registry
	Tooling  *tooling.ToolingAPI
	canc     *cancler.Cancler
}

func Run[ConfigType any](instance Spawnable) {
	app := Application[ConfigType]{
		App:  instance,
		canc: cancler.NewCancler(context.Background()),
	}

	if err := app.Init(); err != nil {
		if app.Log == nil {
			fmt.Fprintf(os.Stderr, "init error: %v\n", err)
		} else {
			app.Log.Error("init error", "error", err)
		}
		os.Exit(1)
	}

	app.Start()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	select {
	case sig := <-sigc:

		app.canc.Cancel()
		app.Log.With("signal", sig.String()).Warn("exiting on signal")

	case <-app.canc.Done():
		app.Log.Warn("context canceled")
	}

	if app.Config.Daemon.GracefulWait > 0 {
		timeout := app.Config.Daemon.GracefulWait
		app.Log.With("timeout", timeout.String()).Info("graceful shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		go app.Shutdown(ctx)

		select {
		case sig := <-sigc:
			app.Log.With("signal", sig).Error("graceful shutdown interrupted")
			cancel()
		case <-ctx.Done():
		}
	}

	app.Log.Info("exit")
	os.Exit(0)
}

func (app *Application[ConfigType]) Init() error {
	rand.Seed(time.Now().Unix())
	app.Opt = ParseCommandLineOdDie()

	err := LoadConfigFileYAML(app.Opt.ConfigPath, &app.Config)
	if err != nil {
		return errors.Wrap(err, ("config load error"))
	}

	if app.Log, err = loggy.NewLogger(app.Config.Daemon.Log); err != nil {
		return errors.Wrap(err, "logger create failed")
	}
	app.Log.Debug("log create success", "config", app.Config)

	app.Registry = tooling.NewRegistry()
	app.Tooling, err = tooling.NewToolingAPI(app.Config.Daemon.Tooling, app.Log, app.Registry)
	if err != nil {
		return errors.Wrap(err, "tooling create failed")
	}

	if userApp, implements := app.App.(Configurable[ConfigType]); implements {
		err = userApp.Configure(app.Config.Application, app.Log, app.Registry)
		if err != nil {
			return errors.Wrap(err, "user application configure error")
		}
	}

	return nil
}

func (app *Application[ConfigType]) Start() {
	go func() {
		if err := app.Tooling.Start(); err != nil {
			app.Log.With("error", err).Error("tooling api start failed")
		}
		app.canc.Cancel()
	}()
	app.App.Spawn(app.canc)
}

func (app *Application[ConfigType]) Shutdown(ctx context.Context) error {
	return app.Tooling.Shutdown(ctx)
}
