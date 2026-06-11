module github.com/tx7do/go-wind-bootstrap/_examples/custom_builder

go 1.26.3

require (
	github.com/tx7do/go-wind v0.0.1
	github.com/tx7do/go-wind-bootstrap v0.0.0
	github.com/tx7do/go-wind-bootstrap/conf v0.0.0
	github.com/tx7do/go-wind-bootstrap/log/zap v0.0.0
	github.com/tx7do/go-wind-bootstrap/transport/grpc v0.0.0
	github.com/tx7do/go-wind-bootstrap/transport/http v0.0.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/cobra v1.10.2 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/tx7do/go-wind-plugins/log/zap v0.0.1 // indirect
	github.com/tx7do/go-wind-plugins/ratelimit v0.0.1 // indirect
	github.com/tx7do/go-wind-plugins/ratelimit/tokenbucket v0.0.1 // indirect
	github.com/tx7do/go-wind-plugins/transport/grpc/middleware/logging v0.0.1 // indirect
	github.com/tx7do/go-wind-plugins/transport/grpc/middleware/recovery v0.0.1 // indirect
	github.com/tx7do/go-wind-plugins/transport/grpc/middleware/tracing v0.0.1 // indirect
	github.com/tx7do/go-wind-plugins/transport/grpc/middleware/validate v0.0.1 // indirect
	github.com/tx7do/go-wind-plugins/transport/grpc/server v0.0.0-20260611000158-aa445b88bea8 // indirect
	github.com/tx7do/go-wind-plugins/transport/http v0.0.1 // indirect
	github.com/tx7do/go-wind-plugins/transport/http/middleware/cors v0.0.0-20260611000158-aa445b88bea8 // indirect
	github.com/tx7do/go-wind-plugins/transport/http/middleware/logging v0.0.0-20260611000158-aa445b88bea8 // indirect
	github.com/tx7do/go-wind-plugins/transport/http/middleware/ratelimit v0.0.1 // indirect
	github.com/tx7do/go-wind-plugins/transport/http/middleware/recovery v0.0.0-20260611000158-aa445b88bea8 // indirect
	github.com/tx7do/go-wind-plugins/transport/http/middleware/requestid v0.0.0-20260611000158-aa445b88bea8 // indirect
	github.com/tx7do/go-wind-plugins/transport/http/middleware/timeout v0.0.1 // indirect
	github.com/tx7do/go-wind-plugins/transport/http/middleware/tracing v0.0.1 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.43.0 // indirect
	go.opentelemetry.io/otel/metric v1.43.0 // indirect
	go.opentelemetry.io/otel/trace v1.43.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.28.0 // indirect
	go.yaml.in/yaml/v2 v2.4.4 // indirect
	golang.org/x/net v0.53.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260427160629-7cedc36a6bc4 // indirect
	google.golang.org/grpc v1.80.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	sigs.k8s.io/yaml v1.6.0 // indirect
)

replace (
	github.com/tx7do/go-wind-bootstrap => ../..
	github.com/tx7do/go-wind-bootstrap/conf => ../../conf
	github.com/tx7do/go-wind-bootstrap/log/zap => ../../log/zap
	github.com/tx7do/go-wind-bootstrap/transport/grpc => ../../transport/grpc
	github.com/tx7do/go-wind-bootstrap/transport/http => ../../transport/http

	// Local plugins replace.
	github.com/tx7do/go-wind-plugins/transport/grpc/server => D:\GoProject\go-wind-plugins\transport\grpc\server
	github.com/tx7do/go-wind-plugins/transport/http => D:\GoProject\go-wind-plugins\transport\http
)
