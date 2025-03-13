package main


import (
	"fmt"
	"image"
	"time"
	//"image/png"
	//"os"
	//"reflect"
	"io/ioutil"
)

import(
	"pixel_restoration/images"
	"pixel_restoration/images/kuwahara"
	"pixel_restoration/contrast"
	"pixel_restoration/visualizations"
	"pixel_restoration/types"
	"pixel_restoration/gridlines"
)

const DEBUG_DIR_PATH string = "../images/DEBUG"

//https://github.com/mpiannucci/peakdetect/blob/master/peakdetect.go
//https://medium.com/@damithadayananda/image-processing-with-golang-8f20d2d243a2
func automaticGridDetectionMain(input_img *image.RGBA, debug bool) (*image.RGBA, error) {
	/*


		!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

	*/
	// input_img = images.ImageUpscaledWithGridlines(input_img, [4]uint8{0,0,0,255},2,1)


	img_width, img_height := input_img.Rect.Dx(), input_img.Rect.Dy()

	var img_preprocessed *image.RGBA = kuwahara.KuwaharaGaussian(input_img, 2, 1.5)
	// img_preprocessed = input_img

	var edge_distances_rows *image.Gray = contrast.CalculatePixelEdgeDistances(img_preprocessed, false)
	var edge_distances_cols *image.Gray = contrast.CalculatePixelEdgeDistances(img_preprocessed, true)

	var min_peak_height_rows uint8 = contrast.CalculateMinPeakHeight(edge_distances_cols.Pix, contrast.GetBasePeakHeightParams())
	var min_peak_height_cols uint8 = contrast.CalculateMinPeakHeight(edge_distances_rows.Pix, contrast.GetBasePeakHeightParams())

	var edge_rows_binary *image.Gray = contrast.ThresholdWithMinHeight(edge_distances_rows, min_peak_height_rows)
	var edge_cols_binary *image.Gray = contrast.ThresholdWithMinHeight(edge_distances_cols, min_peak_height_cols)


	var edge_rows_binary_cleaned *image.Gray
	edge_rows_binary_cleaned, _ = contrast.CleanupEdgeArtifacts(edge_rows_binary)
	var edge_cols_binary_cleaned *image.Gray
	edge_cols_binary_cleaned, _ = contrast.CleanupEdgeArtifacts(edge_cols_binary)

	/*


		!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

	*/
	//temporary for testing of worst case kind of images
	// edge_rows_binary_cleaned = edge_rows_binary
	// edge_cols_binary_cleaned = edge_cols_binary

	var rows_edge_counts []uint = contrast.EdgesToEdgeCounts(edge_rows_binary_cleaned)
	var cols_edge_counts []uint = contrast.EdgesToEdgeCounts(edge_cols_binary_cleaned)

	var most_frequent_rows []int = contrast.SelectMostFrequent(rows_edge_counts, contrast.GetBaseMostFrequentParams())
	var most_frequent_cols []int = contrast.SelectMostFrequent(cols_edge_counts, contrast.GetBaseMostFrequentParams())

	var rows_intervals types.IntervalList = types.IntervalListFromSortedEdgeIndexes(most_frequent_rows, img_width)
	var cols_intervals types.IntervalList = types.IntervalListFromSortedEdgeIndexes(most_frequent_cols, img_height)

	_,_ = gridlines.GuessGridlineParameters(rows_intervals)
	_,_ = gridlines.GuessGridlineParameters(cols_intervals)

	// fmt.Println("Rows\n", most_frequent_rows)
	// fmt.Println("Cols\n", most_frequent_cols)

	// fmt.Println("Interval rows\n",rows_intervals.Intervals)
	// fmt.Println("Interval cols\n",cols_intervals.Intervals)



	if debug {
		fmt.Println("Image size (width, height): ", img_width, img_height)
		fmt.Println("Min peak height (rows, cols): ", min_peak_height_rows, min_peak_height_cols)


		_ = images.RGBASaveToFile(DEBUG_DIR_PATH + "/1_original.png", input_img)
		_ = images.RGBASaveToFile(DEBUG_DIR_PATH + "/2_kuwaharad.png", img_preprocessed)

		edge_distances_cols_trans := images.GrayscaleGetTransposed(edge_distances_cols)
		edges_sidebyside := visualizations.SideBySideGrayscale(
			edge_distances_rows,
			edge_distances_cols_trans,
		)
		_ = images.GraySaveToFile(DEBUG_DIR_PATH + "/3_edges_sidebyside.png", edges_sidebyside)
		_ = images.GraySaveToFile(DEBUG_DIR_PATH + "/3_edges_rows.png", edge_distances_rows)
		_ = images.GraySaveToFile(DEBUG_DIR_PATH + "/3_edges_cols.png", edge_distances_cols_trans)


		edge_cols_binary_trans := images.GrayscaleGetTransposed(edge_cols_binary)
		edges_binary_sidebyside := visualizations.SideBySideGrayscale(
			edge_rows_binary,
			edge_cols_binary_trans,
		)
		_ = images.GraySaveToFile(DEBUG_DIR_PATH + "/4_edges_binary_sidebyside.png", edges_binary_sidebyside)
		_ = images.GraySaveToFile(DEBUG_DIR_PATH + "/4_edges_binary_rows.png", edge_rows_binary)
		_ = images.GraySaveToFile(DEBUG_DIR_PATH + "/4_edges_binary_cols.png", edge_cols_binary_trans)

		edge_cols_binary_cleaned_trans := images.GrayscaleGetTransposed(edge_cols_binary_cleaned)
		edges_binary_sidebyside_cleaned := visualizations.SideBySideGrayscale(
			edge_rows_binary_cleaned,
			edge_cols_binary_cleaned_trans,
		)
		_ = images.GraySaveToFile(DEBUG_DIR_PATH + "/5_edges_cleaned_sidebyside.png", edges_binary_sidebyside_cleaned)
		_ = images.GraySaveToFile(DEBUG_DIR_PATH + "/6_edges_cleaned_rows.png", edge_rows_binary_cleaned)
		_ = images.GraySaveToFile(DEBUG_DIR_PATH + "/7_edges_cleaned_cols.png", edge_cols_binary_cleaned_trans)

		cutout_image := visualizations.ImageWithDrawnGridlinesSimple(
			input_img,
			[2][]int{ most_frequent_cols,most_frequent_rows},
			[4]uint8{255,0,255,255},
		)

		cutout_image_big := visualizations.ImageWithDrawnGridlinesAdvanced(
			input_img,
			[2][]int{ most_frequent_cols,most_frequent_rows},
			[4]uint8{255,0,255,255},
		)

		_ = images.RGBASaveToFile(DEBUG_DIR_PATH + "/8_cutout_base.png", cutout_image)
		_ = images.RGBASaveToFile(DEBUG_DIR_PATH + "/8_cutout_advanced.png", cutout_image_big)


		// fmt.Println("Edge counts, rows:\n" ,rows_edge_counts)
		// for i, value := range rows_edge_counts {
		// 	fmt.Printf("%d: %d\n", i, value)
		// }

		// fmt.Println("Edge counts, columns:\n" ,cols_edge_counts)
		// for i, value := range cols_edge_counts {
		// 	fmt.Printf("%d: %d\n", i, value)
		// }

	}


	return nil, nil

}



func testThroughDirectory(dirname string){
	items, _ := ioutil.ReadDir(dirname)
    for _, item := range items {
	    filename := item.Name()
        filename_full := dirname + "/" + filename
        img, _ := images.RGBALoadFromFile(filename_full)

        //img = images.ImageUpscaledWithGridlines(img,[4]uint8{0,0,0,255}, 2, 1)
        // img = images.ImageUpscaledByFactor(img, 2)
        fmt.Println(filename_full)
        automaticGridDetectionMain(img, false)
    }

}

func main() {

	img, err := images.RGBALoadFromFile("../images/test_set_pixelarts_clean/CLEAN_34.5_paladin.png")
	if(err != nil){
		fmt.Println(err)
		panic(1)
	}

	start := time.Now()
	const DEBUG = true
	automaticGridDetectionMain(img, DEBUG)
    elapsed := time.Since(start)
    fmt.Println(elapsed)
    
    // test:=images.ImageUpscaledWithGridlines(img, [4]uint8{0,0,0,255}, 1,0 )
    // images.RGBASaveToFile("../images/DEBUG/TEST.png", test)

    testThroughDirectory("../images/test_set_pixelarts_clean/")



}

