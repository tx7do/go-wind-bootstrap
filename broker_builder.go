package bootstrap

import (
	"context"
	"fmt"

	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

// resolveBroker 检查 Broker 配置中每个 optional 字段，
// 对已设置的消息代理类型分别调用对应 builder。
// 返回按类型名索引的 broker 实例映射和统一的 cleanup 函数。
func resolveBroker(ctx context.Context, cfg *v1.Broker) (map[string]any, func(), error) {
	type field struct {
		name    string
		builder BrokerBuilder
	}

	var fields []field

	if cfg.GetKafka() != nil {
		b, err := getBrokerBuilder(BrokerTypeKafka)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: BrokerTypeKafka, builder: b})
	}
	if cfg.GetRabbitmq() != nil {
		b, err := getBrokerBuilder(BrokerTypeRabbitMQ)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: BrokerTypeRabbitMQ, builder: b})
	}
	if cfg.GetRedis() != nil {
		b, err := getBrokerBuilder(BrokerTypeRedis)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: BrokerTypeRedis, builder: b})
	}
	if cfg.GetNats() != nil {
		b, err := getBrokerBuilder(BrokerTypeNATS)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: BrokerTypeNATS, builder: b})
	}
	if cfg.GetMqtt() != nil {
		b, err := getBrokerBuilder(BrokerTypeMQTT)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: BrokerTypeMQTT, builder: b})
	}
	if cfg.GetPulsar() != nil {
		b, err := getBrokerBuilder(BrokerTypePulsar)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: BrokerTypePulsar, builder: b})
	}
	if cfg.GetAzuresb() != nil {
		b, err := getBrokerBuilder(BrokerTypeAzureSB)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: BrokerTypeAzureSB, builder: b})
	}
	if cfg.GetGcpubsub() != nil {
		b, err := getBrokerBuilder(BrokerTypeGCPubSub)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: BrokerTypeGCPubSub, builder: b})
	}
	if cfg.GetNsq() != nil {
		b, err := getBrokerBuilder(BrokerTypeNSQ)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: BrokerTypeNSQ, builder: b})
	}
	if cfg.GetRocketmq() != nil {
		b, err := getBrokerBuilder(BrokerTypeRocketMQ)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: BrokerTypeRocketMQ, builder: b})
	}
	if cfg.GetSqs() != nil {
		b, err := getBrokerBuilder(BrokerTypeSQS)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: BrokerTypeSQS, builder: b})
	}
	if cfg.GetStomp() != nil {
		b, err := getBrokerBuilder(BrokerTypeSTOMP)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: BrokerTypeSTOMP, builder: b})
	}
	if cfg.GetActivemq() != nil {
		b, err := getBrokerBuilder(BrokerTypeActiveMQ)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: BrokerTypeActiveMQ, builder: b})
	}

	if len(fields) == 0 {
		return nil, nil, fmt.Errorf("bootstrap: no broker specified")
	}

	instances := make(map[string]any)
	var cleanups []func()
	for _, f := range fields {
		inst, cleanup, err := f.builder(ctx, cfg)
		if err != nil {
			for _, c := range cleanups {
				c()
			}
			return nil, nil, fmt.Errorf("bootstrap: build broker %q: %w", f.name, err)
		}
		if inst != nil {
			instances[f.name] = inst
		}
		if cleanup != nil {
			cleanups = append(cleanups, cleanup)
		}
	}

	finalCleanup := func() {
		for _, c := range cleanups {
			c()
		}
	}
	return instances, finalCleanup, nil
}
