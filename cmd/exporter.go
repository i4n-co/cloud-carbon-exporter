package main

import (
	"cloudcarbonexporter"
	"cloudcarbonexporter/internal/demo"
	"cloudcarbonexporter/internal/gcp"
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

func main() {
	ctx := context.Background()

	cloudProvider := ""
	projectID := ""
	listen := ""
	logLevel := ""
	logFormat := ""
	demoEnabled := false

	flag.StringVar(&cloudProvider, "cloud.provider", "", "cloud provider type (gcp, aws, azure)")
	flag.StringVar(&projectID, "gcp.projectid", "", "gcp project to export data from")
	flag.StringVar(&listen, "listen", "0.0.0.0:2922", "addr to listen to")
	flag.StringVar(&logLevel, "log.level", "info", "log severity (debug, info, warn, error)")
	flag.StringVar(&logFormat, "log.format", "text", "log format (text, json)")
	flag.BoolVar(&demoEnabled, "demo.enabled", false, "return fictive demo data")

	flag.Parse()

	switch logFormat {
	case "text":
		slog.SetDefault(slog.New(tint.NewHandler(os.Stdout, &tint.Options{
			Level:   slogLevel(logLevel),
			NoColor: !isatty.IsTerminal(os.Stdout.Fd()),
		})))
	case "json":
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slogLevel(logLevel),
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				switch a.Key {
				case slog.LevelKey:
					a.Key = "severity"
					return a
				case slog.MessageKey:
					a.Key = "message"
					return a
				default:
					return a
				}
			},
		})))
	}

	var collector cloudcarbonexporter.Collector
	var err error
	switch cloudProvider {
	case "gcp":
		if projectID == "" {
			slog.Error("project id is not set")
			flag.PrintDefaults()
			os.Exit(1)
		}
		collector, err = gcp.NewCollector(ctx, projectID)
		if err != nil {
			slog.Error("failed to create gcp collector", "project_id", projectID, "err", err)
			os.Exit(1)
		}
	case "":
		slog.Error("cloud provider is not set")
		flag.PrintDefaults()
		os.Exit(1)
	default:
		slog.Error("cloud provider is not supported yet", "cloud.provider", cloudProvider)
		os.Exit(1)
	}

	collectors := []cloudcarbonexporter.Collector{collector}
	if demoEnabled {
		collectors = append(collectors, demo.NewCollector())
	}
	http.Handle("/metrics", cloudcarbonexporter.NewHTTPMetricsHandler(collectors...))

	slog.Info("starting cloud carbon exporter", "listen", listen)
	if err := http.ListenAndServe(listen, nil); err != nil {
		slog.Error("failed to start cloud carbon exporter", "err", err)
		os.Exit(1)
	}
}

func slogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	}

	return slog.LevelInfo
}
