package dto

type FramesData struct {
	Src    string  `json:"src"`
	Frames []Frame `json:"frames"`
}

type Frame struct {
	File  string `json:"file"`
	Delay int    `json:"delay"`
}
