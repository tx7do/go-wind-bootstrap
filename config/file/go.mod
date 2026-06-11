module github.com/tx7do/go-wind-bootstrap/config/file

go 1.26.3

require (
	github.com/tx7do/go-wind-bootstrap v0.0.0
	github.com/tx7do/go-wind-bootstrap/conf v0.0.0
	github.com/tx7do/go-wind-plugins/config/file v0.0.1
)

require (
	github.com/fsnotify/fsnotify v1.10.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/cobra v1.10.2 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	github.com/tx7do/go-wind v0.0.1 // indirect
	github.com/tx7do/go-wind-plugins/config v0.0.1 // indirect
	go.yaml.in/yaml/v2 v2.4.2 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	sigs.k8s.io/yaml v1.6.0 // indirect
)

replace (
	github.com/tx7do/go-wind-bootstrap => ../..
	github.com/tx7do/go-wind-bootstrap/conf => ../../conf

	// Local plugins replace (remove before publishing).
	github.com/tx7do/go-wind-plugins/config => D:\GoProject\go-wind-plugins\config
)
