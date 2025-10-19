package service

import (
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"image/color"
)

func Resize(srcPath, dstPath string, width, height int) error {
	img, err := imaging.Open(srcPath)
	if err != nil {
		return err
	}
	resized := imaging.Resize(img, width, height, imaging.Lanczos)
	return imaging.Save(resized, dstPath)
}

func Watermark(srcPath, dstPath, text string) error {
	img, err := imaging.Open(srcPath)
	if err != nil {
		return err
	}
	b := img.Bounds()
	dc := gg.NewContext(b.Dx(), b.Dy())
	dc.DrawImage(img, 0, 0)
	dc.SetColor(color.White)
	dc.LoadFontFace("/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf", 36)
	dc.DrawStringAnchored(text, float64(b.Dx()/2), float64(b.Dy()/2), 0.5, 0.5)
	return dc.SavePNG(dstPath)
}
