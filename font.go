package main

import (
	"fmt"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var fontCache = map[string]*opentype.Font{}

func loadFace(filename string, size float64) (font.Face, error) {
	loadedFont, ok := fontCache[filename]
	if !ok {
		var err error
		loadedFont, err = LoadFont(filename)
		if err != nil {
			return nil, fmt.Errorf("load font: %w", err)
		}

		fontCache[filename] = loadedFont
	}

	fontFace, err := opentype.NewFace(loadedFont, &opentype.FaceOptions{
		Size:    size,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("create font face: %w", err)
	}

	return fontFace, nil
}
