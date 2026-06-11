module github.com/tx7do/go-wind-bootstrap/registry/servicecomb

go 1.26.3

require (
	github.com/go-chassis/sc-client v0.7.0
	github.com/tx7do/go-wind-bootstrap v0.0.0
	github.com/tx7do/go-wind-bootstrap/conf v0.0.0
	github.com/tx7do/go-wind-plugins/registry/servicecomb v0.0.1
)

require (
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/deckarep/golang-set v1.8.0 // indirect
	github.com/go-chassis/cari v0.9.0 // indirect
	github.com/go-chassis/foundation v0.4.0 // indirect
	github.com/go-chassis/openlog v1.1.3 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/karlseguin/ccache/v2 v2.0.8 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/spf13/cobra v1.10.2 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	github.com/tx7do/go-wind v0.0.1 // indirect
	github.com/tx7do/go-wind-plugins/registry v0.0.1 // indirect
	go.yaml.in/yaml/v2 v2.4.2 // indirect
	golang.org/x/sync v0.20.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	sigs.k8s.io/yaml v1.6.0 // indirect
)

replace (
	github.com/tx7do/go-wind-bootstrap => ../..
	github.com/tx7do/go-wind-bootstrap/conf => ../../conf

	// Local plugins replace (remove before publishing).
	github.com/tx7do/go-wind-plugins/registry => D:\GoProject\go-wind-plugins\registry
)
