package images

import (
	"image"
	"math"
	_"fmt"
)

import (
	"pixel_restoration/common"
)


// ( a  a  ab   b  b)
// ( a  a  ab   b  b)
// (ac ac abcd bd bd)
// ( c  c  cd   d  d)
// ( c  c  cd   d  d)

func KuwaharaGaussian(img *image.RGBA, radius int, sigma float32) *image.RGBA{
	if radius < 1 {
		panic("Radius must be bigger than 0")
	}

	img_shape := [2]int{
		img.Rect.Dy(),
		img.Rect.Dx(),
	}
	_ = img_shape
	
	// calculates sigma automatically if not provided with valid value
	if sigma <= 0 {
		sigma = 0.3*(float32(ksize-1)*0.5 - 1.0) + 0.8
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

	greyscale_channel := ImageGetGreyscaledChannel(img)
	greyscale_channel_squared := sliceGetSquared(greyscale_channel)
	greyscale := common.Make2DFromFlat(greyscale_channel, img.Rect.Dy(), img.Rect.Dx())
	greyscale_squared := common.Make2DFromFlat(greyscale_channel_squared, img.Rect.Dy(), img.Rect.Dx())

	// calculating standard deviations of each quadrant
	var standard_deviations[4][][]float32

	for kernel_id := 0; kernel_id < 4 ; kernel_id++{
		greyscale_averages := SepFilter2D(
			greyscale, kernel_quadrants[kernel_id], kernel_anchors[kernel_id],
		)
		standard_deviations[kernel_id] = SepFilter2D(
			greyscale_squared, kernel_quadrants[kernel_id], kernel_anchors[kernel_id],
		)
		sliceSubtractSquared(
			standard_deviations[kernel_id][0][0: img.Rect.Dy() * img.Rect.Dx()],
			greyscale_averages[0][0: img.Rect.Dy() * img.Rect.Dx()],
		)
	}

	// choosing indexes of the quadrants with the lowest variance
	quadrants_chosen := chooseQuadrants(standard_deviations)

	// making space for temporary color averages array
	var color_averages [4][3][]uint8
	var channel_size int = img.Rect.Dy() * img.Rect.Dx()
	color_averages_buffer := make([]uint8, 3 * 4 * channel_size)
	for kernel_id := 0; kernel_id <4 ; kernel_id++{
		for channel_id := 0; channel_id < 3; channel_id++{
			start := (kernel_id * 3 + channel_id) * channel_size
			end := start + channel_size
			color_averages[kernel_id][channel_id] = color_averages_buffer[start: end]
		}
	}

	// calculating color averages
	channels := ImageGetSplitChannels(img)
	channel_float := common.Make2D[float32]( img.Rect.Dy(), img.Rect.Dx())
	channel_float_buffer := channel_float[0][0:img.Rect.Dy()* img.Rect.Dx() ]


	for channel_id := 0; channel_id < 3; channel_id++{
		sliceUint8ToFloat32(channels[channel_id], channel_float_buffer)
		for kernel_id := 0; kernel_id < 4; kernel_id++ {
			averaged := SepFilter2D(
				channel_float, kernel_quadrants[kernel_id], kernel_anchors[kernel_id],
			)
			sliceFloat32ToUint8(averaged[0][0:img.Rect.Dy()* img.Rect.Dx()], color_averages[kernel_id][channel_id])
		}
	}



	// fmt.Println(quadrants_chosen)
	// fmt.Println(color_averages)
	// choosing color averages according to quadrants chosen
	new_rect := image.Rect(0,0,img.Rect.Dx(),img.Rect.Dy())
	new_stride := img.Rect.Dx() * 4
	new_data := make([]uint8, img.Rect.Dy()* img.Rect.Dx() * 4)
	for y := 0; y < img.Rect.Dy() ; y++ {
		for x:= 0;  x<img.Rect.Dx(); x++ {
			flat_id := (y * img.Rect.Dx() + x)
			flat_id_result := flat_id * 4
			chosen_quadrant := quadrants_chosen[y][x]
			new_data[flat_id_result + 0] = color_averages[chosen_quadrant][0][flat_id]
			new_data[flat_id_result + 1] = color_averages[chosen_quadrant][1][flat_id]
			new_data[flat_id_result + 2] = color_averages[chosen_quadrant][2][flat_id]
			new_data[flat_id_result + 3] = 255        // alpha channel constant
		} 
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

func chooseQuadrants(standard_deviations [4][][]float32) [][]uint8{
	y_size := len(standard_deviations[0])
	x_size := len(standard_deviations[0][0])
	quadrants_chosen := common.Make2D[uint8](y_size, x_size)
	for y:= 0; y < y_size ; y++ {
		for x:= 0; x < x_size ; x++ {
			var i uint8
			var min_id uint8 = 0
			var min_deviation float32 = standard_deviations[0][y][x]
			for i = 1; i< 4 ; i++{
				if standard_deviations[i][y][x] < min_deviation {
					min_deviation = standard_deviations[i][y][x]
					min_id = i
				}
			}
			quadrants_chosen[y][x] = min_id
		}
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