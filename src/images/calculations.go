package images
import (
	"image"
)

import (
	"pixel_restoration/common"
)


func CalculateAverageColor(img *image.RGBA) [3]uint8 {
	var result = [3]uint8{}
	channels := ImageGetSplitChannels(img)

	for i := 0; i<3; i++ {
		result[i] = common.MeanOfSliceU8(channels[i].Pix)
	}

	return result
}

func CalculateMedianColor(img *image.RGBA) [3]uint8{
	var result = [3]uint8{}
	channels := ImageGetSplitChannels(img)

	for i := 0; i<3; i++ {
		result[i] = common.MedianOfSliceU8(channels[i].Pix)
	}

	return result
}

