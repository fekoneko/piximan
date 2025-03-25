package downloader

type ImageSize uint8

const (
	ImageSizeThumbnail ImageSize = 0
	ImageSizeSmall     ImageSize = 1
	ImageSizeMedium    ImageSize = 2
	ImageSizeOriginal  ImageSize = 3
)

const ImageSizeDefault = ImageSizeOriginal

func ImageSizeFromUint(size uint8) ImageSize {
	if size <= 3 {
		return ImageSize(size)
	}
	return ImageSizeDefault
}
