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
	"pixel_restoration/common"
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

	temp_image := common.Make2D[float32](1000,1000)
	for y := range temp_image{
		for x := range temp_image[0]{
			temp_image[y][x] = float32(y) * 3.1 + 6.8 * float32(x)
		}
	}

	result := images.SepFilter2D(temp_image,[2][]float32{{1,2,3,4},{1,2,5,6}},[2]int{0,2})
	_ = result

}



/*


[[ 1548.      2092.      2908.      3860.      4812.      5764.
   6716.      7668.      7804.0005]
 [ 1982.      2526.      3342.      4294.      5246.      6198.
   7150.      8102.      8238.    ]
 [ 2416.      2960.      3776.      4728.      5680.      6632.
   7584.      8536.      8672.    ]
 [ 2850.      3394.      4210.      5162.      6114.      7066.
   8018.      8970.      9106.    ]
 [ 3284.      3828.      4644.      5596.      6548.      7500.
   8452.      9404.      9540.    ]
 [ 3718.      4262.      5078.      6030.      6982.      7934.
   8886.      9838.      9974.    ]
 [ 4152.      4696.      5512.      6464.      7416.      8368.
   9320.     10272.     10408.    ]
 [ 4238.8     4782.8     5598.8     6550.8     7502.8     8454.8
   9406.8    10358.8    10494.8   ]
 [ 4065.2     4609.1997  5425.2     6377.2     7329.2     8281.199
   9233.199  10185.2    10321.2   ]
 [ 3718.      4262.      5078.      6030.      6982.      7934.
   8886.      9838.      9974.    ]]


[[1548 2092 2908 3860 4812 5764.0005 6716 7668 7804.0005] 
[1982 2526 3342.0002 4294 5246 6198 7150.0005 8102 8238] 
[2416 2960 3776 4728 5680 6632 7584.0005 8536 8672] 
[2850 3394 4210 5162 6114 7066 8018 8970 9106] 
[3283.9998 3828 4644 5596 6548 7500 8452 9404 9540] 
[3717.9998 4262 5078 6030 6982 7934.0005 8886 9838 9974] 
[4152 4696 5512 6464 7416 8368 9320 10272 10408] 
[4238.8 4782.8 5598.8 6550.8 7502.8 8454.801 9406.801 10358.8 10494.8] 
[4065.2 4609.1997 5425.2 6377.2 7329.1997 8281.201 9233.2 10185.2 10321.2]
 [3717.9998 4261.9995 5078 6030 6982 7934 8886 9838 9974]]

*/
