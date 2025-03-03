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
ImageGetGreyscaledChannel makes a float32 slice. Slice values represent a greyscaled version of input image. 
The returned slice assumes row major order.
Greyscale algorithm used: Y = 0.299 * R + 0.587 * G + 0.114 *B
*/
func ImageGetGreyscaledChannel(img *image.RGBA) []float32 {
	var pixel_count int = img.Rect.Dy() * img.Rect.Dx()
	greyscaled := make([]float32, pixel_count)

	for y := 0 ; y < img.Rect.Dy() ; y++ {
		for x := 0 ; x < img.Rect.Dx() ; x++ {
			flat_id := img.PixOffset(x + img.Rect.Min.X, y + img.Rect.Min.Y)
			greyscale_id := y * img.Rect.Dx() + x

			greyscaled[greyscale_id] = (
				0.299 * float32(img.Pix[flat_id + 0]) +
				0.587 * float32(img.Pix[flat_id + 1]) +
				0.114 * float32(img.Pix[flat_id + 2]))
		}
	}


	return greyscaled

}

/*
	ImageGetNormalized makes an entirely new RGBA image that represets the same image, but the rectangle starts at 0,0 
	and stride is equal to deltaX*4 (underlying slice has only the necessary memory)
*/
func ImageGetNormalized(img *image.RGBA) *image.RGBA {
	new_rect := image.Rect(0,0,img.Rect.Dx(),img.Rect.Dy())
	new_stride := img.Rect.Dx() * 4
	new_data := make([]uint8, img.Rect.Dx() * img.Rect.Dy() * 4)

	for y := 0 ; y < img.Rect.Dy() ; y++ {
		flat_id_origin := img.PixOffset(img.Rect.Min.X, y + img.Rect.Min.Y)
		flat_id_target := y * new_stride 
		copy(
			new_data[flat_id_target: flat_id_target + new_stride],
			img.Pix[flat_id_origin: flat_id_origin + new_stride],
		)

	}
	
	return & image.RGBA{
		Pix : new_data,
		Stride: new_stride,
		Rect: new_rect,
	}
}


/*
	ImageGetTransposed makes an entirely new RGBA image that is a transposed version of input image. 
	Resulting image is always normalized, see ImageGetNormalized for more info.
*/
func ImageGetTransposed(img *image.RGBA) *image.RGBA {
	new_rect := image.Rect(0,0,img.Rect.Dy(), img.Rect.Dx())
	new_stride := img.Rect.Dy() * 4
	new_data := make([]uint8, img.Rect.Dx() * img.Rect.Dy() * 4)

	for y := 0 ; y < img.Rect.Dy() ; y++ {
		for x:= 0; x < img.Rect.Dx(); x++ {
			flat_id_origin := img.PixOffset(x + img.Rect.Min.X, y + img.Rect.Min.Y)
			flat_id_target := (x * img.Rect.Dy() + y) * 4
			copy(
				new_data[flat_id_target: flat_id_target + 4],
				img.Pix[flat_id_origin: flat_id_origin + 4],
			)
		}
	}
	
	return & image.RGBA{
		Pix : new_data,
		Stride: new_stride,
		Rect: new_rect,
	}
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


