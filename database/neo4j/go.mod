module github.com/tx7do/go-wind-bootstrap/database/neo4j

go 1.26.3

require (
	github.com/tx7do/go-crud/neo4j v0.0.1
	github.com/tx7do/go-wind-bootstrap v0.0.0
	github.com/tx7do/go-wind-bootstrap/conf v0.0.0
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jinzhu/copier v0.4.0 // indirect
	github.com/neo4j/neo4j-go-driver/v5 v5.28.5 // indirect
	github.com/spf13/cobra v1.10.2 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/tx7do/go-crud/viewer v0.0.7 // indirect
	github.com/tx7do/go-utils/mapper v0.0.3 // indirect
	github.com/tx7do/go-wind v0.0.3 // indirect
	github.com/tx7do/go-wind-plugins/registry v0.0.1 // indirect
	go.yaml.in/yaml/v2 v2.4.4 // indirect
	golang.org/x/sync v0.22.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	sigs.k8s.io/yaml v1.6.0 // indirect
)

replace (
	github.com/tx7do/go-wind-bootstrap => ../..
	github.com/tx7do/go-wind-bootstrap/conf => ../../conf
)
