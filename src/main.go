package main


import (
	"fmt"
	"image"
	"time"
	// "image/png"
	// "os"
	//"reflect"
)

import(
	"pixel_restoration/images"
	"pixel_restoration/images/kuwahara"
	"pixel_restoration/contrast"
)

const DEBUG = true

//https://github.com/mpiannucci/peakdetect/blob/master/peakdetect.go
//https://medium.com/@damithadayananda/image-processing-with-golang-8f20d2d243a2
func automaticGridDetectionMain(input_img *image.RGBA) (*image.RGBA, error) {
	//preprocessing step with kuwahara (and contrast adjustment (? idk if necessary ?))
	img_preprocessed := input_img

	edge_distances_rows := contrast.CalculatePixelEdgeDistances(img_preprocessed, false)
	edge_distances_cols := contrast.CalculatePixelEdgeDistances(img_preprocessed, true)

	min_peak_height_rows := contrast.CalculateMinPeakHeight(edge_distances_cols[0][0: len(edge_distances_cols) * len(edge_distances_cols[0])])
	min_peak_height_cols := contrast.CalculateMinPeakHeight(edge_distances_rows[0][0: len(edge_distances_rows) * len(edge_distances_rows[0])])


	if DEBUG {
		fmt.Println("Image data: ", input_img)

		fmt.Println("Image size (width, height): ", input_img.Rect.Dx(), input_img.Rect.Dy())

		fmt.Println("Edge distances arrays (rows): ", edge_distances_rows,
		 "\nEdge distances arrays (cols):", edge_distances_cols)

		fmt.Println("Min peak height (rows, cols): ", min_peak_height_rows, min_peak_height_cols)

	}


	return nil, nil
}



func main() {

	img, err := images.ImageLoadFromFile("../images/test_set_pixelarts_clean/CLEAN_5.5_gd_mermaid2.png")
	if(err != nil){
		fmt.Println(err)
		panic(1)
	}

	// _ = img
	// automaticGridDetectionMain(img)
	// fmt.Println(images.ImageGetGreyscaledChannel(img))


	kuwaharad:= kuwahara.KuwaharaGaussian(img, 2, 1.5)

	err = images.ImageSaveToFile("../images/test/RESULT.png", kuwaharad)

	start := time.Now()
    img = images.ImageUpscaledByFactor(img, 3)
    elapsed := time.Since(start)
    fmt.Println(elapsed)

  	var row_ids []int = []int{0,5,19,11}
  	var color [4]uint8 = [4]uint8{255,0,0,255}
  	var color2 [4]uint8 = [4]uint8{0,255,0,255}

  	img_smaller := img.SubImage(image.Rect(100, 125, 366, 500)).(*image.RGBA)
  	images.DrawGridlineRowsOnImage(img_smaller, row_ids, color)
  	images.DrawGridlineColsOnImage(img_smaller, row_ids, color2)
	err = images.ImageSaveToFile("../images/test/RESULT.png", img)


}

