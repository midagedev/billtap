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
	"github.com/hckim/billtap/internal/compatibility"
	"github.com/hckim/billtap/internal/config"
	"github.com/hckim/billtap/internal/scenarios"
	"github.com/hckim/billtap/internal/server"
	"github.com/hckim/billtap/internal/storage"
	"github.com/hckim/billtap/internal/webhooks"
)

func main() {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) > 0 {
		switch args[0] {
		case "scenario":
			os.Exit(runScenario(args[1:]))
		case "compatibility":
			os.Exit(runCompatibility(args[1:]))
		}
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

func runCompatibility(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: billtap compatibility <scorecard|inventory> [flags]")
		return scenarios.ExitInvalidConfig
	}

	switch args[0] {
	case "scorecard":
		return runCompatibilityScorecard(args[1:])
	case "inventory":
		return runCompatibilityInventory(args[1:])
	default:
		fmt.Fprintln(os.Stderr, "usage: billtap compatibility <scorecard|inventory> [flags]")
		return scenarios.ExitInvalidConfig
	}
}

func runCompatibilityScorecard(args []string) int {
	outputDir, err := parseCompatibilityScorecardArgs(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, "usage: billtap compatibility scorecard [--output-dir path]")
		return scenarios.ExitInvalidConfig
	}

	scorecard, err := compatibility.WriteArtifacts(context.Background(), compatibility.Options{OutputDir: outputDir})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return scenarios.ExitRuntimeFailure
	}
	fmt.Fprintf(os.Stdout, "compatibility scorecard wrote %s imported=%d skipped=%d unsupported=%d mismatch=%d error=%d\n",
		outputDir,
		scorecard.Summary.Imported,
		scorecard.Summary.Skipped,
		scorecard.Summary.Unsupported,
		scorecard.Summary.Mismatch,
		scorecard.Summary.Error,
	)
	if !scorecard.Summary.Passed {
		return scenarios.ExitAssertionFailed
	}
	return scenarios.ExitPass
}

func runCompatibilityInventory(args []string) int {
	openAPIPath, outputDir, source, err := parseCompatibilityInventoryArgs(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, "usage: billtap compatibility inventory --openapi path [--output-dir path] [--source label]")
		return scenarios.ExitInvalidConfig
	}

	inventory, err := compatibility.WriteInventoryArtifacts(context.Background(), compatibility.InventoryOptions{
		OpenAPIPath: openAPIPath,
		OutputDir:   outputDir,
		Source:      source,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return scenarios.ExitRuntimeFailure
	}
	fmt.Fprintf(os.Stdout, "compatibility inventory wrote %s operations=%d implemented=%d inventory_only=%d schema_validated=%d billtap_only=%d implemented_percent=%.1f schema_validated_percent=%.1f\n",
		outputDir,
		inventory.Summary.TotalOperations,
		inventory.Summary.ImplementedOperations,
		inventory.Summary.InventoryOnlyOperations,
		inventory.Summary.SchemaValidatedOperations,
		inventory.Summary.BilltapOnlyRoutes,
		inventory.Summary.ImplementedPercent,
		inventory.Summary.SchemaValidatedPercent,
	)
	return scenarios.ExitPass
}

func parseCompatibilityScorecardArgs(args []string) (string, error) {
	outputDir := compatibility.DefaultOutputDir
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--output-dir":
			if i+1 >= len(args) {
				return "", fmt.Errorf("%s requires a value", arg)
			}
			i++
			outputDir = args[i]
		case strings.HasPrefix(arg, "--output-dir="):
			outputDir = strings.TrimPrefix(arg, "--output-dir=")
		case strings.HasPrefix(arg, "-"):
			return "", fmt.Errorf("unknown flag %s", arg)
		default:
			return "", fmt.Errorf("unexpected argument %s", arg)
		}
	}
	if strings.TrimSpace(outputDir) == "" {
		return "", fmt.Errorf("output directory is required")
	}
	return outputDir, nil
}

func parseCompatibilityInventoryArgs(args []string) (openAPIPath string, outputDir string, source string, err error) {
	outputDir = compatibility.DefaultOutputDir
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--openapi" || arg == "--output-dir" || arg == "--source":
			if i+1 >= len(args) {
				return "", "", "", fmt.Errorf("%s requires a value", arg)
			}
			i++
			switch arg {
			case "--openapi":
				openAPIPath = args[i]
			case "--output-dir":
				outputDir = args[i]
			case "--source":
				source = args[i]
			}
		case strings.HasPrefix(arg, "--openapi="):
			openAPIPath = strings.TrimPrefix(arg, "--openapi=")
		case strings.HasPrefix(arg, "--output-dir="):
			outputDir = strings.TrimPrefix(arg, "--output-dir=")
		case strings.HasPrefix(arg, "--source="):
			source = strings.TrimPrefix(arg, "--source=")
		case strings.HasPrefix(arg, "-"):
			return "", "", "", fmt.Errorf("unknown flag %s", arg)
		default:
			return "", "", "", fmt.Errorf("unexpected argument %s", arg)
		}
	}
	if strings.TrimSpace(openAPIPath) == "" {
		return "", "", "", fmt.Errorf("OpenAPI path is required")
	}
	if strings.TrimSpace(outputDir) == "" {
		return "", "", "", fmt.Errorf("output directory is required")
	}
	return openAPIPath, outputDir, source, nil
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
