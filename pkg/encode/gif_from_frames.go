package encode

import (
	"archive/zip"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"

	"github.com/ericpauley/go-quantize/quantize"
)

type Frame struct {
	Filename string
	Duration int
}

var quantizer = quantize.MedianCutQuantizer{}

// TODO: better quality for GIFs
func GifFromFrames(archive []byte, frames []Frame) ([]byte, error) {
	reader, err := zip.NewReader(bytes.NewReader(archive), int64(len(archive)))
	if err != nil {
		return nil, err
	}

	filesLookup := make(map[string]*zip.File)
	for _, file := range reader.File {
		filesLookup[file.Name] = file
	}

	encodedGif := gif.GIF{}
	for _, frame := range frames {
		file, ok := filesLookup[frame.Filename]
		if !ok {
			return nil, fmt.Errorf("frame %v not found in archive", frame.Filename)
		}

		reader, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer reader.Close()

		decodedImage, err := jpeg.Decode(reader)
		if err != nil {
			return nil, err
		}

		bounds := decodedImage.Bounds()
		plt := quantizer.Quantize(make([]color.Color, 0, 256), decodedImage)
		palettedImage := image.NewPaletted(bounds, plt)
		draw.Draw(palettedImage, bounds, decodedImage, bounds.Min, draw.Src)
		encodedGif.Image = append(encodedGif.Image, palettedImage)
		encodedGif.Delay = append(encodedGif.Delay, frame.Duration)
	}
	writer := bytes.NewBuffer(nil)
	gif.EncodeAll(writer, &encodedGif)

	return writer.Bytes(), nil
}
