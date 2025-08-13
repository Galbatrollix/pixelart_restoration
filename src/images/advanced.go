package images

import "image"
import "fmt"

/*
	Given <src> RGBA image, upscale it by expanding each original pixel to <pixel_size> x <pixel_size> square block
	If <grid_size> is more than 0, inserts <grid_color> colored gridlines between each square block, 
	before first-in-row/column square block and after last-in-row/column square block. 

	Resulting image is saved to <dst> RGBA image provided by the caller. 
	Both <src> and <dst> can be subimages (images where rectangle doesn't start at 0,0 and stride doesn't equal 4 * Dx)
	Function will panic if rectangle of <dst> has wrong dimensions.
*/

func AdvancedUpscale(src, dst *image.RGBA, pixel_size, grid_size uint, grid_color [4]uint8){
	advancedUpscaleValidateArguments(src, dst, pixel_size, grid_size)
	RGBAFillColor(dst, grid_color)

	src_base_index := -src.PixOffset(0, 0)
	dst_base_index := -dst.PixOffset(0, 0)

	src_height, src_width := src.Rect.Dy(), src.Rect.Dx()
	pixel_int, grid_int := int(pixel_size), int(grid_size)
	dst_row_len := 4 * (src_width * (pixel_int + grid_int) + grid_int)

	for src_y := 0; src_y < src_height; src_y++ {
		for src_x := 0 ; src_x < src_width; src_x++ {
			// select pixel from source image
			src_pixel_begin := (src_y * src.Stride) + (src_x * 4) + src_base_index
			src_pixel_end := src_pixel_begin + 4
			src_pixel_data := src.Pix[src_pixel_begin: src_pixel_end]

			// find index of topleft corner of corresponding scaled pixel in destination image
			dst_y := grid_int + (grid_int + pixel_int) * src_y
			dst_x := grid_int + (grid_int + pixel_int) * src_x
			dst_pixel_base := (dst_y * dst.Stride) + (dst_x * 4) + dst_base_index

			// copy pixel color from source to destination <pixel_size> times
			for i := 0; i < pixel_int; i++ {
				dst_pixel_begin := dst_pixel_base + i * 4
				dst_pixel_end := dst_pixel_begin + 4
				dst_pixel_data := dst.Pix[dst_pixel_begin: dst_pixel_end]
				copy(dst_pixel_data, src_pixel_data)
			}
		}

		// find index of row filled with pixel colors
		dst_row_y := grid_int + (grid_int + pixel_int) * src_y
		dst_row_x := 0
		dst_row_begin := (dst_row_y * dst.Stride) + (dst_row_x * 4) + dst_base_index
		dst_row_end := dst_row_begin + dst_row_len

		// select entire row as source to copy 
		dst_row_for_copy := dst.Pix[dst_row_begin: dst_row_end]

		// copy entire row <pixel_size - 1> times
		for i := 1 ; i < pixel_int; i++ {
			dst_target_row_begin := dst_row_begin + i * dst.Stride
			dst_target_row_end := dst_target_row_begin + dst_row_len
			dst_target_row := dst.Pix[dst_target_row_begin: dst_target_row_end]
			copy(dst_target_row, dst_row_for_copy)
		}
	}
}

/*
	Does basic validation on input data for function AdvancedUpscale.
	Panics on invalid arguments, does nothing otherwise.

	Does not concern itself with check if dst and src overlap. 
	Does not check if image.RGBA arguments are internally consistent 
	(Doesnt validate against nonsensical structs being passed)

*/
func advancedUpscaleValidateArguments(src, dst *image.RGBA, pixel_size, grid_size uint){
	if src == nil {
		panic("src parameter cannot be nil")
	}
	if dst == nil{
		panic("src parameter cannot be nil")
	}
	if pixel_size == 0 {
		panic("pixel_size parameter must be larger than 0")
	}
	expected_dim := AdvancedUpscaleGetResultDimensions(src.Rect, pixel_size, grid_size)
	var dimensions_ok bool = dst.Rect.Dx() == expected_dim.Dx() && dst.Rect.Dy() == expected_dim.Dy()
	if !dimensions_ok {
		message := fmt.Sprintf(
			"wrong dimensions of destination image\n" +
			"expected [width, height]: [%d, %d], got: [%d, %d]",
			expected_dim.Dx(), expected_dim.Dy(), dst.Rect.Dx(), dst.Rect.Dy(),
		)
		panic(message)
	}

}

/*
	Given <src_rect> as rectangle of source image and <pixel_size>, <grid_size> scaling parameters,
	returns a new rectangle with origin at (0,0) and
	with dimensions required for dst parameter of ImageUpscaleAdvanced function

*/
func AdvancedUpscaleGetResultDimensions(src_rect image.Rectangle, pixel_size, grid_size uint) image.Rectangle {
	src_height, src_width := src_rect.Dy(), src_rect.Dx()
	pixel_int, grid_int := int(pixel_size), int(grid_size)

	dst_height := grid_int + (pixel_int + grid_int) * src_height
	dst_width := grid_int + (pixel_int + grid_int) * src_width

	return image.Rect(0,0, dst_width, dst_height)

}

/*

	Creates and returns an entirely new image upscaled by the rules of "AdvancedUpscale" function.

	Resulting image is not a subimage.
	(result.Rect.Min == 0,0 , result.Stride == result.Rect.Dx() * 4)

*/
func AdvancedUpscaleGetNewImage(src *image.RGBA, pixel_size, grid_size uint, grid_color [4]uint8) *image.RGBA {
	var dst_rect image.Rectangle = AdvancedUpscaleGetResultDimensions(src.Rect, pixel_size, grid_size)
	var dst *image.RGBA = image.NewRGBA(dst_rect)
	AdvancedUpscale(src, dst, pixel_size, grid_size, grid_color)
	return dst
}