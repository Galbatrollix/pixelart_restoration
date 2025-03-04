/*
	Arrangements file contains functions that alter the image files in simple mathematical ways
	- splitting to separate channels
	- combining image from separate channels <not implemented>
	- converting to grayscale
	- transposing images
	- flipping images <not implemented>
*/

package images

import (
	"image"
)


/* 
	ImageGetSplitChannels seturns array of 4 gray image pointers. Each gray image corresponds to 
	a separate channel of RGBA image. Order: R, G, B, A
*/
func ImageGetSplitChannels(img *image.RGBA) [4]*image.Gray {
	height, width := img.Rect.Dy(), img.Rect.Dx()

	var result [4]*image.Gray
	for i:=0; i<4; i++{
		new_data := make([]uint8, height * width)
		new_stride := width
		new_rect := image.Rect(0, 0, width, height)

		result[i] = & image.Gray{
			Pix : new_data,
			Stride: new_stride,
			Rect: new_rect,
		}
	}

	for y := 0 ; y < height ; y++ {
		for x := 0 ; x < width ; x++ {
			flat_id := img.PixOffset(x + img.Rect.Min.X, y + img.Rect.Min.Y)
			id_channel := y * width + x

			result[0].Pix[id_channel] = img.Pix[flat_id + 0]
			result[1].Pix[id_channel] = img.Pix[flat_id + 1]
			result[2].Pix[id_channel] = img.Pix[flat_id + 2]
			result[3].Pix[id_channel] = img.Pix[flat_id + 3]
		}

	}
	return result
}



/*
ImageGetGrayscaled turns an RGBG image into grayscale, returns a pointer to new image.Gray.
Grayscale algorithm used: Result = 0.299*R + 0.587*G + 0.114*B
*/
func ImageGetGrayscaled(img *image.RGBA) *image.Gray {
	height, width := img.Rect.Dy(), img.Rect.Dx()
	new_rect := image.Rect(0,0,width, height)
	new_stride := width
	new_data := make([]uint8, width * height)

	
	for y := 0 ; y < height ; y++ {
		for x := 0 ; x < width ; x++ {
			flat_id := img.PixOffset(x + img.Rect.Min.X, y + img.Rect.Min.Y)
			grayscale_id := y * width + x

			new_data[grayscale_id] = uint8(
				0.299 * float32(img.Pix[flat_id + 0]) +
				0.587 * float32(img.Pix[flat_id + 1]) +
				0.114 * float32(img.Pix[flat_id + 2]) + 0.5)  // 0.5 is to round properly by truncating.

		}
	}


	return & image.Gray{
		Pix : new_data,
		Stride: new_stride,
		Rect: new_rect,
	}
}


/*
	ImageGetNormalized makes an entirely new RGBA image that represets the same image, but the rectangle starts at 0,0 
	and stride is equal to deltaX*4 (underlying slice has only the necessary memory)
*/
func ImageGetNormalized(img *image.RGBA) *image.RGBA {
	height, width := img.Rect.Dy(), img.Rect.Dx()
	new_rect := image.Rect(0, 0, width, height)
	new_stride := width * 4
	new_data := make([]uint8, width * height * 4)

	for y := 0 ; y < height ; y++ {
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
	ImageGetNormalized makes an entirely new Gray image that represets the same image, but the rectangle starts at 0,0 
	and stride is equal to deltaX (underlying slice has only the necessary memory)
*/
func GrayscaleGetNormalized(img *image.Gray) *image.Gray {
	height, width := img.Rect.Dy(), img.Rect.Dx()
	new_rect := image.Rect(0, 0, width, height)
	new_stride := width
	new_data := make([]uint8, width * height)

	for y := 0 ; y < height ; y++ {
		flat_id_origin := img.PixOffset(img.Rect.Min.X, y + img.Rect.Min.Y)
		flat_id_target := y * new_stride 
		copy(
			new_data[flat_id_target: flat_id_target + new_stride],
			img.Pix[flat_id_origin: flat_id_origin + new_stride],
		)

	}

	return & image.Gray{
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
	height, width := img.Rect.Dy(), img.Rect.Dx()
	new_rect := image.Rect(0,0, height, width)
	new_stride := height * 4
	new_data := make([]uint8, width * height * 4)

	for y := 0 ; y < height ; y++ {
		for x:= 0; x < width; x++ {
			flat_id_origin := img.PixOffset(x + img.Rect.Min.X, y + img.Rect.Min.Y)
			flat_id_target := (x * height + y) * 4
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
	GrayscaleGetGrayScaled makes an entirely new Gray image that is a transposed version of input image. 
	Resulting image is always normalized, see ImageGetNormalized for more info.
*/
func GrayscaleGetTransposed(img *image.Gray) *image.Gray {
	height, width := img.Rect.Dy(), img.Rect.Dx()
	new_rect := image.Rect(0,0, height, width)
	new_stride := height
	new_data := make([]uint8, width * height)

	for y := 0 ; y < height ; y++ {
		for x:= 0; x < width; x++ {
			flat_id_origin := img.PixOffset(x + img.Rect.Min.X, y + img.Rect.Min.Y)
			flat_id_target := (x * height + y)
			new_data[flat_id_target] = img.Pix[flat_id_origin]
		}
	}
	
	return & image.Gray{
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


