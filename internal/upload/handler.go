package image

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aabbuukkaarr8/internal/model"
	"github.com/aabbuukkaarr8/internal/repository"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/wb-go/wbf/zlog"
)

type service interface {
	MachineForImg(ctx context.Context, img model.Image) (uuid.UUID, error)
}

type Handler struct {
	service service
}

func NewdHandler(s service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Handle(ctx context.Context, msg kafka.Message) error {
	var img model.Image
	if err := json.Unmarshal(msg.Value, &img); err != nil {
		return fmt.Errorf("unmarshal task: %w", err)
	}

	id, err := h.service.MachineForImg(ctx, img)
	if err != nil {
		if errors.Is(err, repository.ErrImageNotFound) {
			return fmt.Errorf("process task: %w", repository.ErrImageNotFound)
		}

		return fmt.Errorf("process task: %w", err)
	}

	zlog.Logger.Printf("image processed: %s", id)

	return nil
}
