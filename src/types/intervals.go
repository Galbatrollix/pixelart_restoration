package types

/*

	Interval list holds one attribute: Intervals slice
	Intervals:
		slice of N intervals of consecutive pixels without edges between them.

	For example:
		Intervals: [4,3,5,4]
		This means that there are 16 pixels in the image, split into 4 consecutive segments, which represents the following:
		4 pixels > edge > 3 pixels > edge > 5 pixels > edge 4 pixels
*/

type IntervalList struct{
	Intervals []uint
}

/*
	Constructs IntervalList from sorted slice of edge positions and max dimension length
*/
func IntervalListFromSortedEdgeIndexes(edge_indexes []int, dim_length int) IntervalList {
	var new_length int = len(edge_indexes) + 1
	result := IntervalList{
		Intervals: make([]uint, new_length),
	}

	accumulated_length := 0
	for i , edge_pos := range edge_indexes{
		result.Intervals[i] = uint(edge_pos - accumulated_length)
		accumulated_length = edge_pos
	}

	result.Intervals[new_length - 1] = uint(dim_length - accumulated_length)

	return result
}

