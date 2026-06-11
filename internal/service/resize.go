package service

import (
	"skarmdump-backend/internal/config"
	"skarmdump-backend/pkg/utils"
)

type Resizer struct {
	enabled bool
	sharpen bool
}

func NewResizer(cfg *config.Config) *Resizer {
	return &Resizer{
		enabled: cfg.ResizeImages,
		sharpen: cfg.Sharpen,
	}
}

type VariantSet struct {
	OG []byte
	RE []byte
}

func (r *Resizer) Generate(origData []byte) (*VariantSet, error) {
	if !r.enabled {
		return &VariantSet{}, nil
	}

	og, err := utils.ResizeToPNGBounds(origData, 1200, 630, r.sharpen)
	if err != nil {
		return nil, err
	}

	re, err := utils.ResizeToPNGBounds(origData, 1920, 1080, r.sharpen)
	if err != nil {
		return nil, err
	}

	return &VariantSet{OG: og, RE: re}, nil
}
