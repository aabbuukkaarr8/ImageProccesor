package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aabbuukkaarr8/internal/config"
	"github.com/aabbuukkaarr8/internal/model"
	wbfkafka "github.com/wb-go/wbf/kafka"
	"github.com/wb-go/wbf/retry"
)

// Producer represents a Kafka producer.
type Producer struct {
	Client   *wbfkafka.Producer
	strategy retry.Strategy
	cfg      *config.KafkaConfig
}

func NewProducer(
	cfg *config.KafkaConfig,
	s retry.Strategy,
) *Producer {
	producer := wbfkafka.NewProducer(cfg.Brokers, cfg.Topic)

	return &Producer{
		Client:   producer,
		cfg:      cfg,
		strategy: s,
	}
}

func (p *Producer) Produce(ctx context.Context, img model.Image) error {
	data, err := json.Marshal(img)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %v", err)
	}

	key := []byte(img.ID.String())

	if err = p.Client.SendWithRetry(ctx, p.strategy, key, data); err != nil {
		return fmt.Errorf("failed to send task: %v", err)
	}

	return nil
}
