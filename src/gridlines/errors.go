package gridlines 

import (
	"fmt"
	"math"
	"slices"
)
import "pixel_restoration/types"

/*

	This function takes combined list with unknowns and fills the unknown gaps based on correct sections.
	Returns new combined list with no unknown sections.

*/
func GridlinesFixErrors(original_combined_list types.CombinedList, pixel_guess, grid_guess types.IntervalRangeEntry) types.CombinedList {	
	var left_edge_unknown, right_edge_unknown uint
	var middle_unknowns []uint

	left_edge_unknown, right_edge_unknown, middle_unknowns = separateUnknownItems(original_combined_list)
	var mean_pixel, mean_grid float64 = calculateItemAverages(original_combined_list, pixel_guess, grid_guess, [][]uint{})

	// making continous space for all fixed unknown sections, including edges
	fixed_sections := make([][]uint, len(middle_unknowns) + 2)

	middle_fixed := fixed_sections[1:len(middle_unknowns) + 1]
	for i := range middle_unknowns {
		middle_fixed[i] = guessMiddleUnknownSection(middle_unknowns[i], mean_pixel, mean_grid)
	}

	// recalculating averages with fixed sections to improve acuraccy for edge guessing
	mean_pixel, mean_grid = calculateItemAverages(original_combined_list, pixel_guess, grid_guess, middle_fixed)

	left_edge_fixed  := guessEdgeUnknownSection(left_edge_unknown, mean_pixel, mean_grid, true)
	right_edge_fixed := guessEdgeUnknownSection(right_edge_unknown, mean_pixel, mean_grid, false)
	fixed_sections[0] = left_edge_fixed
	fixed_sections[len(middle_unknowns) + 1] = right_edge_fixed

	var fixed_combined_list types.CombinedList = reAssembleCombinedList(original_combined_list, fixed_sections)
	fmt.Println("INTERVALS: ", fixed_combined_list.Intervals)
	return fixed_combined_list

}


/*	
	Given a combined list with unknowns, return 3 values:
	1) length of left edge unknown item
	2) length of right edge unknown item
	3) slice of lengths of non-edge unknown items, in the same order as the ordering in the input list
*/
func separateUnknownItems(combined_list types.CombinedList)(uint, uint, []uint){
	length := len(combined_list.Intervals)
	left, right := combined_list.Intervals[0], combined_list.Intervals[length - 1]

	middle := make([]uint, 0, 32)

	for i:=1 ; i<length - 1; i++{
		if combined_list.IntervalTypes[i] == types.INTERVAL_UNKNOWN{
			middle = append(middle, combined_list.Intervals[i])
		}
	}

	return left, right, middle

}

/*
	Returns average size of pixel and average size of gridline in the correct sections of combined list.
	If non-empty array of fixed middle sections is supplied, treats elements in each section as starting with grid 
		and includes them in average calculation.
	If no items of certain type found, then sets average to midpoint of pixel guess or grid guess respectively.

	Excludes edges.

*/
func calculateItemAverages(
	combined_list types.CombinedList,
	pixel_guess, grid_guess types.IntervalRangeEntry,
	fixed_middle_sections [][]uint,
)(float64, float64) {
	var sum_pixel, sum_grid float64 = 0.0, 0.0
	var count_pixel, count_grid float64 = 0.0, 0.0

	// going over combined list and excluding edge elements
	length := len(combined_list.Intervals)
	for i:=1 ; i<length - 1; i++{
		if combined_list.IntervalTypes[i] == types.INTERVAL_GRID{
			sum_grid += float64(combined_list.Intervals[i])
			count_grid += 1
		}

		if combined_list.IntervalTypes[i] == types.INTERVAL_PIXEL{
			sum_pixel += float64(combined_list.Intervals[i])
			count_pixel += 1
		}

	}

	// going over fixed middle sections and assuming first value of each section to be grid
	for _, section := range fixed_middle_sections {
		for i, length := range section {
			if i % 2 == 0 {
				sum_grid += float64(length)
				count_grid += 1
			}else{
				sum_pixel += float64(length)
				count_pixel += 1
			}

		}
	}

	var pixel_average, grid_average float64

	// If no gridline found 
	if count_grid == 0.0 {
		grid_average = float64(grid_guess.Bounds[0] + grid_guess.Bounds[1]) / 2.0
	}else{
		grid_average = sum_grid / count_grid
	}


	// If no pixels found (very unlikely to happen but whatever)
	if count_pixel == 0.0 {
		pixel_average = float64(pixel_guess.Bounds[0] + pixel_guess.Bounds[1]) / 2.0
	}else{
		pixel_average = sum_pixel / count_pixel
	}

	return pixel_average, grid_average
}

/*

	Given length of unknown item surrounded by correct pixel items, 
	and given mean pixel and mean grid values for correct sections of the image:
	Approximates grid/pixel sequence of the unknown section. 
	Returns sequence as uint slice, first item is assumed to be grid type. 
	(item sequence alternates pixel and grid elements)

*/
func guessMiddleUnknownSection(unknown_length uint, mean_pixel, mean_grid float64) []uint {
	// this funciton assumes that unknown interval has gridlines at both ends
	// n shall be a mathematically expected number of pixels in a sequence
    n := (float64(unknown_length) - mean_grid) / (mean_grid + mean_pixel)

    guessed_pixel_count := int(math.Round(n))
    guessed_grid_count := guessed_pixel_count + 1

    if guessed_pixel_count <= 0 {
    	return []uint{unknown_length}
    }

    // creating grid sections
	grid_base_size := uint(math.Round(mean_grid))
    var grid_sections []uint = makeGridBaseSections(guessed_grid_count, grid_base_size)

    // creating pixel sections
    remaining_length := unknown_length - uint(guessed_grid_count) * grid_base_size
    var pixel_sections []uint = distributeEvenly(guessed_pixel_count, remaining_length)

    // combining grid and pixel sections into a single slice 
   	result := slicesInterleave(grid_sections, pixel_sections)
   	return result

}

/*

	Given length of unknown item on either left or right edge of the picture, 
	and given mean pixel and mean grid values for correct sections of the image:
	Approximates grid/pixel sequence of the edge-located unknown section. 
	Returns sequence as uint slice.
	
	If left edge: 
		Last item is of grid type. 
		First item may be shorter due to snipping, but cannot be 0-length
	If right edge:
		first item is of grid type. 
		Last item may be shorter due to snipping, but cannot be 0-length
	(item sequence alternates pixel and grid elements)

*/

func guessEdgeUnknownSection(unknown_length uint, mean_pixel, mean_grid float64, is_left_edge bool) []uint {
	// n is estimated count of pixels that unknown_length can contain
    n := int(float64(unknown_length)  / (mean_grid + mean_pixel))

    // fabricate a sufficiently large "unknown length" for guessMiddleUnknownSection function to solve
    // chosen "unknown length" should allow solving with as little deviation from averages as possible
    dummy_sequence_pixel_count := n + 10 
    dummy_sequence_length := uint((mean_grid + mean_pixel) * float64(dummy_sequence_pixel_count) + mean_grid + 0.5)

    // solve the prepared dummy sequence 
    var dummy_sequence_solved []uint = guessMiddleUnknownSection(dummy_sequence_length, mean_pixel, mean_grid)

    // trim the solved dummy sequence to desired length 
    trimSequenceFromRight(&dummy_sequence_solved, unknown_length)
 	var trimmed []uint = dummy_sequence_solved

    if is_left_edge {
    	slices.Reverse(trimmed)    
    }

    return trimmed
}

/*
	Returns slice of N grid_base items
*/

func makeGridBaseSections(grid_count int, grid_base uint) []uint {
    grid_sections := make([]uint, grid_count)
    for i := range grid_sections{
    	grid_sections[i] = grid_base
    }
    return grid_sections
}

/*
	Interleaves left and right slice, starting and ending with left. 
	This means that len(left) is assumed to be equal to len(right) + 1
*/
func slicesInterleave(left, right []uint) []uint {
	interleaved := make([]uint, len(left) + len(right))
    
	slices := [2][]uint{left, right}

	for i := 0; i< len(left) + len(right);i++ {
		slice_index := i / 2
		which_slice := i % 2
		interleaved[i] = slices[which_slice][slice_index]
	}
	return interleaved

}

/*
	Given a number of items items, distribute items evenly into buckets.

	Return a buckets-length uint slice with values representing number of items in each bucket.

	Example:
	    num_items: 33
	    num_buckets: 10
	Result:
	    [3, 4, 3, 3, 4, 3, 3, 4, 3, 3]
*/
func distributeEvenly(num_buckets int, num_items uint) []uint{
	if num_buckets == 0 {
		return []uint{}
	}else if num_buckets == 1 {
		return []uint{num_items}
	}

	base := num_items / uint(num_buckets) 
	remainder := num_items % uint(num_buckets) 

	// making result slice and filling it with base item size
	result := make([]uint, num_buckets)
	for i :=  range result {
		result[i] = base - 1
	}

    // incrementing result on each positions pointed by bresenham line indexes
	bresenham_indexes := bresenhamLine(0, 0, int(remainder) + num_buckets - 1, num_buckets - 1)	
	for _, i := range bresenham_indexes {
		result[i] += 1
	}

	// performing additional rotations so the distrubution is even on edge-cases too
	smaller_count := sliceCountMin(result)

	_ = smaller_count

	return result
}

/*
 	Bresenham line algorithm, used for evenly distributing items into buckets.
*/
func bresenhamLine(x1, y1, x2, y2 int) []int{
	dx := int(math.Abs(float64(x2 - x1)))
	dy := int(math.Abs(float64(y2 - y1)))

	gradient := float64(dy) / float64(dx)

	if gradient > 1{
        dx, dy = dy, dx
        x1, y1 = y1, x1
        x2, y2 = y2, x2
	}

	p := 2 * dy - dx

	// x_coordinates := make([]int, 0, dx + 1)
	y_coordinates := make([]int, 0, dx + 1)

	// x_coordinates = append(x_coordinates, x1)
	y_coordinates = append(y_coordinates, y1)

	for k := 2; k < dx + 2; k++{
		if p > 0 {
			if y1 < y2 {
				y1 = y1 + 1
			}else{
				y1 = y1 - 1
			}
			p = p + 2 * (dy - dx)
		}else{
			p = p + 2 * dy
		}

		if x1 < x2 {
			x1 = x1 + 1
		}else{
			x1 = x1 - 1
		}

		// x_coordinates = append(x_coordinates, x1)
		y_coordinates = append(y_coordinates, y1)

	}
	return y_coordinates
}

/*
	Returns count of items equal to min(slice)
*/
func sliceCountMin(slice []uint) int {
	min_val := slices.Min(slice)

	count := 0
	for _, val := range slice {
		if min_val == val {
			count += 1
		}
	}

	return count
}

/*
	Trims interval from the right side, leaving total length of interval sequence equal to target length.
	If required, trimming may reduce the length of the last remaining element.
	Lengths of other remaining elements are guaranteed to not change.
	Sequence after trim is guaranteed to not end with 0-length interval. 

	Interval sequence is modified in place by a pointer
	Function will panic if target length is larger than sum of interval lengths in the sequence

	Example input:
		*interval_sequence:     [0,6,1,6,0,5,1]
		target_item_length:     11
	Example result:
		*interval_sequence:	    [0,6,1,4]

*/
func trimSequenceFromRight(interval_sequence *[]uint, target_item_length uint){
	// find last element that will remain
	var current_index int = 0
	var accumulated_length uint = (*interval_sequence)[0]
	for accumulated_length < target_item_length {
		current_index += 1
		accumulated_length += (*interval_sequence)[current_index]
	}

	// reslice the sequence to the last element
 	*interval_sequence = (*interval_sequence)[:current_index + 1]

 	// modify the length of last element to fit the target 
	var difference uint = accumulated_length - target_item_length
	(*interval_sequence)[current_index] -= difference
}

/*
	Given original combined list with unknown sections and slices describing 
	guessed values for each unknown section, generate a new "fixed" combined list struct.

	Resulting combined list doesn't contain any "unknown" items 
	and must contain alternating "grid" and "pixel" type items only

*/

func reAssembleCombinedList(original_list types.CombinedList, fixed_sections [][]uint) types.CombinedList {
	var result_length int = reAssembleGetTotalLength(original_list, fixed_sections)
	var intervals []uint = reAssembleCreateIntervals(result_length, original_list, fixed_sections)

	var starts_with_grid bool = len(fixed_sections[0]) % 2 == 1
 	var interval_types []uint8 = reAssembleCreateTypes(result_length, starts_with_grid)

 	return types.CombinedList{intervals, interval_types} 
}

/*
	Part of reAssembleCombinedList function. 
	Based on original combined list and fixed sections, calculate exact number of interval items 
	that the resultant combined list will hold.
*/
func reAssembleGetTotalLength(original_list types.CombinedList, fixed_sections [][]uint) int {
	original_count_total := len(original_list.Intervals)
	unknowns_count := len(fixed_sections)

	fixed_sections_element_count := 0
	for _, section := range fixed_sections {
		fixed_sections_element_count += len(section)
	}

	return original_count_total - unknowns_count + fixed_sections_element_count
}

/*
	Part of reAssembleCombinedList function. 
	Based on original combined list and fixed sections,
	construct a fully complete intervals slice for new combined list instance.

	Exact count of elements in result array is provided as first parameter.
*/
func reAssembleCreateIntervals(
	result_length int,
	original_list types.CombinedList,
	fixed_sections [][]uint,
) []uint {
	new_intervals := make([]uint, result_length)
	var intervals_id int = 0
	var current_section int = 0

	var original_list_len int = len(original_list.Intervals)
	for i := 0; i < original_list_len; i++ {
		var is_unknown bool = original_list.IntervalTypes[i] == types.INTERVAL_UNKNOWN
		if is_unknown {
			for _, item_length := range fixed_sections[current_section]{
				new_intervals[intervals_id] = item_length
				intervals_id += 1
			}
			current_section += 1
		}else{ // non-unknown item
			var item_length uint = original_list.Intervals[i]
			new_intervals[intervals_id] = item_length
			intervals_id += 1
		}

	}

	return new_intervals
}

/*
	Part of reAssembleCombinedList function. 
	Constructs interval types slice for combined list 
	based on total list length and type of first element.
*/
func reAssembleCreateTypes(result_length int, starts_with_grid bool) []uint8{
	new_interval_types := make([]uint8, result_length)
	var even_lookup [2]uint8
	if starts_with_grid {
		even_lookup = [2]uint8{types.INTERVAL_GRID, types.INTERVAL_PIXEL}
	}else{
		even_lookup = [2]uint8{types.INTERVAL_PIXEL, types.INTERVAL_GRID}
	}

	for i := range new_interval_types {
		new_interval_types[i] = even_lookup[i % 2]
	}
	return new_interval_types
}