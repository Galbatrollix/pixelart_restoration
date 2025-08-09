package visualizations

import "image"

import "pixel_restoration/images"
import "pixel_restoration/types"

/* 
	Creates copy of an image and draws gridlines at all unknown positions and at all positions where interval type is gridline
*/
func ImageWithDrawnCutoutSimple(img *image.RGBA, combined_lists[2] types.CombinedList,
									 color_unknown, color_gridline [4] uint8) *image.RGBA {
	new_img := images.ImageGetNormalized(img)

	unknowns_cols := getIndexesOfUnknowns(combined_lists[0])
	unknowns_rows := getIndexesOfUnknowns(combined_lists[1])

	grids_cols := getIndexesOfGrid(combined_lists[0])
	grids_rows := getIndexesOfGrid(combined_lists[1])

	images.DrawGridlineRowsOnImage(new_img, unknowns_cols, color_unknown)
	images.DrawGridlineColsOnImage(new_img, unknowns_rows, color_unknown)

	images.DrawGridlineRowsOnImage(new_img, grids_cols, color_gridline)
	images.DrawGridlineColsOnImage(new_img, grids_rows, color_gridline)

	return new_img
}


func getIndexesOfUnknowns(combined_list types.CombinedList)[]int{
	current_index := 0

	result := []int{}

	for i := range combined_list.IntervalTypes {
		interval_type := combined_list.IntervalTypes[i]
		interval_size := int(combined_list.Intervals[i])

		for j:=0; j<interval_size;j++{
			if interval_type == types.INTERVAL_UNKNOWN {
				result = append(result, current_index)
			}

			current_index += 1
		}
	}


	return result
}

func getIndexesOfGrid(combined_list types.CombinedList)[]int{

	current_index := 0
	result := []int{}

	for i := range combined_list.IntervalTypes {
		interval_type := combined_list.IntervalTypes[i]
		interval_size := int(combined_list.Intervals[i])


		for j:=0; j<interval_size;j++{
			if interval_type == types.INTERVAL_GRID {
				result = append(result, current_index)
			}

			current_index += 1
		}
	}


	return result

}



/* 
	Creates copy of an image and draws gridlines at all unknown positions and at all positions where interval type is gridline
*/
func ImageWithDrawnCutoutSimpleWithZeros(img *image.RGBA, combined_lists[2] types.CombinedList,
									 	color_unknown, color_gridline [4] uint8) *image.RGBA {
	new_img := images.ImageGetNormalized(img)

	unknowns_cols := getIndexesOfUnknownsWithZeros(combined_lists[0])
	unknowns_rows := getIndexesOfUnknownsWithZeros(combined_lists[1])

	grids_cols := getIndexesOfGridWithZeros(combined_lists[0])
	grids_rows := getIndexesOfGridWithZeros(combined_lists[1])

	images.DrawGridlineRowsOnImage(new_img, unknowns_cols, color_unknown)
	images.DrawGridlineColsOnImage(new_img, unknowns_rows, color_unknown)

	images.DrawGridlineRowsOnImage(new_img, grids_cols, color_gridline)
	images.DrawGridlineColsOnImage(new_img, grids_rows, color_gridline)

	return new_img
}

func getIndexesOfUnknownsWithZeros(combined_list types.CombinedList)[]int{
	current_index := 0

	result := []int{}

	for i := range combined_list.IntervalTypes {
		interval_type := combined_list.IntervalTypes[i]
		interval_size := int(combined_list.Intervals[i])
		/*

			Temporary before better display is made

		*/
		if interval_size == 0 && interval_type == types.INTERVAL_UNKNOWN{
			result = append(result, current_index)
		}

		for j:=0; j<interval_size;j++{
			if interval_type == types.INTERVAL_UNKNOWN {
				result = append(result, current_index)
			}

			current_index += 1
		}
	}


	return result
}

func getIndexesOfGridWithZeros(combined_list types.CombinedList)[]int{

	current_index := 0
	result := []int{}

	for i := range combined_list.IntervalTypes {
		interval_type := combined_list.IntervalTypes[i]
		interval_size := int(combined_list.Intervals[i])

		/*

			Temporary before better display is made

		*/
		if interval_size == 0 && interval_type == types.INTERVAL_GRID{
			result = append(result, current_index)
		}

		for j:=0; j<interval_size;j++{
			if interval_type == types.INTERVAL_GRID {
				result = append(result, current_index)
			}

			current_index += 1
		}
	}


	return result

}