package contrast

import (
	"image"
	"fmt"
	"math"
)

import (
	"pixel_restoration/common"
)



/*
CalculatePixelEdgeDistances calculates euclidean distance between each adjacent pair of pixels in all rows or colums.
Vertical parameter determines for which axis the distances are computed.

If vertical == false, 
	result is a slice containing Ysize slices of distances between Xsize-1 pixel pairs in a row
If vertical == true, 
	result is a slice containing Xsize slices of distances between Ysize-1 pixel pairs in a column

All inner slices in the returned structure are guaranteed to be allocated in a single memory block,
 which can be reinterpreted as flat slice by reslicing first element to size of total number of elements
*/
func CalculatePixelEdgeDistances(img *image.RGBA, vertical bool) [][]float32{
	var is_vertical int = common.Ternary(vertical, 1, 0)
	sizes := [2]int{
		img.Rect.Dx(),
		img.Rect.Dy(),
	}
	result := common.Make2D[float32](sizes[1 - is_vertical], sizes[is_vertical] - 1)

	for outer := 0; outer < sizes[1 - is_vertical] ; outer++ {
		for inner := 0; inner < sizes[is_vertical] - 1 ; inner++ {
			fmt.Println(outer, inner)
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
			dist := float32(math.Sqrt(r_delta * r_delta + g_delta *g_delta + b_delta * b_delta))
			result[outer][inner] = dist
		}

	}

	return result


}