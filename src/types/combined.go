package types

/*

	Combined list holds two attributes: Intevals and IntervalTypes
	Intervals:
		slice of N intervals of consecutive pixels without edges between them.
		Intervals slice can contain 0 - length items!
	IntervalTypes:
		slice of N uint8 items denoting type of interval item of the same index - Pixel, Gridline or Unknown
		definition of values of this slice can be seen below in enum-like constant group

	For example:
		Intervals: 		[4,1,5,13,4,1,5]
		IntervalTypes: 	[0,1,0,2 ,0,1,0]

		This means that there are 33 pixels in the image, split into 7 consecutive segments, which represents the following:
		4 pixels 'pixel' type -> 1 pixel 'gridline' type -> 5 pixels 'pixel' type, -> 13 pixels 'unknown' type (...)
*/
type CombinedList struct{
	Intervals []uint
	IntervalTypes []uint8

}

/*
	Defines a set of 3 uint8 constants used to distinguish between combined list item types
*/
const (  
		INTERVAL_UNKNOWN uint8 = iota
        INTERVAL_PIXEL uint8 = iota   
       	INTERVAL_GRID uint8 = iota    

)


// func (*CombinedList) SquashItemTypes()

func CombinedFromIntervalList(intervals IntervalList, guessed_params [2]IntervalRangeEntry) CombinedList {
	if len(intervals.Intervals) < 3 {
		return CombinedList{[]uint{}, []uint8{}} 
		//panic("This should never happen, guard against double zero antry and shorter than 3 items interval list")
	}

	var case_grid_1_2 bool = guessed_params[1].Bounds[0] == 1 && guessed_params[1].Bounds[1] == 2
	var temp_intervals []uint

	// if griddy 1_2 edge case, merge double ones surrounded by pixel items
	if case_grid_1_2 {
		temp_intervals = squashSurroundedDoubleOnes(intervals.Intervals, guessed_params[0].Bounds)
	}else{
		temp_intervals = intervals.Intervals
	}
	var new_intervals []uint = insertZeroesBetweenPixelIntervals(temp_intervals, guessed_params[0].Bounds)
	var interval_types []uint8 = createIntervalTypes(new_intervals, guessed_params[0].Bounds, guessed_params[1].Bounds)

	// fmt.Println("BEFORE" ,interval_types)
	// fmt.Println(new_intervals)
	squashConsecutiveSameIntervalTypes(&interval_types, &new_intervals)
	// fmt.Println("AFTER" ,interval_types)
	// fmt.Println(new_intervals)

	return CombinedList{new_intervals, interval_types} 	
}

/*
	Given bounds of pixel interval item, make a new slice which is a copy of intervals, except
	there are 0-sized interval items inserted between each pair of adjacent bounds-conforming intervals
*/
func insertZeroesBetweenPixelIntervals(intervals []uint, pixel_bounds [2]int) []uint {
	// allocating double memory because due to logic of the function there can be at most twice the elements cuz of zero insertions
	new_intervals := make([]uint, 0, len(intervals) * 2 + 1)
	edge_left, edge_right := intervals[0], intervals[len(intervals) - 1]

	// putting left edge interval into the list
	new_intervals = append(new_intervals, edge_left)
	// putting first non-edge interval into the list
	new_intervals = append(new_intervals, intervals[1])

	for i := 2; i < len(intervals) - 1; i++{
		left, right := int(intervals[i - 1]), int(intervals[i])

		var left_belongs bool = pixel_bounds[0] <= left && left <= pixel_bounds[1]
		var right_belongs bool = pixel_bounds[0] <= right && right <= pixel_bounds[1]
		// put zero only if left and right both belong to the given bounds
		if left_belongs && right_belongs {
			new_intervals = append(new_intervals, 0)
		}

		// put right item regardless if zero was inserted or not
		new_intervals = append(new_intervals, uint(right))
	}

	// putting right edge interval into the list
	new_intervals = append(new_intervals, edge_right)
	return new_intervals
}

/*
	Given an intervals slice (with edges) and bounds of interval candidate entry
	Make a new slice, where each [1,1] sequence surrounded by items belonging to the entry bounds 
	is squashed to a single item '2'

	Example: 
		bounds = [6,8]
		intervals = [1,3,6,1,1,7,1,6,1,1,7,1,1,1,6]
	Then result will be:
		[1,3,6,2,7,1,6,2,7,1,1,1,6]

*/
func squashSurroundedDoubleOnes(intervals []uint, pixel_bounds [2]int) []uint {
	result := make([]uint, 0 , len(intervals))

	// if there is not enough elements for even one iteration of sliding window, then return early
	if len(intervals) < 6 {
		result = result[0:len(intervals)]
		copy(result, intervals)
		return result
	}


	// putting left edge interval into the list
	result = append(result, intervals[0])
	// putting first non-edge element to the list
	result = append(result, intervals[1])

	var i int
	for i = 2; i < len(intervals) - 3; i++ {
		left,    right    := intervals[i - 1], intervals[i + 2]
		midleft, midright := intervals[i + 0], intervals[i + 1]

		var middle_ones bool = midleft == 1 && midright == 1
		var left_in_range bool  = int(left ) >= pixel_bounds[0] && int(left ) <= pixel_bounds[1]
		var right_in_range bool = int(right) >= pixel_bounds[0] && int(right) <= pixel_bounds[1]	

		all_conditions_met := middle_ones && left_in_range && right_in_range
		if all_conditions_met {      
			result = append(result, 2)
			i += 1
		}else{
			result = append(result, midleft)
		}

	}

	// append remaining items that werent covered by sliding window (which includes right edge interval)
	// cant just add 3 last items, because if last window iteration hits, there would be too many elements in result
	for ;i < len(intervals); i++ {
		result = append(result, intervals[i])
	}


	return result
}

/*
	Creates interval types slice based on pre-processed intervals slice and bounds of pixel and gridline items.


	1. Mark all items as unknown
	2. Mark pixel-conforming items as pixels
	3. Mark gridline-conforming items ajacent to (or between) pixels as gridlines
	4. potentially remove sequences that are too short 
	  (must make sure that longer exist, otherwise cant remov the short ones)
	5. Squash substrings of the same type into single items.

*/

func createIntervalTypes(intervals []uint, pixel_bounds, gridline_bounds [2]int) []uint8 {
	interval_types := make([]uint8, len(intervals))
	fillWithUnknowns(interval_types)
	markPixelsOnly(interval_types, intervals, pixel_bounds)
	markSurroundedGridlines(interval_types, intervals, gridline_bounds)
	return interval_types

}

/*
	fills interval_types slice in place with INTERVAL_UNKNOWN value
*/
func fillWithUnknowns(interval_types []uint8){
	for i := range interval_types {
		interval_types[i] = INTERVAL_UNKNOWN
	}
}

/*
	sets interval_types[i] to INTERVAL_PIXEL value on positions where intervals[i] is within provided pixel bounds
	interval types slice is changed in place.

	Ignores first and last interval in the collection. 

*/
func markPixelsOnly(interval_types []uint8, intervals []uint, pixel_bounds [2]int){
	for i:= 1; i< len(intervals) - 1; i++{
		var belongs_to_bounds bool = pixel_bounds[0] <= int(intervals[i]) && int(intervals[i]) <= pixel_bounds[1]
		if belongs_to_bounds{
			interval_types[i] = INTERVAL_PIXEL
		}
	}
}

/*
	sets interval_types[i] to INTERVAL_GRID value on positions where intervals[i] is within provided grid bounds
	and at positions i-1 and i+1 the interval_types is set to INTERVAL_PIXEL
	interval types slice is changed in place

*/

func markSurroundedGridlines(interval_types []uint8, intervals []uint, gridline_bounds [2]int){
	for i:= 1; i< len(intervals) - 1; i++{
		var belongs_to_grid bool = gridline_bounds[0] <= int(intervals[i]) && int(intervals[i]) <= gridline_bounds[1]
		var surrounded_by_pixels bool = interval_types[i-1] == INTERVAL_PIXEL && interval_types[i+1] == INTERVAL_PIXEL

		if belongs_to_grid && surrounded_by_pixels {
			interval_types[i] = INTERVAL_GRID
		}
	}
}

/*
	Given pointers to slices of interval types and intervals, modify them and reslice in place such as that:
	Any sequence of N consecutive items where interval type is the same is replaced with
	a 1 item of the same type with length equal to sum of replaced items

	Example input:
		Intervals:     [0,6,1,6,1,2,3,5]
		IntervalTypes: [0,0,1,2,1,0,1,1]

	Example result:
		Intervals:	   [6,1,6,1,2,8]
		IntervalTypes: [0,1,2,1,0,1]
*/

func squashConsecutiveSameIntervalTypes(interval_types *[]uint8, intervals *[]uint){
	new_count := 0

	for i := 1; i<len(*interval_types); i++ {
		if (*interval_types)[new_count] == (*interval_types)[i] {
			(*intervals)[new_count] += (*intervals)[i]

		}else{
			new_count += 1
			(*interval_types)[new_count] = (*interval_types)[i]
			(*intervals)[new_count] = (*intervals)[i]

		}
	}

	resliced_intervals := (*intervals)[0:new_count+1]
	resliced_types := (*interval_types)[0:new_count+1]

	*intervals = resliced_intervals
	*interval_types = resliced_types

}