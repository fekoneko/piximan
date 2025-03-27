package dto

import "github.com/fekoneko/piximan/pkg/encode"

type FramesData struct {
	Src    string  `json:"src"`
	Frames []Frame `json:"frames"`
}

type Frame struct {
	File  string `json:"file"`
	Delay int    `json:"delay"`
}

func (f *FramesData) FromDto() (string, []encode.Frame) {
	frames := make([]encode.Frame, len(f.Frames))
	for i, frame := range f.Frames {
		frames[i] = encode.Frame{
			Filename: frame.File,
			Duration: frame.Delay / 10,
		}
	}

	return f.Src, frames
}
