package dto

import "github.com/fekoneko/piximan/pkg/encode"

type FramesData struct {
	Src    string         `json:"src"`
	Frames []encode.Frame `json:"frames"`
}
