module github.com/tx7do/go-wind-bootstrap/transport/socketio

go 1.26.3

require (
	github.com/tx7do/go-wind v0.0.1
	github.com/tx7do/go-wind-bootstrap v0.0.0
	github.com/tx7do/go-wind-bootstrap/conf v0.0.0
	github.com/tx7do/go-wind-plugins/transport/socketio v0.0.1
)

require (
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/gomodule/redigo v1.9.3 // indirect
	github.com/googollee/go-socket.io v1.7.0 // indirect
	github.com/gorilla/handlers v1.5.2 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/cobra v1.10.2 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/tx7do/go-wind-plugins/encoding v0.0.1 // indirect
	go.yaml.in/yaml/v2 v2.4.4 // indirect
	golang.org/x/sync v0.20.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	sigs.k8s.io/yaml v1.6.0 // indirect
)

replace (
	github.com/tx7do/go-wind-bootstrap => ../..
	github.com/tx7do/go-wind-bootstrap/conf => ../../conf

	// Local plugins replace (remove before publishing).
	github.com/tx7do/go-wind-plugins/transport => D:\GoProject\go-wind-plugins\transport
)
