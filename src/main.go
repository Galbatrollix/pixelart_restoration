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

	// img, err := images.ImageLoadFromFile("../images/test/yellow_red_blue.png")
	// if(err != nil){
	// 	fmt.Println(err)
	// 	panic(1)
	// }
	
	// sliceRect := image.Rect(0,0,3,10)
	// img = img.SubImage(sliceRect).(*image.RGBA)

	// fmt.Println(img)
	// img = images.ImageGetNormalized(img)
	// fmt.Println(img)
	// img = images.ImageGetTransposed(img)
	// fmt.Println(img)
	// err = images.ImageSaveToFile("../images/test/DUPA2.png", img)
	// if(err != nil){
	// 	fmt.Println(err)
	// 	panic(1)
	// }

	img, err := images.ImageLoadFromFile("../images/test/gradient_1_python_preprocessed.png")
	if(err != nil){
		fmt.Println(err)
		panic(1)
	}

	_ = img
	// automaticGridDetectionMain(img)
	fmt.Println(images.ImageGetGreyscaled(img))

	fmt.Println("Gauss", images.GaussianKernel1D(9, -1))

	//Gauss [0.014839455 0.049817294 0.11832251 0.19882901 0.23638351 0.19882901 0.11832251 0.049817294 0.014839455]
	//Gauss [0.06277703 0.21074775 0.50055313 0.8411289 1 0.8411289 0.50055313 0.21074775 0.06277703]

	//Gauss [0.011262847 0.036242973 0.087077424 0.15620445 0.20921235 0.20921235 0.15620445 0.087077424 0.036242973 0.011262847]

}


