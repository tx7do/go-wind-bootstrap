module github.com/tx7do/go-wind-bootstrap

go 1.26.3

require (
	github.com/spf13/cobra v1.10.2
	github.com/tx7do/go-wind v0.0.1
	github.com/tx7do/go-wind-bootstrap/conf v0.0.0
	google.golang.org/protobuf v1.36.11
	sigs.k8s.io/yaml v1.6.0
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	go.yaml.in/yaml/v2 v2.4.2 // indirect
	golang.org/x/sync v0.20.0 // indirect
)

replace github.com/tx7do/go-wind-bootstrap/conf => ./conf
