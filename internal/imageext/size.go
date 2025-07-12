package imageext

type Size uint8

const (
	SizeThumbnail Size = 0
	SizeSmall     Size = 1
	SizeMedium    Size = 2
	SizeOriginal  Size = 3
	SizeDefault        = SizeOriginal
)

func SizeFromUint(size uint) Size {
	if size <= 3 {
		return Size(size)
	}
	return SizeDefault
}
