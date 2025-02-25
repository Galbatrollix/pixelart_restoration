package main


import (
	"fmt"
	"image"
	// "image/png"
	// "os"
	//"reflect"
)

import(
	"pixel_restoration/images"
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


	if DEBUG {
		fmt.Println("Image size (width, height): ", input_img.Rect.Dx(), input_img.Rect.Dy())

		fmt.Println("Edge distances arrays (rows): ", edge_distances_rows,
		 "\nEdge distances arrays (cols):", edge_distances_cols)

	}


	return nil, nil
}

func main() {

	img, err := images.ImageLoadFromFile("../images/test/yellow_red_blue.png")
	if(err != nil){
		fmt.Println(err)
		panic(1)
	}
	
	sliceRect := image.Rect(0,0,3,10)
	img = img.SubImage(sliceRect).(*image.RGBA)

	fmt.Println(img)
	img = images.ImageGetNormalized(img)
	fmt.Println(img)
	img = images.ImageGetTransposed(img)
	fmt.Println(img)
	err = images.ImageSaveToFile("../images/test/DUPA2.png", img)
	if(err != nil){
		fmt.Println(err)
		panic(1)
	}

	// img, err := images.ImageLoadFromFile("../images/test/gradient_1.png")
	// if(err != nil){
	// 	fmt.Println(err)
	// 	panic(1)
	// }
	// fmt.Println(img)

	automaticGridDetectionMain(img)




}