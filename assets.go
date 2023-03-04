package main

import (
	"embed"
	"fmt"
	"io"
	"io/fs"

	"golang.org/x/image/font/opentype"
)

const (
	JetBrainsMonoFont = "JetBrainsMono-Regular.ttf"
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
