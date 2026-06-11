module github.com/tx7do/go-wind-bootstrap/transport/http

go 1.26.3

require (
	github.com/tx7do/go-wind v0.0.1
	github.com/tx7do/go-wind-bootstrap v0.0.0
	github.com/tx7do/go-wind-bootstrap/conf v0.0.0
	github.com/tx7do/go-wind-plugins/ratelimit/tokenbucket v0.0.1
	github.com/tx7do/go-wind-plugins/transport/http v0.0.1
	github.com/tx7do/go-wind-plugins/transport/http/middleware/cors v0.0.0-20260611000158-aa445b88bea8
	github.com/tx7do/go-wind-plugins/transport/http/middleware/logging v0.0.0-20260611000158-aa445b88bea8
	github.com/tx7do/go-wind-plugins/transport/http/middleware/ratelimit v0.0.1
	github.com/tx7do/go-wind-plugins/transport/http/middleware/recovery v0.0.0-20260611000158-aa445b88bea8
	github.com/tx7do/go-wind-plugins/transport/http/middleware/requestid v0.0.0-20260611000158-aa445b88bea8
	github.com/tx7do/go-wind-plugins/transport/http/middleware/timeout v0.0.1
	github.com/tx7do/go-wind-plugins/transport/http/middleware/tracing v0.0.1
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/cobra v1.10.2 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	github.com/tx7do/go-wind-plugins/ratelimit v0.0.1 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.43.0 // indirect
	go.opentelemetry.io/otel/metric v1.43.0 // indirect
	go.opentelemetry.io/otel/trace v1.43.0 // indirect
	go.yaml.in/yaml/v2 v2.4.2 // indirect
	golang.org/x/sync v0.20.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	sigs.k8s.io/yaml v1.6.0 // indirect
)

replace (
	github.com/tx7do/go-wind-bootstrap => ../..
	github.com/tx7do/go-wind-bootstrap/conf => ../../conf

	// Local plugins replace (remove before publishing).
	github.com/tx7do/go-wind-plugins/ratelimit => D:\GoProject\go-wind-plugins\ratelimit
	github.com/tx7do/go-wind-plugins/transport => D:\GoProject\go-wind-plugins\transport
	github.com/tx7do/go-wind-plugins/transport/http => D:\GoProject\go-wind-plugins\transport\http
	github.com/tx7do/go-wind-plugins/transport/http/middleware/cors => D:\GoProject\go-wind-plugins\transport\http\middleware\cors
	github.com/tx7do/go-wind-plugins/transport/http/middleware/logging => D:\GoProject\go-wind-plugins\transport\http\middleware\logging
	github.com/tx7do/go-wind-plugins/transport/http/middleware/ratelimit => D:\GoProject\go-wind-plugins\transport\http\middleware\ratelimit
	github.com/tx7do/go-wind-plugins/transport/http/middleware/recovery => D:\GoProject\go-wind-plugins\transport\http\middleware\recovery
	github.com/tx7do/go-wind-plugins/transport/http/middleware/requestid => D:\GoProject\go-wind-plugins\transport\http\middleware\requestid
	github.com/tx7do/go-wind-plugins/transport/http/middleware/timeout => D:\GoProject\go-wind-plugins\transport\http\middleware\timeout
	github.com/tx7do/go-wind-plugins/transport/http/middleware/tracing => D:\GoProject\go-wind-plugins\transport\http\middleware\tracing
)
