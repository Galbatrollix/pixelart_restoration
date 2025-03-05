package main


import (
	"fmt"
	"image"
	"time"
	//"image/png"
	//"os"
	//"reflect"
)

import(
	"pixel_restoration/images"
	"pixel_restoration/images/kuwahara"
	"pixel_restoration/contrast"
	"pixel_restoration/visualizations"

	//"pixel_restoration/types"
)

const DEBUG = true
const DEBUG_DIR_PATH string = "../images/DEBUG"

//https://github.com/mpiannucci/peakdetect/blob/master/peakdetect.go
//https://medium.com/@damithadayananda/image-processing-with-golang-8f20d2d243a2
func automaticGridDetectionMain(input_img *image.RGBA) (*image.RGBA, error) {
	//consider adding contrast adjustment too but idk if its necessary
	var img_preprocessed *image.RGBA = kuwahara.KuwaharaGaussian(input_img, 2, 1.5)

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


	var rows_edge_counts []uint = contrast.EdgesToEdgeCounts(edge_rows_binary_cleaned)
	var cols_edge_counts []uint = contrast.EdgesToEdgeCounts(edge_cols_binary_cleaned)
	// var rows_interval_list types.IntervalList = ...
	// var cols_interval_list types.IntervalList = ...





	if DEBUG {
		fmt.Println("Image size (width, height): ", input_img.Rect.Dx(), input_img.Rect.Dy())
		fmt.Println("Min peak height (rows, cols): ", min_peak_height_rows, min_peak_height_cols)

		edge_distances_cols_trans := images.GrayscaleGetTransposed(edge_distances_cols)
		edges_sidebyside := visualizations.SideBySideGrayscale(
			edge_distances_rows,
			edge_distances_cols_trans,
		)
		_ = images.GraySaveToFile(DEBUG_DIR_PATH + "/edges_sidebyside.png", edges_sidebyside)
		_ = images.GraySaveToFile(DEBUG_DIR_PATH + "/edges_rows.png", edge_distances_rows)
		_ = images.GraySaveToFile(DEBUG_DIR_PATH + "/edges_cols.png", edge_distances_cols_trans)


		edge_cols_binary_trans := images.GrayscaleGetTransposed(edge_cols_binary)
		edges_binary_sidebyside := visualizations.SideBySideGrayscale(
			edge_rows_binary,
			edge_cols_binary_trans,
		)
		_ = images.GraySaveToFile(DEBUG_DIR_PATH + "/edges_binary_sidebyside.png", edges_binary_sidebyside)
		_ = images.GraySaveToFile(DEBUG_DIR_PATH + "/edges_binary_rows.png", edge_rows_binary)
		_ = images.GraySaveToFile(DEBUG_DIR_PATH + "/edges_binary_cols.png", edge_cols_binary_trans)

		edge_cols_binary_cleaned_trans := images.GrayscaleGetTransposed(edge_cols_binary_cleaned)
		edges_binary_sidebyside_cleaned := visualizations.SideBySideGrayscale(
			edge_rows_binary_cleaned,
			edge_cols_binary_cleaned_trans,
		)
		_ = images.GraySaveToFile(DEBUG_DIR_PATH + "/edges_cleaned_sidebyside.png", edges_binary_sidebyside_cleaned)
		_ = images.GraySaveToFile(DEBUG_DIR_PATH + "/edges_cleaned_rows.png", edge_rows_binary_cleaned)
		_ = images.GraySaveToFile(DEBUG_DIR_PATH + "/edges_cleaned_cols.png", edge_cols_binary_cleaned_trans)


		fmt.Println("Edge counts, rows:\n" ,rows_edge_counts)
		for i, value := range rows_edge_counts {
			fmt.Printf("%d: %d\n", i, value)
		}

		fmt.Println("Edge counts, columns:\n" ,cols_edge_counts)
		for i, value := range cols_edge_counts {
			fmt.Printf("%d: %d\n", i, value)
		}

	}


	return nil, nil
}



func main() {

	img, err := images.RGBALoadFromFile("../images/test_set_pixelarts_clean/CLEAN_34.5_paladin.png")
	if(err != nil){
		fmt.Println(err)
		panic(1)
	}

	_ = img
	automaticGridDetectionMain(img)

	kuwaharad:= kuwahara.KuwaharaGaussian(img, 2, 1.5)

	err = images.RGBASaveToFile("../images/test/RESULT.png", kuwaharad)

	start := time.Now()
    img = images.ImageUpscaledByFactor(img, 3)
    elapsed := time.Since(start)
    fmt.Println(elapsed)



}

