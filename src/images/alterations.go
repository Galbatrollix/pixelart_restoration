/*
	Alterations file contains functions that alter the image files in the following ways:
	- upscale (also upscale with gridlines)
	- draw gridlines
	- add noise to image <not yet added>
	- add watermakrs to image <not yet added>
*/

package images

import (
	"image"
	"image/draw"
	"image/color"
	"fmt"
)


/*
	ImageUpscaledByFactor returns new RGBA image where each pixel in original cooresponds to
	factor times factor square of the same color in the output

	The solution is rather overengineered, but i wondered
	 if it could be done this way instead of iterating over each destination pixel,
	 getting modulo to map it to original pixel and copying the value.

	and it certainly could have been done
*/

func ImageUpscaledByFactor(img *image.RGBA, factor int) *image.RGBA{
	if factor < 0 {
		panic("Upscale by negative factor attempted")
	}
	if factor > 100{
		panic("Upscale by factor higher than 100 attempted.")
	}
	// simplifies iteration logic if image rect starts at 0,0 and stride is equal to 4 * width
	img = ImageGetNormalized(img)

	height, width := img.Rect.Dy(), img.Rect.Dx()
	upscaled := image.NewRGBA(image.Rect(0, 0, width * factor, height * factor))

	for y:=0; y<height; y++ {
		// make new row
		for x:=0 ;x<width; x++ {
			source_pixel_begin := (y * width + x) * 4
			source_pixel_end := source_pixel_begin + 4
			// selecting pixel from original image
			source_pixel_data := img.Pix[source_pixel_begin: source_pixel_end]

			dest_y := factor * factor * y
			dest_x := factor * x
			dest_pixel_begin_base := (dest_y * width + dest_x) * 4
			// copy selected pixel (factor) times
			for i:=0; i<factor;i++{
				dest_pixel_begin := dest_pixel_begin_base + i * 4
				dest_pixel_end := dest_pixel_begin + 4;
				dest_pixel_data := upscaled.Pix[dest_pixel_begin: dest_pixel_end]
				copy(dest_pixel_data, source_pixel_data)
			}

		}
		// copy new row (factor - 1) times
		source_begin := factor * factor * y * 4 * width
		source_end := source_begin + upscaled.Stride
		source := upscaled.Pix[source_begin:source_end]
		for i:=1; i<factor; i++ {
			dest_begin := source_begin + i * upscaled.Stride
			dest_end := dest_begin +  upscaled.Stride
			dest := upscaled.Pix[dest_begin:dest_end]
			copy(dest, source)
		}
	}
	return upscaled

}


	
/*

	DrawGridlineRowsOnImage fills selected rows of input image with provided RGBA color.
	Selected rows are chosen based on y_indexes slice. 
	If y_indexes holds invalid row index, function will panic.

*/
func DrawGridlineRowsOnImage(img *image.RGBA, y_indexes []int, color [4]uint8) {
	height, width := img.Rect.Dy(), img.Rect.Dx()
	//validate y_indexes
	for _, y_id := range y_indexes{
		if y_id < 0 {
			panic(fmt.Sprintf("Provided y_index %d is less than 0", y_id))
		}
		if y_id > height - 1 {
			panic(fmt.Sprintf("Provided y_index %d is larger than image height - 1 (%d)", y_id, height - 1))
		}
	}

	// for each provided row, fill with given color
	for _, y := range y_indexes{
		flat_id_base := img.PixOffset(0 + img.Rect.Min.X, y + img.Rect.Min.Y)
		for i:=0; i<width; i++{
			destination := img.Pix[flat_id_base + i*4:]
			copy(destination, color[:])
		}
	}
}

/*

	DrawGridlineColsOnImage fills selected columns of input image with provided RGBA color.
	Selected columns are chosen based on x_indexes slice. 
	If x_indexes holds invalid row index, function will panic.

*/
func DrawGridlineColsOnImage(img *image.RGBA, x_indexes []int, color [4]uint8){
	height, width := img.Rect.Dy(), img.Rect.Dx()
	//validate y_indexes
	for _, x_id := range x_indexes{
		if x_id < 0 {
			panic(fmt.Sprintf("Provided y_index %d is less than 0", x_id))
		}
		if x_id > width - 1 {
			panic(fmt.Sprintf("Provided y_index %d is larger than image width - 1 (%d)", x_id, width - 1))
		}
	}

	// for each provided column, fill with given color
	for _, x := range x_indexes{
		flat_id_base := img.PixOffset(x + img.Rect.Min.X, 0 + img.Rect.Min.Y)
		for i:=0; i<height; i++{
			destination := img.Pix[flat_id_base + i*img.Stride:]
			copy(destination, color[:])
		}
	}
}


/*
	ImageUpscaledWithGridlines returns new RGBA image where each pixel in original 
	cooresponds to a pixel-width sized square. Between each square there is a
	gridline_width-sized gridline drawn with provided color.

*/
func ImageUpscaledWithGridlines(img *image.RGBA, color [4]uint8, pixel_width int, gridline_width int) *image.RGBA {
	if pixel_width < 0 {
		panic("Pixel width negative")
	}
	if gridline_width < 0{
		panic("Gridline width negative")
	}
	if gridline_width + pixel_width < 1 {
		panic("Total upscale factor less than one")
	}


	height, width := img.Rect.Dy(), img.Rect.Dx()

	scaling_factor := pixel_width + gridline_width
	upscaled := ImageUpscaledByFactor(img, scaling_factor)

	x_ids := gridlineIdsForUpscaled(width, scaling_factor, gridline_width)
	y_ids := gridlineIdsForUpscaled(height, scaling_factor, gridline_width)


	DrawGridlineRowsOnImage(upscaled, y_ids, color)
	DrawGridlineColsOnImage(upscaled, x_ids, color)

	return upscaled
}

func gridlineIdsForUpscaled(dim_size int, scaling_factor int, gridline_width int) []int {
	ids := make([]int, 0, gridline_width * dim_size)
	for n:=0 ;n<dim_size; n++{
		base := (n + 1) * scaling_factor - 1
		for i:=0; i<gridline_width; i++{
			ids = append(ids, base - i)
		}
	}

	return ids
}


/*
	Fills entirety of a RGBA image with provided color.
	Temporarily works using builtin not-in-place method.
	TODO: Rework to handrolled in-place method.

*/
func RGBAFillColor2(img *image.RGBA, fill_color [4]uint8){
	color_rgba := color.RGBA{
		R: fill_color[0],
		G: fill_color[1],
		B: fill_color[2],
		A: fill_color[3],
	}
	draw.Draw(img, img.Bounds(), &image.Uniform{color_rgba}, img.Bounds().Min, draw.Src)
}

/*
	Fills entirety of a RGBA image with provided color.
	Faster than using draw package from standard library
*/
func RGBAFillColor(img *image.RGBA, fill_color [4]uint8){
	color_as_slice := fill_color[:]
	base_index := -img.PixOffset(0, 0)
	// fill fist row by copying values from color one by one
	for x := 0; x < img.Rect.Dx(); x++ {
		dest_pixel := img.Pix[base_index + x*4: base_index + (x + 1)* 4]
		copy(dest_pixel, color_as_slice)
	}

	// fill remaining rows (if any) by copying first row into them
	row_length_bytes := img.Rect.Dx() * 4
	first_row := img.Pix[base_index: base_index + row_length_bytes]
	for y := 1; y < img.Rect.Dy(); y++{
		row_begin := base_index + y * img.Stride
		row := img.Pix[row_begin: row_begin + row_length_bytes]
		copy(row, first_row)
	}
}