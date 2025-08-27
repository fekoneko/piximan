package imageext

type Size uint8

const (
	SizeThumbnail Size = 0
	SizeSmall     Size = 1
	SizeMedium    Size = 2
	SizeOriginal  Size = 3
	SizeDefault        = SizeOriginal
)

func ValidSizeUint(size uint64) bool {
	return size <= 3
}

func SizeFromUint(size uint64) Size {
	if ValidSizeUint(size) {
		return Size(size)
	}
	return SizeDefault
}

func (s Size) ToUint() uint64 {
	return uint64(s)
}
