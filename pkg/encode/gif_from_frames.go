package encode

import (
	"archive/zip"
	"bytes"
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
)

type Frame struct {
	Filename string
	Duration int
}

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

		// NOTE: assume the files in the archive to be of JPEG format
		decodedImage, err := jpeg.Decode(reader)
		if err != nil {
			return nil, err
		}

		bounds := decodedImage.Bounds()
		palettedImage := image.NewPaletted(bounds, palette.Plan9)
		draw.Draw(palettedImage, bounds, decodedImage, bounds.Min, draw.Src)
		encodedGif.Image = append(encodedGif.Image, palettedImage)
		encodedGif.Delay = append(encodedGif.Delay, frame.Duration)
	}
	writer := bytes.NewBuffer(nil)
	gif.EncodeAll(writer, &encodedGif)

	return writer.Bytes(), nil
}
