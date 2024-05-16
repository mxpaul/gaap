package main

import (
	"log/slog"
	"time"

	"github.com/mxpaul/cancler"
	"github.com/mxpaul/gaap"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

type Config struct {
	SleepInterval time.Duration `yaml:"sleep_interval"`
	SleepLimit    int           `yaml:"sleep_limit"`
}

type App struct {
	Config Config
	Log    *slog.Logger
	Stat   Stat
}

type Stat struct {
	WakeCounter *prometheus.CounterVec
}

func (a *App) Configure(
	cfg Config,
	logger *slog.Logger,
	registry *prometheus.Registry,
) error {
	a.Config = cfg
	a.Log = logger

	a.Stat.WakeCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "wake_up_counter",
			Help: "Number of sleeper wakeups",
		},
		[]string{"sleeper_ident"},
	)

	if err := registry.Register(a.Stat.WakeCounter); err != nil {
		return errors.Wrap(err, "counter register failed")
	}

	return nil
}

func (a *App) Spawn(canc *cancler.Cancler) {
	go func() {
		for i := 1; ; i++ {
			select {
			case <-canc.Done():
				a.Log.Info("finish loop")
				return
			case <-time.After(a.Config.SleepInterval):
				a.Stat.WakeCounter.With(prometheus.Labels{"sleeper_ident": "ident123"}).Inc()
				a.Log.With("count", i).Info("wake up")
				if a.Config.SleepLimit > 0 && i >= a.Config.SleepLimit {
					a.Log.Info("signal daemon to exit")
					canc.Cancel()
				}
			}
		}
	}()
}

func main() {
	gaap.Run[Config](&App{})
}
