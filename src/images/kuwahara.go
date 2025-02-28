package images

import (
	"image"
	"math"
	_"fmt"
)

import (
	"pixel_restoration/common"
)

func KuwaharaGaussian(image *image.RGBA, radius int, sigma float32){
	if radius < 1 {
		panic("Radius must be bigger than 0")
	}
	panic("NOT IMPLEMENTED")
}	



func KuwaharaMean(image *image.RGBA, radius int){
	if radius < 1 {
		panic("Radius must be bigger than 0")
	}
	panic("NOT IMPLEMENTED")
}


/*

	Computes flat gaussian, implementation based on opencv
	https://docs.opencv.org/4.x/d4/d86/group__imgproc__filter.html#gac05a120c1ae92a6060dd0db190a61afa

*/
func GaussianKernel1D(ksize int, sigma float32) []float32 {
	if ksize < 1 {
		panic("Non positive ksize provided to gaussianKernel1D")
	}

	// calculates sigma automatically if not provided with valid value
	if sigma <= 0 {
		sigma = 0.3*(float32(ksize-1)*0.5 - 1.0) + 0.8
	}
	var denominator = - (2.0 * sigma * sigma)

	var total_sum float32 = 0.0
	kernel := make([]float32, ksize)
	for i:=0 ; i<ksize ; i++{
		nominator_sqrt := (float32(i) - (float32(ksize) - 1.0)/2.0)
		nominator := nominator_sqrt * nominator_sqrt
		full_value := float32(math.Exp(float64(nominator / denominator)))
		total_sum += full_value
		kernel[i] = full_value
	}

	// dividing everything by total sum to get sum equal to 1
	for i:=0 ; i<ksize ; i++{
		kernel[i] /= total_sum
	}

	return kernel
}
/*
SepFilter2D applies a separable linear filter to the single channel image. 


*/
func SepFilter2D(img [][]float32, kernels [2][]float32, kernel_anchors [2]int) [][]float32{
	if len(img) < len(kernels[0]) || len(img[0]) < len(kernels[1]) {
		panic("Kernel dimension larger than image dimension. Not implemented.")
	}
	intermediate := filterHorizontal1D(img, kernels[1], kernel_anchors[1])
	result := filterVertical1D(intermediate, kernels[0], kernel_anchors[0])

	return result
}


func filterVertical1D(img [][]float32, kernel []float32, kernel_anchor int) [][]float32{
	result := common.Make2D[float32](len(img), len(img[0]))

	// offsets say which row/ columns from the start/end where not all kernel values are in range of image
	kernel_offset_L := kernel_anchor
	kernel_offset_R := len(kernel) - 1 - kernel_anchor

	KernelRange := [2]int{-kernel_offset_L, kernel_offset_R + 1}
	Xrange := [2]int{0, len(img[0])}
	// ranges for Y in each of three loops
	Yranges := [3][2]int{
		{0, kernel_offset_L},
		{kernel_offset_L, len(img) - kernel_offset_R},
		{len(img) - kernel_offset_R, len(img)},
	}

	// first loop - left hands side kernel positions are out of bounds
	for y := Yranges[0][0]; y<Yranges[0][1]; y++{
		for x := Xrange[0]; x<Xrange[1]; x++ {
			var sum float32 = 0.0
			for y_offset := KernelRange[0]; y_offset < KernelRange[1] ; y_offset++ {
				kernel_weight := kernel[kernel_anchor + y_offset]
				img_y := y + y_offset
				img_y_reflected := common.Ternary(img_y < 0, -img_y, img_y)
				sum += img[img_y_reflected][x] * kernel_weight
			}
			result[y][x] = sum
		}
	}


	// middle loop - all pixels accessible
	for y := Yranges[1][0]; y<Yranges[1][1]; y++{
		for x := Xrange[0]; x<Xrange[1]; x++ {
			var sum float32 = 0.0
			for y_offset := KernelRange[0]; y_offset < KernelRange[1] ; y_offset++ {
				kernel_weight := kernel[kernel_anchor + y_offset]
				sum += img[y + y_offset][x] * kernel_weight
			}
			result[y][x] = sum
		}
	}


	// third loop - right hands side kernel positions are out of bounds
	for y := Yranges[2][0]; y<Yranges[2][1]; y++{
		for x := Xrange[0]; x<Xrange[1]; x++ {
			var sum float32 = 0.0
			for y_offset := KernelRange[0]; y_offset < KernelRange[1] ; y_offset++ {
				kernel_weight := kernel[kernel_anchor + y_offset]
				img_y := y + y_offset
				img_y_reflected := common.Ternary(
					img_y > Yranges[2][1] - 1,
					(Yranges[2][1] - 1) - img_y + (Yranges[2][1] - 1),
					img_y,
				)
				sum += img[img_y_reflected][x] * kernel_weight
			}
			result[y][x] = sum
		}
	}


	return result
}

func filterHorizontal1D(img [][]float32, kernel []float32, kernel_anchor int) [][]float32{
	result := common.Make2D[float32](len(img), len(img[0]))

	// offsets say which row/ columns from the start/end where not all kernel values are in range of image
	kernel_offset_L := kernel_anchor
	kernel_offset_R := len(kernel) - 1 - kernel_anchor


	KernelRange := [2]int{-kernel_offset_L, kernel_offset_R + 1}
	Yrange := [2]int{0, len(img)}
	// ranges for X in each of three loops
	Xranges := [3][2]int{
		{0, kernel_offset_L},
		{kernel_offset_L, len(img[0]) - kernel_offset_R},
		{len(img[0]) - kernel_offset_R, len(img[0])},
	}

	// first loop - left hands side kernel positions are out of bounds
	for y := Yrange[0]; y<Yrange[1]; y++{
		for x := Xranges[0][0]; x<Xranges[0][1]; x++ {
			var sum float32 = 0.0
			for x_offset := KernelRange[0]; x_offset < KernelRange[1] ; x_offset++ {
				kernel_weight := kernel[kernel_anchor + x_offset]
				// reflect the index from the left edge by taking absolute value of index 
				img_x := x + x_offset
				img_x_reflected := common.Ternary(img_x < 0, -img_x, img_x)
				sum += img[y][img_x_reflected] * kernel_weight
			}
			result[y][x] = sum
		}
	}


	// middle loop - all pixels accessible
	for y := Yrange[0]; y<Yrange[1]; y++{
		for x := Xranges[1][0]; x<Xranges[1][1]; x++ {
			var sum float32 = 0.0
			for x_offset := KernelRange[0]; x_offset < KernelRange[1] ; x_offset++ {
				kernel_weight := kernel[kernel_anchor + x_offset]
				sum += img[y][x + x_offset] * kernel_weight
			}
			result[y][x] = sum
		}
	}


	// third loop - right hands side kernel positions are out of bounds
	for y := Yrange[0]; y<Yrange[1]; y++{
		for x := Xranges[2][0]; x<Xranges[2][1]; x++ {
			var sum float32 = 0.0
			for x_offset := KernelRange[0]; x_offset < KernelRange[1] ; x_offset++ {
				kernel_weight := kernel[kernel_anchor + x_offset]
				// reflect the index from the right edge
				img_x := x + x_offset
				img_x_reflected := common.Ternary(
					img_x > Xranges[2][1] - 1,
					(Xranges[2][1] - 1) - img_x + (Xranges[2][1] - 1),
					img_x,
				)
				sum += img[y][img_x_reflected] * kernel_weight
			}
			result[y][x] = sum
		}
	}
	return result
}
