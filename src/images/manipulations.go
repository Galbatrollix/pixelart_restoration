package images

import (
	"image"
)


/* 
	ImageGetSplitChannels seturns array of 3 slices of uint8, each slice has different memory block
	each uint8 slice corresponds to flattened R, G, B channel of the image, alpha is discarded
*/
func ImageGetSplitChannels(img *image.RGBA) [3][]uint8 {
	var channel_size int = img.Rect.Dy() * img.Rect.Dx()

	channels := [3][]uint8{ 
		make([]uint8, channel_size),
		make([]uint8, channel_size),
		make([]uint8, channel_size),
	}

	for y := 0 ; y < img.Rect.Dy() ; y++ {
		for x := 0 ; x < img.Rect.Dx() ; x++ {
			flat_id := img.PixOffset(x + img.Rect.Min.X, y + img.Rect.Min.Y)
			id_channel := y * img.Rect.Dx() + x

			channels[0][id_channel] = img.Pix[flat_id + 0]
			channels[1][id_channel] = img.Pix[flat_id + 1]
			channels[2][id_channel] = img.Pix[flat_id + 2]
		}

	}
	return channels
}



/*
	ImageGetNormalized makes an entirely new RGBA image that represets the same image, but the rectangle starts at 0,0 
	and stride is equal to deltaX*4 (underlying slice has only the necessary memory)
*/
func ImageGetNormalized(img *image.RGBA) *image.RGBA {
	panic("NOT IMPLEMENTED")
	return nil
}


/*
	ImageGetTransposed makes an entirely new RGBA image that is a transposed version of input image. 
	Resulting image must be normalized, see ImageGetNormalized for more info.
*/
func ImageGetTransposed(img *image.RGBA) *image.RGBA {
	panic("NOT IMPLEMENTED")
	return nil
}

/*
	ImageFlipVertical flips RGBA image pixels along the Y axis 
*/
func ImageFlipVertical(img *image.RGBA) {
	panic("NOT IMPLEMENTED")
}

/*
	ImageFlipHorizontal flips RGBA image pixels along the X axis 
*/
func ImageFlipHorizontal(img *image.RGBA) {
	panic("NOT IMPLEMENTED")
}


