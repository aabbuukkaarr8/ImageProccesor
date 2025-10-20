package machine

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aabbuukkaarr8/internal/model"
	"github.com/disintegration/imaging"
	"strconv"
)

func (p *Machine) resize(ctx context.Context, img model.Image) (model.Image, error) {
	params := img.Action.Params

	width, err := strconv.Atoi(params["width"])
	if err != nil {
		return model.Image{}, fmt.Errorf("invalid width: %v", err)
	}
	height, err := strconv.Atoi(params["height"])
	if err != nil {
		return model.Image{}, fmt.Errorf("invalid height: %v", err)
	}

	// Load the original image from storage.
	srcReader, err := p.storageMinio.Load(ctx, img.Path)
	if err != nil {
		return model.Image{}, fmt.Errorf("failed to load original image: %w", err)
	}
	defer srcReader.Close()

	// Decode into an image object.
	image, err := imaging.Decode(srcReader)
	if err != nil {
		return model.Image{}, fmt.Errorf("failed to decode image: %w", err)
	}

	// Perform resizing.
	resized := imaging.Resize(image, width, height, imaging.Lanczos)

	// Encode resized image into buffer for storage.
	buf := bytes.NewBuffer(nil)
	if err := imaging.Encode(buf, resized, imaging.JPEG); err != nil {
		return model.Image{}, fmt.Errorf("failed to encode resized image: %w", err)
	}

	// Save resized version.
	dst, err := p.storageMinio.Save(ctx, "resized", img.Filename, buf)
	if err != nil {
		return model.Image{}, fmt.Errorf("failed to save resized image: %w", err)
	}

	img.Path = dst
	img.Status = "processed"

	return img, nil
}
