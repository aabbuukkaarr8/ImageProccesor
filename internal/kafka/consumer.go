package kafka

import (
	"context"
	"sync"
	"time"

	"github.com/aabbuukkaarr8/internal/config"
	"github.com/segmentio/kafka-go"
	wbfkafka "github.com/wb-go/wbf/kafka"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
)

type uploadedHandler interface {
	Handle(ctx context.Context, msg kafka.Message) error
}

type Consumer struct {
	Client          *wbfkafka.Consumer
	uploadedHandler uploadedHandler
	cfg             *config.KafkaConfig
	strategy        retry.Strategy
}

func NewConsumer(
	cfg *config.KafkaConfig,
	s retry.Strategy,
	uh uploadedHandler,
) *Consumer {
	consumer := wbfkafka.NewConsumer(cfg.Brokers, cfg.Topic, cfg.GroupID)

	return &Consumer{
		Client:          consumer,
		uploadedHandler: uh,
		cfg:             cfg,
		strategy:        s,
	}
}

func (c *Consumer) Consume(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	zlog.Logger.Info().
		Str("topic", c.cfg.Topic).
		Msg("starting consumer")

	for {
		// Exit if context is canceled (graceful shutdown).
		if ctx.Err() != nil {
			zlog.Logger.Info().Msg("shutdown signal received, stopping consumer")
			return
		}

		// Fetch a message from Kafka with retries.
		var msg kafka.Message
		err := retry.Do(func() error {
			var fetchErr error
			msg, fetchErr = c.Client.Fetch(ctx)
			return fetchErr
		}, c.strategy)

		if err != nil {
			// Log error and retry after a short backoff.
			zlog.Logger.Err(err).Msg("failed to fetch message")
			time.Sleep(500 * time.Millisecond)
			continue
		}

		// Process message using the uploadedHandler.
		if err := c.uploadedHandler.Handle(ctx, msg); err != nil {
			zlog.Logger.Err(err).
				Str("message", string(msg.Value)).
				Msg("failed to process image")
			continue
		}

		// Commit the message with retries.
		err = retry.Do(func() error {
			return c.Client.Commit(ctx, msg)
		}, c.strategy)
		if err != nil {
			zlog.Logger.Err(err).Msg("failed to commit message after retries")
		}

		zlog.Logger.Info().
			Int64("offset", msg.Offset).
			Str("message", string(msg.Value)).
			Msg("message handled successfully")
	}
}
