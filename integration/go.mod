module github.com/tx7do/go-wind-bootstrap/integration

go 1.26.3

require (
	github.com/tx7do/go-wind-bootstrap v0.0.0
	github.com/tx7do/go-wind-bootstrap/conf v0.0.0
	github.com/tx7do/go-wind-bootstrap/registry/etcd v0.0.0
	github.com/tx7do/go-wind-bootstrap/transport/http v0.0.0
	go.etcd.io/etcd/client/v3 v3.6.10
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/coreos/go-systemd/v22 v22.7.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.29.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/cobra v1.10.2 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/tx7do/go-wind v0.0.3 // indirect
	github.com/tx7do/go-wind-plugins/ratelimit v0.0.1 // indirect
	github.com/tx7do/go-wind-plugins/ratelimit/tokenbucket v0.0.1 // indirect
	github.com/tx7do/go-wind-plugins/registry v0.0.1 // indirect
	github.com/tx7do/go-wind-plugins/registry/etcd v0.0.2 // indirect
	github.com/tx7do/go-wind-plugins/transport/http v0.0.1 // indirect
	github.com/tx7do/go-wind-plugins/transport/http/middleware/cors v0.0.0-20260611000158-aa445b88bea8 // indirect
	github.com/tx7do/go-wind-plugins/transport/http/middleware/logging v0.0.0-20260611000158-aa445b88bea8 // indirect
	github.com/tx7do/go-wind-plugins/transport/http/middleware/ratelimit v0.0.1 // indirect
	github.com/tx7do/go-wind-plugins/transport/http/middleware/recovery v0.0.0-20260611000158-aa445b88bea8 // indirect
	github.com/tx7do/go-wind-plugins/transport/http/middleware/requestid v0.0.0-20260611000158-aa445b88bea8 // indirect
	github.com/tx7do/go-wind-plugins/transport/http/middleware/timeout v0.0.1 // indirect
	github.com/tx7do/go-wind-plugins/transport/http/middleware/tracing v0.0.1 // indirect
	go.etcd.io/etcd/api/v3 v3.6.10 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.6.10 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.43.0 // indirect
	go.opentelemetry.io/otel/metric v1.43.0 // indirect
	go.opentelemetry.io/otel/trace v1.43.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.28.0 // indirect
	go.yaml.in/yaml/v2 v2.4.4 // indirect
	golang.org/x/net v0.53.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260427160629-7cedc36a6bc4 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260427160629-7cedc36a6bc4 // indirect
	google.golang.org/grpc v1.80.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	sigs.k8s.io/yaml v1.6.0 // indirect
)

replace github.com/tx7do/go-wind-bootstrap => ../

replace github.com/tx7do/go-wind-bootstrap/conf => ../conf

replace github.com/tx7do/go-wind-bootstrap/registry/etcd => ../registry/etcd

replace github.com/tx7do/go-wind-bootstrap/transport/http => ../transport/http
