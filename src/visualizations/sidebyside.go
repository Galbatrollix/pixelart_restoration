package visualizations


import "image"

import "pixel_restoration/images"

/*
	Makes new grayscale image out of two provided grayscale images

	As long as provided images have width < 2 * height, images are arranged side by side
	Otherwise images are arranged up and below

	Input images must have the same dimensions

*/
func SideBySideGrayscale(img1, img2 *image.Gray) *image.Gray{
	if img1.Rect.Dx() != img2.Rect.Dx() || img1.Rect.Dy() != img2.Rect.Dy(){
		panic("Different sized grayscale images provided")
	}

	img1 = images.GrayscaleGetNormalized(img1)
	img2 = images.GrayscaleGetNormalized(img2)

	height, width := img1.Rect.Dy(), img1.Rect.Dx()
	if  width > height * 2 {
		return combineVertically(img1, img2)
	}else{
		return combineHorizontally(img1, img2)
	}


}


func combineVertically(img1, img2 *image.Gray) *image.Gray{
	height, width := img1.Rect.Dy() * 2, img1.Rect.Dx()

	new_rect := image.Rect(0,0, width, height)
	new_stride := width
	new_data := make([]uint8, width * height)

	subimage_pixel_count := img1.Rect.Dy() * img1.Rect.Dx()

	copy(new_data[0:], img1.Pix)
	copy(new_data[subimage_pixel_count:], img2.Pix)

	return & image.Gray{
		Pix : new_data,
		Stride: new_stride,
		Rect: new_rect,
	}

}

func combineHorizontally(img1, img2 *image.Gray) *image.Gray{
	height, width := img1.Rect.Dy(), img1.Rect.Dx() * 2

	new_rect := image.Rect(0,0, width, height)
	new_stride := width
	new_data := make([]uint8, width * height)

	old_stride := img1.Rect.Dx()
	for y:=0; y<height; y++{
		copy(
			new_data[y*new_stride:],
			img1.Pix[y*old_stride: (y+1) * old_stride],
		)
		copy(
			new_data[y*new_stride + old_stride:],
			img2.Pix[y*old_stride: (y+1) * old_stride],
		)
	}


	return & image.Gray{
		Pix : new_data,
		Stride: new_stride,
		Rect: new_rect,
	}

}