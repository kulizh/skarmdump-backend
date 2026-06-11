package utils

import (
	"bytes"
	"image"
	"image/png"

	"golang.org/x/image/draw"
)

func ResizeToPNG(data []byte, maxWidth int) ([]byte, error) {
	return ResizeToPNGBounds(data, maxWidth, 0, false)
}

func ResizeToPNGBounds(data []byte, maxWidth, maxHeight int, sharpen bool) ([]byte, error) {
	src, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	bounds := src.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	newW, newH := fitBounds(w, h, maxWidth, maxHeight)

	var dst *image.RGBA
	if newW != w || newH != h {
		dst = image.NewRGBA(image.Rect(0, 0, newW, newH))
		draw.BiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	} else {
		dst = image.NewRGBA(image.Rect(0, 0, w, h))
		draw.Draw(dst, dst.Bounds(), src, src.Bounds().Min, draw.Over)
	}

	if sharpen {
		applySharpen(dst)
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, dst)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func fitBounds(w, h, maxW, maxH int) (int, int) {
	if maxH == 0 {
		if w <= maxW {
			return w, h
		}
		return maxW, h * maxW / w
	}

	if w <= maxW && h <= maxH {
		return w, h
	}

	ratioW := float64(maxW) / float64(w)
	ratioH := float64(maxH) / float64(h)
	ratio := min(ratioW, ratioH)

	newW := int(float64(w) * ratio)
	newH := int(float64(h) * ratio)
	if newW < 1 {
		newW = 1
	}
	if newH < 1 {
		newH = 1
	}
	return newW, newH
}

func applySharpen(img *image.RGBA) {
	kernel := [9]float64{
		0, -1, 0,
		-1, 5, -1,
		0, -1, 0,
	}

	bounds := img.Bounds()
	srcData := make([]byte, len(img.Pix))
	copy(srcData, img.Pix)

	for y := bounds.Min.Y + 1; y < bounds.Max.Y-1; y++ {
		for x := bounds.Min.X + 1; x < bounds.Max.X-1; x++ {
			var r, g, b, a float64
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					k := kernel[(ky+1)*3+(kx+1)]
					stride := img.Stride
					idx := (y+ky)*stride + (x+kx)*4
					r += float64(srcData[idx+0]) * k
					g += float64(srcData[idx+1]) * k
					b += float64(srcData[idx+2]) * k
					a += float64(srcData[idx+3]) * k
				}
			}

			idx := y*img.Stride + x*4
			img.Pix[idx+0] = clamp(r)
			img.Pix[idx+1] = clamp(g)
			img.Pix[idx+2] = clamp(b)
			img.Pix[idx+3] = clamp(a)
		}
	}
}

func clamp(v float64) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v)
}

const ThumbnailWidth = 1200
