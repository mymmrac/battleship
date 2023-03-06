package main

import (
	"embed"
	"fmt"
	"image"
	"io"
	"io/fs"

	"golang.org/x/image/font/opentype"
)

const (
	JetBrainsMonoFont = "JetBrainsMono-Regular.ttf"

	Icon16x16 = "icon-16x16.png"
	Icon32x32 = "icon-32x32.png"
	Icon48x48 = "icon-48x48.png"
)

//go:embed assets
var assetDirFS embed.FS

var assetFS fs.FS

func init() {
	var err error
	assetFS, err = fs.Sub(assetDirFS, "assets")
	if err != nil {
		panic(fmt.Errorf("load assets dir: %w", err))
	}
}

func LoadData(filename string) ([]byte, error) {
	file, err := assetFS.Open(filename)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func LoadFont(filename string) (*opentype.Font, error) {
	data, err := LoadData(filename)
	if err != nil {
		return nil, err
	}

	fontFace, err := opentype.Parse(data)
	if err != nil {
		return nil, err
	}

	return fontFace, nil
}

func LoadImage(filename string) (image.Image, error) {
	file, err := assetFS.Open(filename)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}
