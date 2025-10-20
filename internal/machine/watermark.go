package machine

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aabbuukkaarr8/internal/model"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"image/color"
)

const defaultFontPath = "internal/assets/fonts/DejaVuSans.ttf"

func (p *Machine) watermark(ctx context.Context, img model.Image) (model.Image, error) {
	params := img.Action.Params

	text := params["text"]
	if text == "" {
		text = "Watermark"
	}

	// Load the original image.
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

	// Draw watermark text on top of the image.
	dc := gg.NewContextForImage(image)
	dc.SetColor(color.White)

	fontSize := float64(dc.Width()) * 0.05 // 5% of the image width

	err = dc.LoadFontFace(defaultFontPath, fontSize)
	if err != nil {
		return model.Image{}, fmt.Errorf("failed to load font: %w", err)
	}

	tw, th := dc.MeasureString(text) // calculate font size

	margin := 10.0
	x := float64(dc.Width()) - tw - margin
	y := float64(dc.Height()) - th - margin

	dc.DrawStringAnchored(text, x, y, 1, 1) // bottom-right corner
	dc.Fill()

	// Encode modified image.
	buf := new(bytes.Buffer)
	if err := imaging.Encode(buf, dc.Image(), imaging.JPEG); err != nil {
		return model.Image{}, fmt.Errorf("failed to encode watermarked image: %w", err)
	}

	// Save watermarked version.
	dst, err := p.storageMinio.Save(ctx, "watermarked", img.Filename, buf)
	if err != nil {
		return model.Image{}, fmt.Errorf("failed to save watermarked image: %w", err)
	}

	img.Path = dst
	img.Status = "processed"

	return img, nil
}
