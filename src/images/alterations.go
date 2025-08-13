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
	Fills entirety of a RGBA image with provided color.
	Doesnt consume extra memory, as image uniform is just a color value and bounds

*/
func RGBAFillColor(img *image.RGBA, fill_color [4]uint8){
	color_rgba := color.RGBA{
		R: fill_color[0],
		G: fill_color[1],
		B: fill_color[2],
		A: fill_color[3],
	}
	draw.Draw(img, img.Bounds(), &image.Uniform{color_rgba}, img.Bounds().Min, draw.Src)
}
