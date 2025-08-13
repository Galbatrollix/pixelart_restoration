package visualizations

import "image"

import "pixel_restoration/images"
import "pixel_restoration/types"

/* 
	Creates copy of an image and draws gridlines at selected position indexes with given color
*/
func ImageWithDrawnGridlinesSimple(img *image.RGBA, indexes[2][]int, color [4]uint8) *image.RGBA {
	new_img := images.ImageGetNormalized(img)

	images.DrawGridlineRowsOnImage(new_img, indexes[0], color)
	images.DrawGridlineColsOnImage(new_img, indexes[1], color)

	return new_img
}

/*
	 Creates an upscaled copy of an image and draws black gridlines between each pixel of original image

	 Draws <color> gridlines after selected pixels represented by positions from original image
*/
func ImageWithDrawnGridlinesAdvanced(img *image.RGBA, indexes[2][]int, color [4]uint8) *image.RGBA{
	const pixel_size = 5
	const grid_size = 1
	color_black := [4]uint8{0,0,0,255}

	var img_big *image.RGBA = images.AdvancedUpscaleGetNewImage(img, pixel_size, grid_size, color_black)
	
	// Drawing selected edges in provided color
	indexes_scaled := [2][]int{
		indexesConvertToScaled(indexes[0], pixel_size),
		indexesConvertToScaled(indexes[1], pixel_size),
	}

	images.DrawGridlineRowsOnImage(img_big, indexes_scaled[0], color)
	images.DrawGridlineColsOnImage(img_big, indexes_scaled[1], color)

	return img_big

}

/*
	Creates an upscaled copy of image, visualising contents of combined list data structure.
	Draws black gridlines between each square representing a pixel of original image.
	Blocks of pixels (and edges) that represent unknown items are colored according to <color_unknown>
	Blocks of pixels (and edges) that represent grid items are colored according to <color_gridline> 
	Gridline colors are applied on top of unknown colors. 

*/

func ImageWithDrawnCombinedListAdvanced(
	img *image.RGBA, combined_lists[2] types.CombinedList, color_unknown, color_gridline [4]uint8,
) *image.RGBA {

	// obtain upscaled image with black grid
	const pixel_size = 5
	const grid_size = 1
	color_black := [4]uint8{0,0,0,255}

	var img_big *image.RGBA = images.AdvancedUpscaleGetNewImage(img, pixel_size, grid_size, color_black)

	unknown_ranges := [2][][2]int{
		getIntervalTypePixelRanges(combined_lists[0], types.INTERVAL_UNKNOWN),
		getIntervalTypePixelRanges(combined_lists[1], types.INTERVAL_UNKNOWN),
	}
	unknown_ranges_scaled := [2][][2]int{
		pixelRangesToScaled(unknown_ranges[0], pixel_size),
		pixelRangesToScaled(unknown_ranges[1], pixel_size),
	}
	unknown_indexes := [2][]int{
		scaledRangesToIndexes(unknown_ranges_scaled[0]),
		scaledRangesToIndexes(unknown_ranges_scaled[1]),
	}

	images.DrawGridlineRowsOnImage(img_big, unknown_indexes[0], color_unknown)
	images.DrawGridlineColsOnImage(img_big, unknown_indexes[1], color_unknown)

	grid_ranges := [2][][2]int{
		getIntervalTypePixelRanges(combined_lists[0], types.INTERVAL_GRID),
		getIntervalTypePixelRanges(combined_lists[1], types.INTERVAL_GRID),
	}

	grid_ranges_scaled := [2][][2]int{
		pixelRangesToScaled(grid_ranges[0], pixel_size),
		pixelRangesToScaled(grid_ranges[1], pixel_size),
	}


	grid_indexes := [2][]int{
		scaledRangesToIndexes(grid_ranges_scaled[0]),
		scaledRangesToIndexes(grid_ranges_scaled[1]),
	}
	images.DrawGridlineRowsOnImage(img_big, grid_indexes[0], color_gridline)
	images.DrawGridlineColsOnImage(img_big, grid_indexes[1], color_gridline)

	return img_big
}

/*
	Converts edge-detection indexes to indexes of edges on advanced scaled picture.
*/
func indexesConvertToScaled(indexes []int, pixel_size int) []int{
	result := make([]int, len(indexes))
	for i, value := range indexes{
		result[i] = value * (pixel_size + 1)
	}
	return result
}	

/*
	Generates indexes of all edges between pixels in an advanced scaled picture
*/

func indexesAllScaled(dimension int, scaling_factor int) []int {
	result := make([]int, dimension + 1)
	result[0] = 0
	for i := 1 ; i <= dimension ; i++{
		result[i] = i * scaling_factor
	}
	return result
}

/*
	Makes slice of ranges [int, int] describing all combined list intervals of selected type.
	 where each range describes:
		range[0] is the index of last pixel before the interval
		range[1] is the index of first pixel after the interval
	Therefore, 0-sized intervals have form [n, n+1]
*/
func getIntervalTypePixelRanges(combined_list types.CombinedList, target_interval_type uint8) [][2]int {
	var range_count int = 0
	for _, interval_type := range combined_list.IntervalTypes {
		if interval_type == target_interval_type {
			range_count += 1
		}
	}

	ranges := make([][2]int, 0, range_count)

	var pixel_id int = -1
	var list_len int = len(combined_list.Intervals)

	for i:=0; i<list_len; i++ {
		var interval_type uint8 = combined_list.IntervalTypes[i]
		var interval_length int = int(combined_list.Intervals[i])

		if interval_type == target_interval_type {
			var range_start int = pixel_id
			var range_end int = range_start + interval_length + 1
			ranges = append(ranges, [2]int{range_start, range_end})
		}

		pixel_id += interval_length
	}
	return ranges
}

/*
	Converts ranges of pixels described in "getIntervalTypePixelRanges" function into
	index ranges in an advanced visualization image with pixel_size scaling. 

	Resulting range format: 
		range[0] is the index of first pixel (in scaled image) belonging to the range
		range[1] is the index of last pixel (in scaled image) belonging to the range
*/	
func pixelRangesToScaled(pixel_ranges [][2]int, pixel_size int) [][2]int{
	new_ranges := make([][2]int, len(pixel_ranges))

	for i:=0; i < len(pixel_ranges); i++{
		// how many pixels the range covers
		var pixel_length int = pixel_ranges[i][1] - pixel_ranges[i][0] - 1

		var new_range_start int = (pixel_ranges[i][0] + 1 ) * (pixel_size + 1)
		var new_range_end int = new_range_start + pixel_length * (pixel_size + 1)
		new_ranges[i] = [2]int{new_range_start, new_range_end}
	}
	return new_ranges
}

/*
	Given index ranges obtained from "pixelRangesToScaled" function, 
	return slice of all indexes beloning to ranges in the collection
*/
func scaledRangesToIndexes(scaled_ranges [][2]int) []int{
	var index_count int = 0

	for i := range scaled_ranges {
		start := scaled_ranges[i][0]
		end := scaled_ranges[i][1]
		index_count += end - start + 1
	}

	indexes := make([]int, 0, index_count)

	for _, scaled_range := range scaled_ranges {
		start := scaled_range[0]
		end := scaled_range[1]
		for i := start; i<= end ; i++{
			indexes = append(indexes, i)
		}
	}

	return indexes

}