package contrast


import "pixel_restoration/common"

import (
	"sort"
)

/*
	ClipTop: float32
		(After filtering out positions with 0 edges)
		removes <ClipTop * 100%> most common edge positions before calculating minimum count for qualifying edge
	CutoffMultiplier: float32
		(After applying both clip top)
		The lowest permissable edge count for a position is set to CutoffMultiplier * 100% *<most common edge position count>)

	Constraints:
		0 <= ClipTop < 1.0
		0 <= CutoffMultipler < 1.0
*/
type MostFrequentParams struct {
	ClipTop float32
	CutoffMultiplier float32
}


func GetBaseMostFrequentParams() MostFrequentParams {
	return MostFrequentParams{
		ClipTop: 0.2,
		CutoffMultiplier: 0.3,
	}
}

/*
	SelectMostFrequent selects indexes most likely to be correctly detected edges by heuristics described in params struct.
	Takes <edge_counts[id] = count> mapping slice as a parameter and returns a slice of edge ids that passed the heuristic.
*/
func SelectMostFrequent(edge_counts []uint, params MostFrequentParams) []int {
	if len(edge_counts) == 0 {
		return []int{}
	}

	counter := edgeCountsSortedNonzero(edge_counts)
	
	clip_amount := int(params.ClipTop * float32(len(counter)))
	// index of highest value after clipping
	sample_val_id := min(clip_amount, len(counter) - 1)

	// all indexes with less count than this will not pass
	threshold_quantity := uint(float32(counter[sample_val_id]) * params.CutoffMultiplier)

	result := make([]int, 0, len(counter))

	// put thresholded counts in the result
	for i, count := range edge_counts{
		if count >= threshold_quantity {
			result = append(result, i)
		}
	}

	return result
}


// returns a slice of edge counts with filtered out 0 values, sorted in descending order
func edgeCountsSortedNonzero(edge_counts []uint) []uint {
	positive_count := common.CountNonZeroU(edge_counts)
	new_counts := make([]uint, positive_count)
	populateNonzero(edge_counts, new_counts)

	comparator :=  func(i, j int) bool {
		return new_counts[i] > new_counts[j]
	}
	sort.Slice(new_counts, comparator)

	return new_counts

}
// Put all nonzero items from edge counts to positive counts, memory for positive counts is already allocated
func populateNonzero(edge_counts, positive_counts []uint){
	new_id := 0
	for _, count := range edge_counts{
		if count != 0 {
			positive_counts[new_id] = count
			new_id += 1
		}
	}
}

