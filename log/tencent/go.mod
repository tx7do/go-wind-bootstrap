module github.com/tx7do/go-wind-bootstrap/log/tencent

go 1.26.3

require (
	github.com/tx7do/go-wind v0.0.1
	github.com/tx7do/go-wind-bootstrap v0.0.0
	github.com/tx7do/go-wind-bootstrap/conf v0.0.0
	github.com/tx7do/go-wind-plugins/log/tencent v0.0.1
)

require (
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/klauspost/compress v1.18.5 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/spf13/cobra v1.10.2 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/tencentcloud/tencentcloud-cls-sdk-go v1.0.14 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.yaml.in/yaml/v2 v2.4.4 // indirect
	golang.org/x/sync v0.20.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	sigs.k8s.io/yaml v1.6.0 // indirect
)

replace (
	github.com/tx7do/go-wind-bootstrap => ../..
	github.com/tx7do/go-wind-bootstrap/conf => ../../conf
)
