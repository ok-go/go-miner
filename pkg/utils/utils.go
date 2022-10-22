package utils

import (
	"github.com/golang/freetype/truetype"
	goMiner "go-miner"
	"golang.org/x/image/font"
	"path"
)

func LoadTTF(name string, size float64) (font.Face, error) {
	bytes, err := goMiner.Ttf.ReadFile(path.Join("assets", name))
	if err != nil {
		return nil, err
	}

	f, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(f, &truetype.Options{
		GlyphCacheEntries: 1,
		Size:              size,
	}), nil
}
