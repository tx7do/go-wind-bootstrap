module github.com/tx7do/go-wind-bootstrap/script_engine/source/etcd

go 1.26.3

require (
	github.com/tx7do/go-scripts v0.0.8
	github.com/tx7do/go-scripts/source/etcd v0.0.0-00010101000000-000000000000
	github.com/tx7do/go-wind-bootstrap/conf v0.0.0
	github.com/tx7do/go-wind-bootstrap/script_engine v0.0.0
)

require (
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/coreos/go-systemd/v22 v22.7.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.29.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/cobra v1.10.2 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/tx7do/go-wind v0.0.3 // indirect
	github.com/tx7do/go-wind-bootstrap v0.0.0 // indirect
	github.com/tx7do/go-wind-plugins/registry v0.0.1 // indirect
	go.etcd.io/etcd/api/v3 v3.7.1 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.7.1 // indirect
	go.etcd.io/etcd/client/v3 v3.7.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.28.0 // indirect
	go.yaml.in/yaml/v2 v2.4.4 // indirect
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260727163830-6c54dddc4772 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260727163830-6c54dddc4772 // indirect
	google.golang.org/grpc v1.82.1 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	sigs.k8s.io/yaml v1.6.0 // indirect
)

replace (
	github.com/tx7do/go-scripts/source/etcd => D:\GoProject\go-scripts\source\etcd
	github.com/tx7do/go-wind-bootstrap => ../../..
	github.com/tx7do/go-wind-bootstrap/conf => ../../../conf
	github.com/tx7do/go-wind-bootstrap/script_engine => ../..
)
