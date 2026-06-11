// Package bootstrap provides a declarative application bootstrapper for go-wind.
//
// It reads a [BootstrapConfig] (defined via Protobuf) and assembles the
// corresponding go-wind [App] with the right combination of servers,
// registries, config sources, loggers, tracers, metrics, and brokers
// from go-wind-plugins.
//
// The core flow is:
//
//	BootstrapConfig → Builder registry → wind.Option chain → wind.App
package bootstrap

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"sigs.k8s.io/yaml"

	wind "github.com/tx7do/go-wind"

	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

// Bootstrap reads a [BootstrapConfig], resolves all configured subsystems
// through the builder registry, constructs a [*wind.App], and returns it
// along with a cleanup function.
//
// The caller is responsible for calling [wind.App.Run].
//
// cleanup must be called after the app stops to release resources
// (e.g. deregister from registry, flush tracers).
func Bootstrap(ctx context.Context, cfg *v1.BootstrapConfig) (*wind.App, func(), error) {
	if cfg == nil {
		return nil, nil, fmt.Errorf("bootstrap: config is nil")
	}

	var (
		opts    []wind.Option
		cleanup = func() {}
	)

	// 1. App metadata.
	if appCfg := cfg.GetApp(); appCfg != nil {
		opts = append(opts, resolveApp(appCfg)...)
	}

	// 2. Logger — set early so subsequent builders can log.
	if logCfg := cfg.GetLogger(); logCfg != nil {
		logger, logCleanup, err := resolveLog(logCfg)
		if err != nil {
			return nil, nil, fmt.Errorf("bootstrap: resolve log: %w", err)
		}
		if logger != nil {
			opts = append(opts, wind.WithLogger(logger))
		}
		if logCleanup != nil {
			prev := cleanup
			cleanup = func() { logCleanup(); prev() }
		}
	}

	// 3. Server.
	if srvCfg := cfg.GetServer(); srvCfg != nil {
		servers, srvCleanup, err := resolveServer(srvCfg)
		if err != nil {
			return nil, nil, fmt.Errorf("bootstrap: resolve server: %w", err)
		}
		if len(servers) > 0 {
			opts = append(opts, wind.WithServer(servers...))
		}
		if srvCleanup != nil {
			prev := cleanup
			cleanup = func() { srvCleanup(); prev() }
		}
	}

	// 4. Registry — wire BeforeStop for deregistration.
	if regCfg := cfg.GetRegistry(); regCfg != nil {
		regCleanup, err := resolveRegistry(ctx, regCfg, cfg.GetApp())
		if err != nil {
			return nil, nil, fmt.Errorf("bootstrap: resolve registry: %w", err)
		}
		if regCleanup != nil {
			prev := cleanup
			cleanup = func() { regCleanup(); prev() }
		}
	}

	// 5. Config source.
	if confCfg := cfg.GetConfig(); confCfg != nil {
		confCleanup, err := resolveConfig(ctx, confCfg)
		if err != nil {
			return nil, nil, fmt.Errorf("bootstrap: resolve config source: %w", err)
		}
		if confCleanup != nil {
			prev := cleanup
			cleanup = func() { confCleanup(); prev() }
		}
	}

	// 6. Tracer.
	if tracerCfg := cfg.GetTracer(); tracerCfg != nil {
		tp, tracerCleanup, err := resolveTracer(tracerCfg)
		if err != nil {
			return nil, nil, fmt.Errorf("bootstrap: resolve tracer: %w", err)
		}
		if tp != nil {
			prev := cleanup
			cleanup = func() { tracerCleanup(); prev() }
		}
	}

	// 7. Metrics.
	if metricsCfg := cfg.GetMetrics(); metricsCfg != nil {
		metricsCleanup, err := resolveMetrics(metricsCfg)
		if err != nil {
			return nil, nil, fmt.Errorf("bootstrap: resolve metrics: %w", err)
		}
		if metricsCleanup != nil {
			prev := cleanup
			cleanup = func() { metricsCleanup(); prev() }
		}
	}

	// 8. Broker.
	if brokerCfg := cfg.GetBroker(); brokerCfg != nil {
		brokerCleanup, err := resolveBroker(ctx, brokerCfg)
		if err != nil {
			return nil, nil, fmt.Errorf("bootstrap: resolve broker: %w", err)
		}
		if brokerCleanup != nil {
			prev := cleanup
			cleanup = func() { brokerCleanup(); prev() }
		}
	}

	app := wind.New(opts...)
	return app, cleanup, nil
}

// Run is a convenience function that calls [Bootstrap] and then [wind.App.Run].
// It is intended for simple use cases; for more control, use [Bootstrap] directly.
func Run(ctx context.Context, cfg *v1.BootstrapConfig) error {
	app, cleanup, err := Bootstrap(ctx, cfg)
	if err != nil {
		return err
	}
	defer cleanup()

	return app.Run(ctx)
}

// ---------------------------------------------------------------------------
// High-level (sealed) API
// ---------------------------------------------------------------------------

// BootstrapWithContext is like [Bootstrap] but returns a [*Context] that wraps
// the app, config, cleanup, and a cancellable context.
//
// The caller should call ctx.Cleanup() after the app stops.
func BootstrapWithContext(ctx context.Context, cfg *v1.BootstrapConfig) (*Context, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)

	app, cleanup, err := Bootstrap(ctx, cfg)
	if err != nil {
		cancel()
		return nil, err
	}
	return newContext(cfg, app, cleanup, cancel), nil
}

// RunApp is the sealed one-call entry point. It:
//  1. Creates a signal-aware context (SIGINT/SIGTERM)
//  2. Loads config from the given file path
//  3. Bootstraps the application
//  4. Runs the app and cleans up on exit
//
// This is the simplest way to start an application:
//
//	func main() {
//	    if err := bootstrap.RunApp("config.yaml"); err != nil {
//	        log.Fatal(err)
//	    }
//	}
func RunApp(configPath string) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := LoadConfigFromFile(configPath)
	if err != nil {
		return fmt.Errorf("bootstrap: %w", err)
	}

	app, cleanup, err := Bootstrap(ctx, cfg)
	if err != nil {
		return err
	}
	defer cleanup()

	return app.Run(ctx)
}

// RunAppWithFlags creates a cobra root command that loads config from the
// --conf flag and runs the application. It allows further customisation of
// the root command (adding sub-commands, extra flags, etc.).
//
// Usage:
//
//	func main() {
//	    flags := bootstrap.NewCommandFlags()
//	    if err := bootstrap.RunAppWithFlags(flags); err != nil {
//	        os.Exit(1)
//	    }
//	}
func RunAppWithFlags(flags *CommandFlags, opts ...func(root *cobra.Command)) error {
	if flags == nil {
		flags = NewCommandFlags()
	}

	root := NewRootCmd(flags, func(cmd *cobra.Command, args []string) error {
		return RunApp(flags.Conf)
	})

	for _, opt := range opts {
		if opt != nil {
			opt(root)
		}
	}

	return root.Execute()
}

// LoadConfigFromFile reads a [BootstrapConfig] from the given file path.
// It auto-detects the format by extension:
//   - .yaml, .yml → YAML
//   - .json       → JSON
//   - .bin, .pb   → Protobuf binary
func LoadConfigFromFile(path string) (*v1.BootstrapConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("bootstrap: read config file: %w", err)
	}

	ext := filepath.Ext(path)
	switch ext {
	case ".yaml", ".yml":
		return LoadConfigFromYAML(data)
	case ".bin", ".pb":
		return LoadConfigBinary(data)
	default:
		// .json or unknown → try JSON.
		return LoadConfig(data)
	}
}

// LoadConfigFromYAML unmarshals a YAML-encoded [BootstrapConfig].
func LoadConfigFromYAML(data []byte) (*v1.BootstrapConfig, error) {
	jsonData, err := yaml.YAMLToJSON(data)
	if err != nil {
		return nil, fmt.Errorf("bootstrap: convert YAML to JSON: %w", err)
	}
	return LoadConfig(jsonData)
}

// LoadConfig unmarshals a JSON-encoded [BootstrapConfig].
func LoadConfig(data []byte) (*v1.BootstrapConfig, error) {
	cfg := &v1.BootstrapConfig{}
	if err := protojson.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("bootstrap: unmarshal config: %w", err)
	}
	return cfg, nil
}

// LoadConfigBinary unmarshals a binary-encoded (proto wire format) [BootstrapConfig].
func LoadConfigBinary(data []byte) (*v1.BootstrapConfig, error) {
	cfg := &v1.BootstrapConfig{}
	if err := proto.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("bootstrap: unmarshal binary config: %w", err)
	}
	return cfg, nil
}
