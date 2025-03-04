package contrast

import (
	"image"
	"math"
)

import (
	"pixel_restoration/common"
)



/*
CalculatePixelEdgeDistances calculates euclidean distance between each adjacent pair of pixels in all rows or colums.
Vertical parameter determines for which axis the distances are computed.

If vertical == false, 
	result is an image containing Ysize rows of distances between Xsize-1 pixel pairs in a row Y
If vertical == true, 
	result is an image containing Xsize rows of distances between Ysize-1 pixel pairs in a column X

First value of each row is 0 padding used to preserve image shape

*/
func CalculatePixelEdgeDistances(img *image.RGBA, vertical bool) *image.Gray{
	var is_vertical int = common.Ternary(vertical, 1, 0)
	sizes := [2]int{
		img.Rect.Dx(),
		img.Rect.Dy(),
	}

	height, width := sizes[1 - is_vertical], sizes[is_vertical]
	new_rect := image.Rect(0,0,width, height)
	new_stride := width
	new_data := make([]uint8, width * height)


	for outer := 0; outer < height ; outer++ {
		for inner := 0; inner < width - 1 ; inner++ {
			curr := [2]int {inner, outer}
			next := [2]int {inner + 1, outer}

			curr_flat_id := img.PixOffset(
				curr[is_vertical] 	  + img.Rect.Min.X,
			 	curr[1 - is_vertical] + img.Rect.Min.Y,
			)
			next_flat_id := img.PixOffset(
				next[is_vertical] 	  + img.Rect.Min.X,
				next[1 - is_vertical] + img.Rect.Min.Y,
			)

			var r_delta float64 = float64(img.Pix[curr_flat_id + 0]) - float64(img.Pix[next_flat_id + 0])
			var g_delta float64 = float64(img.Pix[curr_flat_id + 1]) - float64(img.Pix[next_flat_id + 1])
			var b_delta float64 = float64(img.Pix[curr_flat_id + 2]) - float64(img.Pix[next_flat_id + 2])
			dist := math.Sqrt(r_delta * r_delta + g_delta *g_delta + b_delta * b_delta)

			new_data[outer * new_stride + inner + 1] = distMapToUint8(dist)
		}

	}

	return & image.Gray{
		Pix : new_data,
		Stride: new_stride,
		Rect: new_rect,
	}


}


func distMapToUint8(dist float64) uint8 {
	const max_possible_color_diff = 441.674
	return uint8(255.0 * dist / max_possible_color_diff + 0.5)
}