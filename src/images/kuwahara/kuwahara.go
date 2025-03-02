package kuwahara

import (
	"image"
)

import (
	"pixel_restoration/images"
)


func KuwaharaGaussian(img *image.RGBA, radius int, sigma float32) *image.RGBA{
	if radius < 1 {
		panic("Radius must be bigger than 0")
	}
	if img.Rect.Dx() * img.Rect.Dy() == 0 {
		panic("Image must have at least one pixel")
	}

	img_shape := [2]int{
		img.Rect.Dy(),
		img.Rect.Dx(),
	}
	total_count := img_shape[0] * img_shape[1]

	// calculates sigma automatically if not provided with valid value
	if sigma <= 0 {
		sigma = 0.3* (float32(radius) - 1.0) + 0.8
	}
	
	// making two semi-kernels for kuwahara filter
	kernel_forward, kernel_reverse := makeSemikernels(radius, sigma)

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

	// get standard deviations 
	greyscale := images.ImageGetGreyscaledChannel(img)
	var standard_deviations [4][]float32 = calculateStandardDeviations(
		greyscale, img_shape, total_count,
		kernel_quadrants, kernel_anchors,
	)


	// calculating color averages
	channels := images.ImageGetSplitChannels(img)
	var color_averages [4][3][]uint8 = getColorAverages(
		channels, img_shape, total_count,
		kernel_quadrants, kernel_anchors,
	)
	
	// choosing indexes of the quadrants with the lowest variance
	var quadrants_chosen []uint8 = chooseQuadrants(standard_deviations)

	// getting output image_data
	var averages_chosen []uint8 = takeAveragesFromChosenQuadrants(color_averages, quadrants_chosen)


	// making values used to assemble final image struct
	new_rect := image.Rect(0,0, img_shape[1], img_shape[0])
	new_stride := img_shape[1] * 4
	new_data := averages_chosen

	return & image.RGBA{
		Pix : new_data,
		Stride: new_stride,
		Rect: new_rect,
	}
}	






func makeSemikernels(radius int, sigma float32) ([]float32, []float32){
	kernel_base := gaussianKernel1D(radius * 2 + 1, sigma)
	kernel_forward := kernel_base[:radius + 1]
	kernel_reverse := kernel_base[radius:]

	// normalizing the semi-kernels so they add up to 1
	var semikernel_sum float32 = sliceSumFloat32(kernel_forward)
	sliceDivbyFloat32(kernel_base, semikernel_sum)

	return kernel_forward, kernel_reverse

}


func calculateStandardDeviations(greyscale []float32, img_shape [2] int , total_count int,
			kernel_quadrants [4][2][]float32, kernel_anchors [4][2]int) [4][]float32 {

	// make space for result
	var deviations [4][]float32
	stddevs_buffer := make([]float32, total_count * 4)
	for i:= 0; i<4; i++{
		deviations[i] = stddevs_buffer[i * total_count : (i+1) * total_count]
	}

	greyscale_squared := sliceGetSquared(greyscale)

	// making required temporary memory
	temporary := make([]float32, total_count)
	greyscale_averages := make([]float32, total_count)

	// calculating standard deviations of each quadrant
	for kernel_id := 0; kernel_id < 4; kernel_id++{
		sepFilter2D(
			greyscale, greyscale_averages, 
			temporary, img_shape,
			kernel_quadrants[kernel_id], kernel_anchors[kernel_id],
		)
        sepFilter2D(
			greyscale_squared, deviations[kernel_id],  
			temporary, img_shape,
			kernel_quadrants[kernel_id], kernel_anchors[kernel_id],
		)
		sliceSubtractSquared(
			deviations[kernel_id], greyscale_averages,
		)
	}
	return deviations
}


func getColorAverages(channels [3][]uint8, img_shape [2]int, total_count int,
				kernel_quadrants [4][2][]float32, kernel_anchors [4][2]int) [4][3][]uint8{

	// making space for result array
	var color_averages [4][3][]uint8
	color_averages_buffer := make([]uint8, 3 * 4 * total_count)
	for kernel_id := 0; kernel_id < 4 ; kernel_id++{
		for channel_id := 0; channel_id < 3; channel_id++{
			start := (kernel_id * 3 + channel_id) * total_count
			end := start + total_count
			color_averages[kernel_id][channel_id] = color_averages_buffer[start: end]
		}
	}

	// allocating temporary stores
	channel_float := make([]float32, total_count)
	channel_averaged := make([]float32, total_count)
	temporary := make([]float32, total_count)

	// calculating color averages
	for channel_id := 0; channel_id < 3; channel_id++{
		sliceUint8ToFloat32(channels[channel_id], channel_float)
		for kernel_id := 0; kernel_id < 4; kernel_id++ {
		  	sepFilter2D(
				channel_float , channel_averaged,
				temporary, img_shape,
				kernel_quadrants[kernel_id], kernel_anchors[kernel_id],
			)
			sliceFloat32ToUint8(channel_averaged, color_averages[kernel_id][channel_id])
		}
	}	

	return color_averages
}

func chooseQuadrants(standard_deviations [4][]float32) []uint8{
	item_count := len(standard_deviations[0])
	quadrants_chosen := make([]uint8, item_count)
	for i:= 0; i<item_count; i++{
		var quadrant_id uint8
		var min_id uint8 = 0
		var min_deviation float32 = standard_deviations[0][i]
		for quadrant_id = 1; quadrant_id < 4; quadrant_id++{
			if standard_deviations[quadrant_id][i] < min_deviation {
				min_deviation = standard_deviations[quadrant_id][i]
				min_id = quadrant_id
			}
		}
		quadrants_chosen[i] = min_id
	}

	return quadrants_chosen
}


func takeAveragesFromChosenQuadrants(color_averages [4][3][]uint8, quadrants_chosen []uint8 ) []uint8 {
	count := len(quadrants_chosen)
	result := make([]uint8, count * 4)

	// choosing color averages according to quadrants chosen
	for flat_id := 0; flat_id < count ; flat_id++ {
		flat_id_result := flat_id * 4
		chosen_quadrant := quadrants_chosen[flat_id]
		result[flat_id_result + 0] = color_averages[chosen_quadrant][0][flat_id]
		result[flat_id_result + 1] = color_averages[chosen_quadrant][1][flat_id]
		result[flat_id_result + 2] = color_averages[chosen_quadrant][2][flat_id]
		result[flat_id_result + 3] = 255        // alpha channel constant
		
	}
	return result
}



