module github.com/tx7do/go-wind-bootstrap/script_engine/source/redis

go 1.26.3

require (
	github.com/tx7do/go-scripts v0.0.6
	github.com/tx7do/go-scripts/source/redis v0.0.0-00010101000000-000000000000
	github.com/tx7do/go-wind-bootstrap/conf v0.0.0
	github.com/tx7do/go-wind-bootstrap/script_engine v0.0.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/redis/go-redis/v9 v9.18.0 // indirect
	github.com/spf13/cobra v1.10.2 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	github.com/tx7do/go-wind v0.0.1 // indirect
	github.com/tx7do/go-wind-bootstrap v0.0.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.yaml.in/yaml/v2 v2.4.2 // indirect
	golang.org/x/sync v0.20.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	sigs.k8s.io/yaml v1.6.0 // indirect
)

replace (

	github.com/tx7do/go-scripts/source/redis => D:\GoProject\go-scripts\source\redis
	github.com/tx7do/go-wind-bootstrap => ../../..
	github.com/tx7do/go-wind-bootstrap/conf => ../../../conf
	github.com/tx7do/go-wind-bootstrap/script_engine => ../..
)
