package dto

import (
	"github.com/fekoneko/piximan/internal/imageext"
	"github.com/fekoneko/piximan/internal/utils"
)

type FramesData struct {
	Src    *string `json:"src"`
	Frames []Frame `json:"frames"`
}

type Frame struct {
	File  *string `json:"file"`
	Delay *int    `json:"delay"`
}

func (f *FramesData) FromDto() (framesSrc *string, frames *[]imageext.Frame) {
	frames = utils.ToPtr(make([]imageext.Frame, len(f.Frames)))
	for i, frame := range f.Frames {
		if frame.File == nil || frame.Delay == nil {
			return f.Src, nil
		}

		(*frames)[i] = imageext.Frame{
			Filename: *frame.File,
			Duration: *frame.Delay / 10,
		}
	}

	return f.Src, frames
}
