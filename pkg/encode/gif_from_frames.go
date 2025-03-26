package encode

import (
	"archive/zip"
	"bytes"
)

type Frame struct {
	Filename string `json:"file"`
	Duration uint64 `json:"delay"`
}

func GifFromFrames(archive []byte, frames []Frame) ([]byte, error) {
	reader, err := zip.NewReader(bytes.NewReader(archive), int64(len(archive)))
	if err != nil {
		return nil, err
	}

	for _, file := range reader.File {
		reader, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer reader.Close()

		// TODO: implement gif encoding
	}

	// TODO: return gif bytes
	return nil, nil
}
