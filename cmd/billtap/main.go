package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/hckim/billtap/internal/billing"
	"github.com/hckim/billtap/internal/config"
	"github.com/hckim/billtap/internal/scenarios"
	"github.com/hckim/billtap/internal/server"
	"github.com/hckim/billtap/internal/storage"
	"github.com/hckim/billtap/internal/webhooks"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "scenario" {
		os.Exit(runScenario(os.Args[2:]))
	}

	configPath := flag.String("config", "", "optional path to a JSON config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		slog.Error("load config", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	store, err := storage.Open(ctx, storage.Options{Driver: storage.DriverSQLite, DSN: cfg.DatabaseURL})
	if err != nil {
		slog.Error("open storage", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := store.Close(); err != nil {
			slog.Warn("close storage", "error", err)
		}
	}()

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           server.New(server.Options{Config: cfg, Store: store}),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.Warn("shutdown server", "error", err)
		}
	}()

	slog.Info("starting billtap", "addr", cfg.Addr, "environment", cfg.Environment)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("serve", "error", err)
		os.Exit(1)
	}
}

func runScenario(args []string) int {
	if len(args) == 0 || args[0] != "run" {
		fmt.Fprintln(os.Stderr, "usage: billtap scenario run <file> [--report-json path] [--report-md path] [--database-url dsn]")
		return scenarios.ExitInvalidConfig
	}

	scenarioPath, reportJSON, reportMD, databaseURL, err := parseScenarioRunArgs(args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, "usage: billtap scenario run <file> [--report-json path] [--report-md path] [--database-url dsn]")
		return scenarios.ExitInvalidConfig
	}

	scenario, err := scenarios.LoadFile(scenarioPath)
	if err != nil {
		report := scenarios.Report{
			Name:        scenarioPath,
			Status:      "failed",
			FailureType: scenarios.FailureInvalidConfig,
			StartedAt:   time.Now().UTC(),
			FinishedAt:  time.Now().UTC(),
			Errors:      []string{err.Error()},
		}
		writeScenarioReports(report, reportJSON, reportMD)
		fmt.Fprintln(os.Stderr, err)
		return report.ExitCode()
	}

	ctx := context.Background()
	store, err := storage.OpenSQLite(ctx, databaseURL)
	if err != nil {
		report := runtimeFailureReport(scenario.Name, err)
		writeScenarioReports(report, reportJSON, reportMD)
		fmt.Fprintln(os.Stderr, err)
		return report.ExitCode()
	}
	defer func() {
		if err := store.Close(); err != nil {
			slog.Warn("close scenario store", "error", err)
		}
	}()

	runner := scenarios.NewRunner(billing.NewService(store), webhooks.NewService(store))
	report, runErr := runner.Run(ctx, scenario)
	writeScenarioReports(report, reportJSON, reportMD)
	if runErr != nil {
		fmt.Fprintln(os.Stderr, runErr)
	}
	fmt.Fprintf(os.Stdout, "scenario %s %s\n", report.Name, report.Status)
	return report.ExitCode()
}

func parseScenarioRunArgs(args []string) (scenarioPath string, reportJSON string, reportMD string, databaseURL string, err error) {
	databaseURL = ":memory:"
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--report-json" || arg == "--report-md" || arg == "--database-url":
			if i+1 >= len(args) {
				return "", "", "", "", fmt.Errorf("%s requires a value", arg)
			}
			i++
			switch arg {
			case "--report-json":
				reportJSON = args[i]
			case "--report-md":
				reportMD = args[i]
			case "--database-url":
				databaseURL = args[i]
			}
		case strings.HasPrefix(arg, "--report-json="):
			reportJSON = strings.TrimPrefix(arg, "--report-json=")
		case strings.HasPrefix(arg, "--report-md="):
			reportMD = strings.TrimPrefix(arg, "--report-md=")
		case strings.HasPrefix(arg, "--database-url="):
			databaseURL = strings.TrimPrefix(arg, "--database-url=")
		case strings.HasPrefix(arg, "-"):
			return "", "", "", "", fmt.Errorf("unknown flag %s", arg)
		default:
			if scenarioPath != "" {
				return "", "", "", "", fmt.Errorf("multiple scenario files provided")
			}
			scenarioPath = arg
		}
	}
	if scenarioPath == "" {
		return "", "", "", "", fmt.Errorf("scenario file is required")
	}
	return scenarioPath, reportJSON, reportMD, databaseURL, nil
}

func writeScenarioReports(report scenarios.Report, jsonPath string, markdownPath string) {
	if jsonPath != "" {
		body, err := report.JSON()
		if err != nil {
			fmt.Fprintf(os.Stderr, "render JSON report: %v\n", err)
		} else if err := os.WriteFile(jsonPath, body, 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "write JSON report: %v\n", err)
		}
	}
	if markdownPath != "" {
		if err := os.WriteFile(markdownPath, []byte(report.Markdown()), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "write Markdown report: %v\n", err)
		}
	}
}

func runtimeFailureReport(name string, err error) scenarios.Report {
	now := time.Now().UTC()
	return scenarios.Report{
		Name:        name,
		Status:      "failed",
		FailureType: scenarios.FailureRunner,
		StartedAt:   now,
		FinishedAt:  now,
		Errors:      []string{err.Error()},
	}
}
