package images

import (
	"image"
	"slices"
)

/*
	Create greyscale image from slice of float values.
	max_float_val specifies the lowest value that will be mapped to 255
	if parameter is not positive, its calculated automatically based on highest value in slice

*/
func GreyscaleImageFromFloatSlice(slice []float32, shape [2]int, max_float_val float32) *image.Gray{
	if(max_float_val <= 0){
		max_float_val = slices.Max(slice)
	}

	greyscale_data := make([]uint8, shape[1] * shape[0])
	greyscale_stride := shape[1]
	greyscale_rect := image.Rect(0,0, shape[1], shape[0])


	for id , float_val := range slice {
		float_val = min(float_val, max_float_val)
		as_uint8 := uint8(255.0 * float_val/max_float_val)
		greyscale_data[id] = as_uint8
	}


	return & image.Gray{
		Pix : greyscale_data,
		Stride: greyscale_stride,
		Rect: greyscale_rect,
	}
}
