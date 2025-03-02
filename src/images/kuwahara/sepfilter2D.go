package kuwahara


import "pixel_restoration/common"


/*
	sepFilter2D applies a separable linear filter to the single channel image. 
	Function is insipered by similarly named function in OpenCV

	img: 
		slice of row-wise greyscale data of n*m float items 
	dest:
		an output parameter for the result, in the same format as img
		must hold at least n*m values
	temp_buffer:
		Additional slice of memory needed by the function to perform the operation
		Provided as parameter to avoid allocating huge chunks of memory over and over if function is called repeatedly
		must hold at least n*m values
	shape:
		two integer array denoting input (and output) image shape
		shape[0] == n, shape[1] = m
		(where n = number of rows, m = number of columns)
	kernels:
		array of two slices representing column-wise kernel and row-wise kernel respectively
	kerenel_anchors:
		array of two integers representing anchor index in column wise and row-wise kernel respectively
*/


func sepFilter2D(img []float32, dest[]float32, temp_buffer[]float32,
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


/* 
	this is done if X kernel size is larger than X dimension,
	very inefficient, computes advanced reflection logic in each loop iteration
*/
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
				img_y_reflected := reflectIndex101Full(img_y, y_shape-1)
				img_flat_id := img_y_reflected * x_shape + x
				sum += img[img_flat_id] * kernel_weight
			}
			result[y * x_shape + x] = sum
		}
	}

}

/* 
	this is done if X kernel size is larger than X dimension,
	very inefficient, computes advanced reflection logic in each loop iteration
*/
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
				img_x_reflected := reflectIndex101Full(img_x, x_shape-1)
				img_flat_id := y * x_shape + img_x_reflected
				sum += img[img_flat_id] * kernel_weight
			}
			result[y * x_shape + x] = sum
		}
	}
}

/* used by edge case variants of filter, can reflect off walls multiple times */
func reflectIndex101Full(index_check, max_index int) int {
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


/* 
	this is done if Y kernel size is smaller or equal to Y dimension,
	omits checking for reflection where not necessary 
*/
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


/* 
	this is done if X kernel size is smaller or equal to X dimension,
	omits checking for reflection where not necessary 
*/
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
