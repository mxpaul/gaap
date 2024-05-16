package tooling

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	slogecho "github.com/samber/slog-echo"
)

type ToolingAPI struct {
	Echo              *echo.Echo
	HTTPListenAddress string
	Log               *slog.Logger
}

func NewToolingAPI(
	cfg Config,
	logger *slog.Logger,
	registry *prometheus.Registry,
) (*ToolingAPI, error) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	if cfg.LogRequests {
		e.Use(slogecho.New(logger.WithGroup("http")))
	}

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	e.GET(cfg.MetricsPath,
		echoprometheus.NewHandlerWithConfig(
			echoprometheus.HandlerConfig{Gatherer: registry},
		),
	)

	t := &ToolingAPI{
		HTTPListenAddress: cfg.HTTPListenAddress,
		Echo:              e,
		Log:               logger,
	}

	return t, nil
}

func (t *ToolingAPI) Start() (err error) {
	t.Log.With("listen_address", t.HTTPListenAddress).Warn("start tooling http server")
	err = t.Echo.Start(t.HTTPListenAddress)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors.Wrap(err, "echo start failed")
	}
	return nil
}

func (t *ToolingAPI) Shutdown(ctx context.Context) (err error) {
	err = t.Echo.Shutdown(ctx)
	return err
}
