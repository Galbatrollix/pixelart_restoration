package gridlines

import (
	"slices"
	"math"
	"fmt"
)
import (
	"pixel_restoration/types"
)


//todo: const big_bias = ...

/*
	Returns interval range for pixel and gridline if griddy image,
	otherwise returns interval and zeroed range entry, 
	in case of being unable to detect any range (not enough intervals) returns two zeroed ranges
*/

func GuessGridlineParameters(intervals types.IntervalList) (types.IntervalRangeEntry, types.IntervalRangeEntry) {
	if intervals.TotalCount < 3 {
		return types.GetZeroRangeEntry(), types.GetZeroRangeEntry()
	}

	interval_counts := getIntervalCounts(intervals)
	largest_interval := len(interval_counts) - 1
	_ = largest_interval
	// fmt.Println("Interval counts: ",interval_counts)
	// fmt.Println("Largest interval:",largest_interval)

	interval_ranges := getIntervalRanges(interval_counts)
	candidate1 := mostCommonIntervalRange(interval_ranges)

	interval_ranges_modified := rangesWithCollisionsZeroed(interval_ranges, candidate1)
	candidate2 := mostCommonIntervalRange(interval_ranges_modified)
	// todo hide once testing is done
	//types.PrintRangeArrayZeroFiltered(interval_ranges)
	fmt.Printf("%-v,%-v\n", candidate1, candidate2)

	var one_involved bool = candidate1.Bounds[0] == 1 || candidate2.Bounds[0] == 1

	fmt.Println(one_involved)

	var pixel_guess, grid_guess types.IntervalRangeEntry
	if one_involved {
		pixel_guess, grid_guess = guessParametersWithOne(
			intervals, interval_counts,
			[2]types.IntervalRangeEntry{candidate1, candidate2},
		)
	}else{
		pixel_guess, grid_guess = guessParametersNoOne(
			intervals, interval_counts,
			[2]types.IntervalRangeEntry{candidate1, candidate2},
		)
	}

	return pixel_guess, grid_guess
}


func guessParametersWithOne(intervals types.IntervalList, inteval_counts []int,
	 candidates [2]types.IntervalRangeEntry) (types.IntervalRangeEntry, types.IntervalRangeEntry) {




	return types.IntervalRangeEntry{[2]int{0,0},0,0}, types.IntervalRangeEntry{[2]int{0,0},0,0}
}

func guessParametersNoOne(intervals types.IntervalList, inteval_counts []int, 
	candidates [2]types.IntervalRangeEntry) (types.IntervalRangeEntry, types.IntervalRangeEntry){

	// if no intervals were left for calculatinng second candidate, then candidate 1 is assumed to be pixel size
	var second_empty bool = candidates[1].Count == 0
	if second_empty {
		zero_range := types.IntervalRangeEntry{[2]int{0,0},0,0}
		return candidates[0] , zero_range
	}
	// TODO edge case for "roughly double", potentially change the filtrting out algoritm too, to accomodate that
	// probably remove if midpoint of entry is within extended range, not entire entry
	/*


	!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

	*/


	// Base case: calculate score for both cadidates and score for candidate 1 and 2 alternating, and choose the winner
	score1 := singleCandidateRunlengthScore(intervals, candidates[0])
	score2 := singleCandidateRunlengthScore(intervals, candidates[1])

	score_alternating := alternatingCandidatesRunlengthScore(intervals, candidates)
	fmt.Printf("Score1: %d\nScore2: %d\nScoreAlternating: %d\n",  score1,score2,score_alternating)

	return types.IntervalRangeEntry{[2]int{0,0},0,0}, types.IntervalRangeEntry{[2]int{0,0},0,0}
}

/*
	Makes a lookup where index is a size of interval and value is the count of detected intervals of this size
	Lookup size is enough to hold up to the largest encountered interval. (so size of largest interval + 1)
	First and last (if any) intervals are removed before the calculation
*/
func getIntervalCounts(intervals types.IntervalList) []int {
	intervals_noedges := intervals.Intervals[1:intervals.TotalCount-1]

	largest_interval := int(slices.Max(intervals_noedges))
	counts := make([]int, largest_interval + 1)

	for _, value := range intervals_noedges{
		counts[value] += 1
	}
	return counts
}


/*
	Given a lookup with interval counts:
	
	Get slice of interval size ranges where each range describes range of sizes pertaining to the "same" item.
	Range window grows the higher the interval size.

	More info in definition of types.IntervalRangeEntry

	Example:
	[0,0], [1,2], [2,3], [3,5], [4,6], [5,7] (...) [14,17]

*/
func getIntervalRanges(interval_counts []int) []types.IntervalRangeEntry {
	// allocating memory
	deviation_ranges := make([]types.IntervalRangeEntry, len(interval_counts))

	// constructing bounds for ranges
	for i := range deviation_ranges {
		range_offset := getRangeOffset(i)

		start := i
		end := min(len(deviation_ranges) - 1, start + range_offset )
		deviation_ranges[i].Bounds = [2]int{ start, end }
	}

	// calculating means and counts for each range bound
	for i := range deviation_ranges {
		first, last := deviation_ranges[i].Bounds[0], deviation_ranges[i].Bounds[1]
		var sum float64 = 0.0
		var count int = 0
		for interval := first ; interval <= last; interval++ {
			count += interval_counts[interval]
			sum += float64(interval_counts[interval] * interval)
		}
		var mean float64 = 0.0
		if count != 0 {
			mean = sum / float64(count)
		}
		deviation_ranges[i].Mean = mean
		deviation_ranges[i].Count = count
	}

	return deviation_ranges
}


/* calculates length of range depending on its left bound */
func getRangeOffset(interval_base int) int {
	if interval_base <= 2 {
		return 1
	}else if interval_base <= 10 {
		return 2
	}else if interval_base <= 25 {
		return 3
	}else if interval_base <= 50 {
		return 4
	}else{
		return 5
	}
}


/*
	Returns entry with the highest count. 
	If multiple have the same count, returns the last of them
*/
func mostCommonIntervalRange(entries []types.IntervalRangeEntry) types.IntervalRangeEntry {
	var max_count int = 0
	var result types.IntervalRangeEntry
	for _, entry := range entries {
		if entry.Count >= max_count {
			max_count = entry.Count
			result = entry
		}
	}
	return result	
}

/*
	
	Creates a copy of entries slice and sets all values colliding with provided sample entry to count 0.

	Colliding items are: 
		- All items in range of sample
		- All items that would validate the constraint on maximum gridline-to-pixel ratio (1/2)
		(second point has some special cases around interval sizes 1 and 2)
	
*/
func rangesWithCollisionsZeroed(entries []types.IntervalRangeEntry, sample types.IntervalRangeEntry) []types.IntervalRangeEntry{
	new_entries := make([]types.IntervalRangeEntry, len(entries))
	copy(new_entries, entries)

	
	sample_min, sample_max := sample.Bounds[0], sample.Bounds[1]

	sample_midpoint := int(math.Round(sample.Mean))
	var sample_big bool = sample_midpoint >= 4
	ratio_min, ratio_max := sample_midpoint / 2 + 1, sample_midpoint * 2 - 1

	// fmt.Printf("Sample: %+v\n", sample)
	for i, entry := range new_entries {
		// Checking if items directly overlap
		entry_min, entry_max := entry.Bounds[0], entry.Bounds[1]
		var direct_overlap bool = entry_min <= sample_max && sample_min <= entry_max

		// Removing if item would violate gridline to pixel ratio and is big enough for that
		var violates_ratio bool = sample_big && (entry_min <= ratio_max && ratio_min <= entry_max)

		if direct_overlap || violates_ratio {
			new_entries[i].Count = 0
			new_entries[i].Mean = 0.0
			// fmt.Printf("Removed entry: %+v\n", entry)
		}
	}

	return new_entries
}

/*
	Calculates score for comparing alternating interval arrangement with other candidates
*/
func alternatingCandidatesRunlengthScore(intervals types.IntervalList, candidates [2]types.IntervalRangeEntry) int {
	intervals_noedges := intervals.Intervals[1:intervals.TotalCount-1]
	var runlengths []int = alternatingCandidatesRunlengths(
		intervals_noedges,
		[2][2]int{candidates[0].Bounds ,candidates[1].Bounds},
	)

	score := sumLargerThan(runlengths, 1)
	// test functionality, considers edges in 2-item sequences with 50% weight
	score_with_ones := sumLargerThan(runlengths, 0)
	score += ( score_with_ones - score) / 2 

	return score
}


/*
	Calculates score for comparing selected interval range entry candidate with other candidates and arrangements
*/
func singleCandidateRunlengthScore(intervals types.IntervalList, candidate types.IntervalRangeEntry) int {
	intervals_noedges := intervals.Intervals[1:intervals.TotalCount-1]

	var runlengths []int = singleCandidateRunlengths(intervals_noedges, candidate.Bounds)
	score := sumLargerThan(runlengths, 0)

	return score

}

/*
	Given list of intervals without first and last elements, and bounds of two interval range entries,
	Find all pairs of intervals such that:
		 each pair consists of one interval conforming to first range and one interval conforming to second range.
	Returned slice contains lenghts of consecutive runlengths of such (overlapping) pairs 
	Example : 
		bounds = [4,5], [1,2]
		intervals = [2,5,1,7,4,4,5,4,1,5,7,1,5]

	Then the substrings with all pairs conforming would be the following:
		[2,5,1]
		[4,1,5]
		[1,5]
	Hence the result slice would be: [2,2,1]

*/ 

func alternatingCandidatesRunlengths(intervals_noedges []uint, cadidate_bounds [2][2]int) []int {
	lookup_size := len(intervals_noedges) - 1
	lookup := make([]bool, lookup_size)

	first := cadidate_bounds[0]
	second := cadidate_bounds[1]

	for i := 0; i<lookup_size; i++ {
		left := int(intervals_noedges[i])
		right := int(intervals_noedges[i+1])

		var left_belongs_to_first bool = first[0] <= left && left <= first[1]
		var right_belongs_to_first bool = first[0] <= right && right <= first[1]
		var left_belongs_to_second bool = second[0] <= left && left <= second[1]
		var right_belongs_to_second bool = second[0] <= right && right <= second[1]

		var alternating_edge_found bool =  left_belongs_to_first && right_belongs_to_second || 
										   left_belongs_to_second && right_belongs_to_first 
		lookup[i] = alternating_edge_found

	}

	runlengths := consecutiveTrueRunlengths(lookup)
	return runlengths

}


/*
	Given list of intervals without first and last elements, and bounds of interval range entry,
	Find all substrings of intervals data that are composed only of items within candidate bounds.
	Returned slice contains lenghts of all such substrings. 
	Example : 
		bounds = [4,5]
		intervals = [2,5,1,7,4,4,5,4,3,5,7]

	Then relevant substrings would be the following:
		[5]
		[4,4,5,4]
		[5]
	Hence the result slice would be: [1,4,1]

*/ 
func singleCandidateRunlengths(intervals_noedges []uint, cadidate_bounds [2]int) []int {
	lookup := make([]bool, len(intervals_noedges))
	for i, interval := range intervals_noedges {
		var is_in_range bool = int(interval) >= cadidate_bounds[0] && int(interval) <= cadidate_bounds[1]
		lookup[i] = is_in_range
	}

	runlengths := consecutiveTrueRunlengths(lookup)


	return runlengths
}


/*
	Returns lengths of all true-only substrings in bool slice
*/
func consecutiveTrueRunlengths(slice []bool) []int{
	runlengths := make([]int, 0, len(slice)/2 + 1)

	previous_belongs := false
	runlength_count := 0
	for _, belongs := range slice{
		// if item doesnt belong to range, reset counting and append if it ended the streak
		if ! belongs {
			if previous_belongs {
				runlengths = append(runlengths, runlength_count)
			}
			previous_belongs = false
			runlength_count = 0
			continue
		}

		// Item belongs to range, then continue counting
		previous_belongs = true
		runlength_count += 1
	}

	// if loop ended and previous item was true, then append the accumulated sum to the runlengths
	if previous_belongs {
		runlengths = append(runlengths, runlength_count)
	}

	return runlengths
}

/*
 	Returns sum of all elements in slice that are larger than cutoff value 
*/
func sumLargerThan(slice []int, cutoff_value int) int{
	sum := 0
	for _, value := range slice {
		if value > cutoff_value {
			sum += value
		}
	}
	return sum
}