package visualizations


import "image"

import "pixel_restoration/images"

/* 
	Creates copy of an image and draws gridlines at selected position indexes with given color
*/
func ImageWithDrawnGridlinesSimple(img *image.RGBA, indexes[2][]int, color [4]uint8) *image.RGBA {
	new_img := images.ImageGetNormalized(img)

	images.DrawGridlineRowsOnImage(new_img, indexes[0], color)
	images.DrawGridlineColsOnImage(new_img, indexes[1], color)

	return new_img
}

/*
	 Creates an upscaled copy of an image and draws gridlines after selected positions of original image
*/
func ImageWithDrawnGridlinesAdvanced(img *image.RGBA, indexes[2][]int, color [4]uint8) *image.RGBA{
	const scaling_factor = 6

	img_big := images.ImageUpscaledByFactor(img, scaling_factor)

	// Drawing all edges in black
	all_indexes := [2][]int{
		indexesAllScaled(img.Rect.Dy(), scaling_factor),
		indexesAllScaled(img.Rect.Dx(), scaling_factor),
	}
	color_black := [4]uint8{0,0,0,255}
	images.DrawGridlineRowsOnImage(img_big, all_indexes[0], color_black)
	images.DrawGridlineColsOnImage(img_big, all_indexes[1], color_black)

	// Drawing selected edges in provided color
	indexes = [2][]int{
		indexesConvertToScaled(indexes[0], scaling_factor),
		indexesConvertToScaled(indexes[1], scaling_factor),
	}

	images.DrawGridlineRowsOnImage(img_big, indexes[0], color)
	images.DrawGridlineColsOnImage(img_big, indexes[1], color)

	// cutting last pixel out since its garbage empty gridline

	img_big = img_big.SubImage(
		image.Rect(
			0, 0, img_big.Rect.Dx() - 1, img_big.Rect.Dy() - 1,
		),
	).(*image.RGBA)

	return img_big

}


func indexesConvertToScaled(indexes []int, scaling_factor int) []int{
	result := make([]int, len(indexes))
	for i, value := range indexes{
		result[i] = value * scaling_factor - 1
	}
	return result
}	

func indexesAllScaled(dimension int, scaling_factor int) []int {
	result := make([]int, dimension)
	for i := 1 ; i <= dimension; i++{
		result[i-1] = i * scaling_factor - 1
	}
	return result
}