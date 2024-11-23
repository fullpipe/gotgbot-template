package di

import (
	"context"
	"log/slog"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/fx"
)

func PrometheusMetricsServer() fx.Option {
	return fx.Module(
		"MetricsServer",
		fx.Provide(fx.Private, initPrometheusMetricsServer),
		fx.Invoke(func(*http.Server) {}),
	)
}

func initPrometheusMetricsServer(lc fx.Lifecycle) *http.Server {
	r := http.NewServeMux()
	r.Handle("/metrics", promhttp.Handler())

	ms := &http.Server{
		Addr:    ":9090",
		Handler: r,
	}

	lc.Append(fx.StartHook(func(ctx context.Context) error {
		ln, err := net.Listen("tcp", ms.Addr)
		if err != nil {
			return err
		}

		slog.Info("starting prometheus metrics server", slog.Any("addr", ms.Addr))

		go ms.Serve(ln)

		return nil
	}))

	lc.Append(fx.StopHook(func(ctx context.Context) error {
		return ms.Shutdown(ctx)
	}))

	return ms
}
