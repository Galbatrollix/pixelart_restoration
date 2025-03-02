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
	_"pixel_restoration/common"
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

	img, err := images.ImageLoadFromFile("../images/test_set_pixelarts_clean/CLEAN_4_gigantic_difficulty_faces.png")
	if(err != nil){
		fmt.Println(err)
		panic(1)
	}

	// _ = img
	// automaticGridDetectionMain(img)
	// fmt.Println(images.ImageGetGreyscaledChannel(img))

	start := time.Now()
	kuwaharad:= kuwahara.KuwaharaGaussian(img, 2, 1.5)
    elapsed := time.Since(start)
    fmt.Println(elapsed)

	err = images.ImageSaveToFile("../images/test/RESULT.png", kuwaharad)


	// temp_image := common.Make2D[float32](100,100)
	// for y := range temp_image{
	// 	for x := range temp_image[0]{
	// 		temp_image[y][x] = float32(y) * 3.1 + 6.8 * float32(x)
	// 	}
	// }

	// result := images.SepFilter2D(temp_image,[2][]float32{{1,2,3,4},{1,2,5,6}},[2]int{0,2})
	// _ = result

}

