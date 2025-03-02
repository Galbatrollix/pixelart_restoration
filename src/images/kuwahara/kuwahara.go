package kuwahara

import (
	"image"
	"math"
)

import (
	"pixel_restoration/common"
	"pixel_restoration/images"
)


func KuwaharaGaussian(img *image.RGBA, radius int, sigma float32) *image.RGBA{
	if radius < 1 {
		panic("Radius must be bigger than 0")
	}
	if img.Rect.Dx() * img.Rect.Dy() == 0 {
		panic("Image must have at least one pixel")
	}

	// calculates sigma automatically if not provided with valid value
	if sigma <= 0 {
		sigma = 0.3* (float32(radius) - 1.0) + 0.8
	}

	// making two semi-kernels
	kernel_base := GaussianKernel1D(radius * 2 + 1, sigma)
	kernel_forward := kernel_base[:radius + 1]
	kernel_reverse := kernel_base[radius:]

	// normalizing the semi-kernels so they add up to 1
	var semikernel_sum float32 = sliceSumFloat32(kernel_forward)
	sliceDivbyFloat32(kernel_base, semikernel_sum)

	kernel_quadrants := [4][2][]float32{
		// {kernel y , kernel x}
		{kernel_forward, kernel_forward},
		{kernel_forward, kernel_reverse},
		{kernel_reverse, kernel_reverse},
		{kernel_reverse, kernel_forward},
	}

	kernel_anchors := [4][2]int{
		// {anchor y, anchor x}
		{radius, radius},
		{radius, 0     },
		{0     , 0     },
		{0     , radius},
	} 

	img_shape := [2]int{
		img.Rect.Dy(),
		img.Rect.Dx(),
	}
	total_count := img_shape[0] * img_shape[1]

	greyscale := images.ImageGetGreyscaledChannel(img)
	greyscale_squared := sliceGetSquared(greyscale)

	// reserving space for temporary buffers for SepFilter2D
	temporary := make([]float32, total_count)
	temporary2 := make([]float32, total_count)
	temporary3 := make([]float32, total_count)

	// making space for standard deviations array
	var standard_deviations [4][]float32
	stddevs_buffer := make([]float32, total_count * 4)
	for i:= 0; i<4; i++{
		standard_deviations[i] = stddevs_buffer[i * total_count : (i+1) * total_count]
	}

	// calculating standard deviations of each quadrant
	for kernel_id := 0; kernel_id < 4 ; kernel_id++{
		greyscale_averages := temporary2
		SepFilter2D(
			greyscale, greyscale_averages, temporary,
			img_shape, kernel_quadrants[kernel_id], kernel_anchors[kernel_id],
		)
        SepFilter2D(
			greyscale_squared,standard_deviations[kernel_id], temporary,
			img_shape, kernel_quadrants[kernel_id], kernel_anchors[kernel_id],
		)
		sliceSubtractSquared(
			standard_deviations[kernel_id],
			greyscale_averages,
		)
	}

	// choosing indexes of the quadrants with the lowest variance
	quadrants_chosen := chooseQuadrants(standard_deviations)

	// making space for temporary color averages array
	var color_averages [4][3][]uint8
	var channel_size int = img.Rect.Dy() * img.Rect.Dx()
	color_averages_buffer := make([]uint8, 3 * 4 * channel_size)
	for kernel_id := 0; kernel_id < 4 ; kernel_id++{
		for channel_id := 0; channel_id < 3; channel_id++{
			start := (kernel_id * 3 + channel_id) * channel_size
			end := start + channel_size
			color_averages[kernel_id][channel_id] = color_averages_buffer[start: end]
		}
	}

	// calculating color averages
	channels := images.ImageGetSplitChannels(img)
	channel_float := temporary2
	averaged := temporary3

	for channel_id := 0; channel_id < 3; channel_id++{
		sliceUint8ToFloat32(channels[channel_id], channel_float)
		for kernel_id := 0; kernel_id < 4; kernel_id++ {
		  	SepFilter2D(
				channel_float , averaged, temporary,
				img_shape, kernel_quadrants[kernel_id], kernel_anchors[kernel_id],
			)
			sliceFloat32ToUint8(averaged, color_averages[kernel_id][channel_id])
		}
	}

	// choosing color averages according to quadrants chosen
	new_rect := image.Rect(0,0, img_shape[1], img_shape[0])
	new_stride := img_shape[1] * 4
	new_data := make([]uint8, total_count * 4)
	for flat_id := 0; flat_id < total_count ; flat_id++ {
		flat_id_result := flat_id * 4
		chosen_quadrant := quadrants_chosen[flat_id]
		new_data[flat_id_result + 0] = color_averages[chosen_quadrant][0][flat_id]
		new_data[flat_id_result + 1] = color_averages[chosen_quadrant][1][flat_id]
		new_data[flat_id_result + 2] = color_averages[chosen_quadrant][2][flat_id]
		new_data[flat_id_result + 3] = 255        // alpha channel constant
		
	}

	return & image.RGBA{
		Pix : new_data,
		Stride: new_stride,
		Rect: new_rect,
	}
}	


/*

	Computes flat gaussian, implementation based on opencv
	https://docs.opencv.org/4.x/d4/d86/group__imgproc__filter.html#gac05a120c1ae92a6060dd0db190a61afa

*/
func GaussianKernel1D(ksize int, sigma float32) []float32 {
	if ksize < 1 {
		panic("Non positive ksize provided to gaussianKernel1D")
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


func sliceSumFloat32(arr []float32) float32{
    var sum float32 = 0.0
    for _, value := range arr{
        sum += value
    }
    return sum
}

func sliceDivbyFloat32(arr []float32, divisor float32){
    for id := range arr{
        arr[id] /= divisor
    }
}

func sliceGetSquared(arr []float32) []float32{
	result := make([]float32, len(arr))
	for id := range arr {
		result[id] = arr[id] * arr[id]
	}
	return result
}

func sliceSubtractSquared(arr []float32, other []float32){
	for id := range arr {
		arr[id] -= other[id] * other[id]
	}
}

func chooseQuadrants(standard_deviations [4][]float32) []uint8{
	item_count := len(standard_deviations[0])
	quadrants_chosen := make([]uint8, item_count)
	for i:= 0; i<item_count; i++{
		var quadrant_id uint8
		var min_id uint8 = 0
		var min_deviation float32 = standard_deviations[0][i]
		for quadrant_id = 1; quadrant_id< 4 ; quadrant_id++{
			if standard_deviations[quadrant_id][i] < min_deviation {
				min_deviation = standard_deviations[quadrant_id][i]
				min_id = quadrant_id
			}
		}
		quadrants_chosen[i] = min_id
	}

	return quadrants_chosen
}


func sliceFloat32ToUint8(floats []float32, uints []uint8){
	for id := range floats {
		uints[id] = uint8(floats[id] + 0.5)
	}
}

func sliceUint8ToFloat32(uints []uint8, floats []float32){
	for id := range floats {
		floats[id] = float32(uints[id])
	}
}


/*
SepFilter2D applies a separable linear filter to the single channel image. 
Function is insipered by similarly named function in OpenCV

*/
func SepFilter2D(img []float32, dest[]float32, temp_buffer[]float32,
			 	shape [2]int,  kernels [2][]float32, kernel_anchors [2]int){
	if len(kernels[1]) <= shape[1]{
		filterHorizontal1D(img, temp_buffer, shape, kernels[1], kernel_anchors[1])
	}else{
		filterEdgeCaseHorizontal1D(img, temp_buffer, shape, kernels[1], kernel_anchors[1])
	}

	if len(kernels[0]) <=  shape[0]{
		filterVertical1D(temp_buffer, dest, shape, kernels[0], kernel_anchors[0])
	}else{
		filterEdgeCaseVertical1D(temp_buffer, dest, shape, kernels[0], kernel_anchors[0])
	}

}
/* this is done if kernel size is larger than dimension*/

func filterEdgeCaseVertical1D(img []float32, result[]float32, shape [2]int, kernel []float32, kernel_anchor int){
	var y_shape int = shape[0]
	var x_shape int = shape[1]

	// offsets say which row/ columns from the start/end where not all kernel values are in range of image
	kernel_offset_L := kernel_anchor
	kernel_offset_R := len(kernel) - 1 - kernel_anchor
	KernelRange := [2]int{-kernel_offset_L, kernel_offset_R + 1}

	for y:=0; y < y_shape; y++{
		for x:=0; x < x_shape; x++{
			var sum float32 = 0.0
			for y_offset := KernelRange[0]; y_offset < KernelRange[1] ; y_offset++ {
				kernel_weight := kernel[kernel_anchor + y_offset]
				img_y := y + y_offset
				img_y_reflected := reflectIndex101(img_y, y_shape-1)
				img_flat_id := img_y_reflected * x_shape + x
				sum += img[img_flat_id] * kernel_weight
			}
			result[y * x_shape + x] = sum
		}
	}

}

func filterEdgeCaseHorizontal1D(img []float32, result[]float32, shape [2]int, kernel []float32, kernel_anchor int){
	var y_shape int = shape[0]
	var x_shape int = shape[1]

	// offsets say which row/ columns from the start/end where not all kernel values are in range of image
	kernel_offset_L := kernel_anchor
	kernel_offset_R := len(kernel) - 1 - kernel_anchor
	KernelRange := [2]int{-kernel_offset_L, kernel_offset_R + 1}

	for y:=0; y < y_shape; y++{
		for x:=0; x < x_shape; x++{
			var sum float32 = 0.0
			for x_offset := KernelRange[0]; x_offset < KernelRange[1] ; x_offset++ {
				kernel_weight := kernel[kernel_anchor + x_offset]
				img_x := x + x_offset
				img_x_reflected := reflectIndex101(img_x, x_shape-1)
				img_flat_id := y * x_shape + img_x_reflected
				sum += img[img_flat_id] * kernel_weight
			}
			result[y * x_shape + x] = sum
		}
	}
}

/* used by edge case variants of filter, can reflect off walls multiple times */
func reflectIndex101(index_check, max_index int) int {
	remainder := index_check % max(max_index * 2, 1)
	if remainder < 0 {
		remainder = - remainder
	}
	var result int = 0
	if(remainder > max_index){
		result = - remainder + max_index * 2
	}else{
		result =  remainder
	}
	return result
}

func filterVertical1D(img []float32, result[]float32, shape [2]int, kernel []float32, kernel_anchor int){
	var y_shape int = shape[0]
	var x_shape int = shape[1]

	// offsets say which row/ columns from the start/end where not all kernel values are in range of image
	kernel_offset_L := kernel_anchor
	kernel_offset_R := len(kernel) - 1 - kernel_anchor

	KernelRange := [2]int{-kernel_offset_L, kernel_offset_R + 1}
	Xrange := [2]int{0, x_shape}
	// ranges for Y in each of three loops
	Yranges := [3][2]int{
		{0, kernel_offset_L},
		{kernel_offset_L, y_shape - kernel_offset_R},
		{y_shape - kernel_offset_R, y_shape},
	}

	// first loop - left hands side kernel positions are out of bounds
	for y := Yranges[0][0]; y<Yranges[0][1]; y++{
		for x := Xrange[0]; x<Xrange[1]; x++ {
			var sum float32 = 0.0
			for y_offset := KernelRange[0]; y_offset < KernelRange[1] ; y_offset++ {
				kernel_weight := kernel[kernel_anchor + y_offset]
				img_y := y + y_offset
				img_y_reflected := common.Ternary(img_y < 0, -img_y, img_y)
				img_flat_id := img_y_reflected * x_shape + x
				sum += img[img_flat_id] * kernel_weight
			}
			result[y * x_shape + x] = sum
		}
	}


	// middle loop - all pixels accessible
	for y := Yranges[1][0]; y<Yranges[1][1]; y++{
		for x := Xrange[0]; x<Xrange[1]; x++ {
			var sum float32 = 0.0
			for y_offset := KernelRange[0]; y_offset < KernelRange[1] ; y_offset++ {
				kernel_weight := kernel[kernel_anchor + y_offset]
				img_flat_id := (y + y_offset) * x_shape + x
				sum += img[img_flat_id] * kernel_weight
			}
			result[y * x_shape + x] = sum
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
				img_flat_id := img_y_reflected * x_shape + x
				sum += img[img_flat_id] * kernel_weight
			}
			result[y * x_shape + x] = sum
		}
	}
}

func filterHorizontal1D(img []float32, result[]float32, shape [2]int, kernel []float32, kernel_anchor int){
	var y_shape int = shape[0]
	var x_shape int = shape[1]

	// offsets say which row/ columns from the start/end where not all kernel values are in range of image
	kernel_offset_L := kernel_anchor
	kernel_offset_R := len(kernel) - 1 - kernel_anchor


	KernelRange := [2]int{-kernel_offset_L, kernel_offset_R + 1}
	Yrange := [2]int{0, y_shape}
	// ranges for X in each of three loops
	Xranges := [3][2]int{
		{0, kernel_offset_L},
		{kernel_offset_L, x_shape - kernel_offset_R},
		{x_shape - kernel_offset_R, x_shape},
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
				img_flat_id := y * x_shape + img_x_reflected
				sum += img[img_flat_id] * kernel_weight
			}
			result[y * x_shape + x] = sum
		}
	}


	// middle loop - all pixels accessible
	for y := Yrange[0]; y<Yrange[1]; y++{
		for x := Xranges[1][0]; x<Xranges[1][1]; x++ {
			var sum float32 = 0.0
			for x_offset := KernelRange[0]; x_offset < KernelRange[1] ; x_offset++ {
				kernel_weight := kernel[kernel_anchor + x_offset]
				img_flat_id := y * x_shape + x + x_offset
				sum += img[img_flat_id] * kernel_weight
			}
			result[y * x_shape + x] = sum
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
				img_flat_id := y * x_shape + img_x_reflected
				sum += img[img_flat_id] * kernel_weight
			}
			result[y * x_shape + x] = sum
		}
	}
}

