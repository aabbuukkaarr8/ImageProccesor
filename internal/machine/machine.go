package machine

import (
	"context"
	"fmt"
	"github.com/aabbuukkaarr8/internal/model"
	"io"
)

type StorageMinio interface {
	Save(ctx context.Context, subdir, filename string, src io.Reader) (string, error)
	Load(ctx context.Context, path string) (io.ReadCloser, error)
}

type Machine struct {
	storageMinio StorageMinio
}

func New(sm StorageMinio) *Machine {
	return &Machine{storageMinio: sm}
}

func (m *Machine) Machine(ctx context.Context, img model.Image) (model.Image, error) {
	switch img.Action.Name {
	case "resize":
		return m.resize(ctx, img)
	case "thumbnail":
		return m.thumbnail(ctx, img)
	case "watermark":
		return m.watermark(ctx, img)
	default:
		return model.Image{}, fmt.Errorf("unknown task action: %s", img.Action.Name)
	}
}
