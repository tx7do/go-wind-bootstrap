package bootstrap

// Server type 常量，用于注册 Builder 时的 key。
const (
	ServerTypeHTTP         = "http"
	ServerTypeHTTP3        = "http3"
	ServerTypeGRPC         = "grpc"
	ServerTypeGraphQL      = "graphql"
	ServerTypeSSE          = "sse"
	ServerTypeWebSocket    = "websocket"
	ServerTypeTCP          = "tcp"
	ServerTypeUDP          = "udp"
	ServerTypeKCP          = "kcp"
	ServerTypeThrift       = "thrift"
	ServerTypeTRPC         = "trpc"
	ServerTypeWebTransport = "webtransport"
	ServerTypeCron         = "cron"
	ServerTypeHPTimer      = "hptimer"
	ServerTypeMCP          = "mcp"
	ServerTypeSignalR      = "signalr"
	ServerTypeSocketIO     = "socketio"
	ServerTypeWebRTC       = "webrtc"
	ServerTypeAsynq        = "asynq"
	ServerTypeMachinery    = "machinery"

	// Broker-based transport servers.
	ServerTypeKafka    = "kafka"
	ServerTypeRabbitMQ = "rabbitmq"
	ServerTypeRedis    = "redis"
	ServerTypeNATS     = "nats"
	ServerTypeMQTT     = "mqtt"
	ServerTypePulsar   = "pulsar"
	ServerTypeActiveMQ = "activemq"
	ServerTypeAzureSB  = "azuresb"
	ServerTypeNSQ      = "nsq"
	ServerTypeRocketMQ = "rocketmq"
	ServerTypeSQS      = "sqs"
)

// Config type 常量，用于注册 ConfigAction 时的 key。
const (
	ConfigTypeFile       = "file"
	ConfigTypeFs         = "fs"
	ConfigTypeEtcd       = "etcd"
	ConfigTypeNacos      = "nacos"
	ConfigTypeConsul     = "consul"
	ConfigTypeApollo     = "apollo"
	ConfigTypeKubernetes = "kubernetes"
	ConfigTypeRedis      = "redis"
	ConfigTypeZookeeper  = "zookeeper"
	ConfigTypeVault      = "vault"
	ConfigTypeHTTP       = "http"
	ConfigTypeEnv        = "env"
	ConfigTypeOSS        = "oss"
	ConfigTypePolaris    = "polaris"
)

// Registry type 常量。
const (
	RegistryTypeConsul      = "consul"
	RegistryTypeEtcd        = "etcd"
	RegistryTypeNacos       = "nacos"
	RegistryTypeZookeeper   = "zookeeper"
	RegistryTypePolaris     = "polaris"
	RegistryTypeEureka      = "eureka"
	RegistryTypeKubernetes  = "kubernetes"
	RegistryTypeServiceComb = "service_comb"
)

// Logger type 常量。
const (
	LoggerTypeZap        = "zap"
	LoggerTypeZerolog    = "zerolog"
	LoggerTypeSlog       = "slog"
	LoggerTypeLogrus     = "logrus"
	LoggerTypeCharm      = "charm"
	LoggerTypePhuslu     = "phuslu"
	LoggerTypeGlog       = "glog"
	LoggerTypeHclog      = "hclog"
	LoggerTypeFluent     = "fluent"
	LoggerTypeLoki       = "loki"
	LoggerTypeSentry     = "sentry"
	LoggerTypeAliyun     = "aliyun"
	LoggerTypeTencent    = "tencent"
	LoggerTypeCloudWatch = "cloudwatch"
)

// Tracer type 常量。
const (
	TracerTypeOTLP = "otlp"
)

// Metrics type 常量。
const (
	MetricsTypePrometheus = "prometheus"
	MetricsTypeOTLP       = "otlp"
)

// Broker type 常量，用于注册 BrokerBuilder 时的 key。
const (
	BrokerTypeKafka    = "kafka"
	BrokerTypeRabbitMQ = "rabbitmq"
	BrokerTypeRedis    = "redis"
	BrokerTypeNATS     = "nats"
	BrokerTypeMQTT     = "mqtt"
	BrokerTypePulsar   = "pulsar"
	BrokerTypeAzureSB  = "azuresb"
	BrokerTypeGCPubSub = "gcpubsub"
	BrokerTypeNSQ      = "nsq"
	BrokerTypeRocketMQ = "rocketmq"
	BrokerTypeSQS      = "sqs"
	BrokerTypeSTOMP    = "stomp"
	BrokerTypeActiveMQ = "activemq"
)

// Storage type 常量，用于注册 StorageBuilder 时的 key。
const (
	StorageTypeMinio = "minio"
	StorageTypeS3    = "s3"
)

// AI type 常量，用于注册 AiBuilder 时的 key。
const (
	AiTypeOpenAI      = "openai"
	AiTypeLangChainGo = "langchaingo"
	AiTypeEino        = "eino"
)

// Workflow type 常量，用于注册 WorkflowBuilder 时的 key。
const (
	WorkflowTypeTemporal    = "temporal"
	WorkflowTypeArgo        = "argo"
	WorkflowTypeConductor   = "conductor"
	WorkflowTypeGoWorkflows = "goworkflows"
)

// Cache type 常量，用于注册 CacheBuilder 时的 key。
const (
	CacheTypeLocal = "local"
	CacheTypeRedis = "redis"
)

// Script Engine type 常量。
const (
	ScriptEngineLua        = "lua"
	ScriptEngineJavaScript = "javascript"
	ScriptEngineGPython    = "gpython"
	ScriptEngineYaegi      = "yaegi"
	ScriptEngineWazero     = "wazero"
	ScriptEngineCEL        = "cel"
	ScriptEngineExpr       = "expr"
	ScriptEngineStarlark   = "starlark"
	ScriptEngineTcl        = "tcl"
)

// Database type 常量，用于注册 DatabaseBuilder 时的 key。
const (
	DatabaseTypeGorm          = "gorm"
	DatabaseTypeMongodb       = "mongodb"
	DatabaseTypeClickhouse    = "clickhouse"
	DatabaseTypeDoris         = "doris"
	DatabaseTypeElasticsearch = "elasticsearch"
	DatabaseTypeOpensearch    = "opensearch"
	DatabaseTypeInfluxdb      = "influxdb"
	DatabaseTypeCassandra     = "cassandra"
)
