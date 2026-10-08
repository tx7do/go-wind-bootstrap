module github.com/tx7do/go-wind-bootstrap/database/qdrant

go 1.26.3

require (
	github.com/tx7do/go-crud/qdrant v0.0.1
	github.com/tx7do/go-wind-bootstrap v0.0.0
	github.com/tx7do/go-wind-bootstrap/conf v0.0.0
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jinzhu/copier v0.4.0 // indirect
	github.com/qdrant/go-client v1.19.3 // indirect
	github.com/spf13/cobra v1.10.2 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/tx7do/go-crud/vector v0.0.1 // indirect
	github.com/tx7do/go-crud/viewer v0.0.7 // indirect
	github.com/tx7do/go-utils/mapper v0.0.3 // indirect
	github.com/tx7do/go-wind v0.0.3 // indirect
	github.com/tx7do/go-wind-plugins/registry v0.0.1 // indirect
	go.yaml.in/yaml/v2 v2.4.4 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260706201446-f0a921348800 // indirect
	google.golang.org/grpc v1.84.0 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
	sigs.k8s.io/yaml v1.6.0 // indirect
)

replace (
	github.com/tx7do/go-wind-bootstrap => ../..
	github.com/tx7do/go-wind-bootstrap/conf => ../../conf
)
