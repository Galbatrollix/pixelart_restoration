package gridlines

import (
	"slices"
	"math"
	"fmt"
)
import (
	"pixel_restoration/types"
)



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

	interval_ranges := getIntervalRanges(interval_counts)
	candidate1 := mostCommonIntervalRange(interval_ranges)

	interval_ranges_modified := rangesWithCollisionsZeroed(interval_ranges, candidate1)
	candidate2 := mostCommonIntervalRange(interval_ranges_modified)
	// todo hide once testing is done
	// types.PrintRangeArrayZeroFiltered(interval_ranges)
	// fmt.Printf("%-v,%-v\n", candidate1, candidate2)

	var one_involved bool = candidate1.Bounds[0] == 1 || candidate2.Bounds[0] == 1


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

	fmt.Printf("Pixel Guess: %-v\nGrid guess: %-v\n", pixel_guess, grid_guess)

	return pixel_guess, grid_guess
}


func guessParametersWithOne(intervals types.IntervalList, interval_counts []int,
	 candidates [2]types.IntervalRangeEntry) (types.IntervalRangeEntry, types.IntervalRangeEntry) {

	// Preparing statistics and interval entry values
	var candidate_smaller, candidate_bigger types.IntervalRangeEntry
	if candidates[0].Mean > candidates[1].Mean{
		candidate_smaller, candidate_bigger = candidates[1], candidates[0]
	}else{
		candidate_smaller, candidate_bigger = candidates[0], candidates[1]
	}

	entry_1_2 := types.IntervalRangeEntry{
		Bounds: [2]int{1,2},
		Count: candidate_smaller.Count,
		Mean: candidate_smaller.Mean,
	}
	var count_ones_only int
	var count_twos_only int
	if len(interval_counts) > 1 {
		count_ones_only = interval_counts[1]
	}else{
		count_ones_only = 0
	}

	if len(interval_counts) > 2 {
		count_twos_only = interval_counts[2]
	}else{
		count_twos_only = 0
	}
	entry_1_1 := types.IntervalRangeEntry{
		Bounds: [2]int{1,1},
		Count : count_twos_only,
		Mean: 1.0,
	}
	entry_2_2 := types.IntervalRangeEntry{
		Bounds: [2]int{2,2},
		Count: count_ones_only,
		Mean: 2.0,
	}

	// PROBABLY go to one and two if other candidate count is just 1, no use in calculating scores on low values
	var second_empty bool = candidates[1].Count <= 1
	if second_empty {
		goto ConsiderOnlyOneAndTwo
	}
	
	
	{		

		// multiply by scores involving big candidate to get more accurate prediction
		var big_bias float64 = 1 + 0.1 * candidate_bigger.Mean

		// score for [1,2] gridline size and second candidate as pixel size		
		score_candidate_1_2 := float64(alternatingCandidatesRunlengthScore_1_2(intervals, candidate_bigger)) * big_bias
		score_candidate_0_1 := float64(alternatingCandidatesRunlengthScore_0_1(intervals, candidate_bigger)) * big_bias
		score_only_1_2 := float64(singleCandidateRunlengthScore(intervals, entry_1_2, 1)) 
		// fmt.Println("SCORE CANDIDATE 1_2 ",score_candidate_1_2)
		// fmt.Println("SCORE CANDIDATE 0_1", score_candidate_0_1)
		// fmt.Println("SCORE ONLY 1_2", score_only_1_2)
		
		highest_score := max(score_candidate_0_1, score_candidate_1_2, score_only_1_2)

		switch highest_score {
		case score_only_1_2:
			goto ConsiderOnlyOneAndTwo
		case score_candidate_1_2:
			return candidate_bigger, entry_1_2
		default: // score_candidate_0_1
			return candidate_bigger, entry_1_1

		}

	}


	ConsiderOnlyOneAndTwo:{

		score_1 := singleCandidateRunlengthScore(intervals, entry_1_1, 0)
		score_2 := singleCandidateRunlengthScore(intervals, entry_2_2, 0)
		score_alternating_1_2 := alternatingCandidatesRunlengthScore(
			intervals,
			[2]types.IntervalRangeEntry{entry_1_1, entry_2_2},
		)

		// fmt.Println("CONSIDERING ONLY ONE AND TWO")
		// fmt.Println("Scores: ", score_1, score_2, score_alternating_1_2)
		highest_score := max(score_1, score_2, score_alternating_1_2)

		switch highest_score {
		case score_alternating_1_2:
			return entry_2_2, entry_1_1
		case score_2:
			return entry_2_2, types.GetZeroRangeEntry()
		default: // score_1
			return entry_1_1, types.GetZeroRangeEntry()

		}

	}

}

func guessParametersNoOne(intervals types.IntervalList, inteval_counts []int, 
	candidates [2]types.IntervalRangeEntry) (types.IntervalRangeEntry, types.IntervalRangeEntry){

	// if no intervals were left for calculatinng second candidate, then candidate 1 is assumed to be pixel size
	var second_empty bool = candidates[1].Count == 0
	if second_empty {
		return candidates[0] , types.GetZeroRangeEntry()
	}

	var candidate_smaller, candidate_bigger types.IntervalRangeEntry
	if candidates[0].Mean > candidates[1].Mean{
		candidate_smaller, candidate_bigger = candidates[1], candidates[0]
	}else{
		candidate_smaller, candidate_bigger = candidates[0], candidates[1]
	}

	// Edge case for happens when one candidate is roughly 2x the other one.
	// If that happens, assume its a gridless image and guess which guess is the correct one
	// 1. Second candidate count must be at least 30% of the first candidate count
	double_width_count_condition := float64(candidates[1].Count) / float64(candidates[0].Count) >= 0.3
	// 2. Average of bigger of the candidates must be very close to double of the average of smaller candidate
	double_width_size_condition := candidate_smaller.Mean * 2.0 + 1.0 >= candidate_bigger.Mean
	if double_width_count_condition && double_width_size_condition {
		
		// If arrangement of larger objects suggests gridline mismatch, choose smaller item as pixel
		var properly_aligned bool = isDoubleSizedIntervalAligned(intervals, candidate_bigger)
		if ! properly_aligned {
			return candidate_smaller, types.GetZeroRangeEntry()
		}


		// if alignment check succeeded, then guess based on runlength scores
		bigger_score := singleCandidateRunlengthScore(intervals, candidate_bigger, 1)
		smaller_score := singleCandidateRunlengthScore(intervals, candidate_smaller, 1)
		if bigger_score > smaller_score {
			return candidate_bigger, types.GetZeroRangeEntry()
		}else if smaller_score > bigger_score {
			return candidate_smaller, types.GetZeroRangeEntry()
		}else{ // equal scores, choose 1st candidate
			return candidates[0], types.GetZeroRangeEntry()
		}
	}

	// Base case: calculate score for both cadidates and score for candidate 1 and 2 alternating, and choose the winner
	score1 := singleCandidateRunlengthScore(intervals, candidates[0], 0)
	score2 := singleCandidateRunlengthScore(intervals, candidates[1], 0)

	score_alternating := alternatingCandidatesRunlengthScore(intervals, candidates)
	//fmt.Printf("Score1: %d\nScore2: %d\nScoreAlternating: %d\n",  score1,score2,score_alternating)

	highest_score := max(score1, score2, score_alternating)
	//fmt.Println("Highest Score: ", highest_score)
	switch highest_score {
	case score_alternating:
		return candidate_bigger, candidate_smaller
	case score1:
		return candidates[0], types.GetZeroRangeEntry()
	default: // score2, currently will never execute but might if algorithm for score is modified
		return candidates[1], types.GetZeroRangeEntry()

	}

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
		entry_midpoint := (entry_max - entry_min ) / 2 + entry_min
		var violates_ratio bool = sample_big && (ratio_min <= entry_midpoint && entry_midpoint <= ratio_max)

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
	Calculates special case for alternating interval items score, where width items are only 1 or 2 sized
	and there is additional pre-processing performed.
*/ 
func alternatingCandidatesRunlengthScore_1_2(intervals types.IntervalList, candidate types.IntervalRangeEntry) int {
	intervals_noedges := intervals.Intervals[1:intervals.TotalCount-1]
	intervals_squashed := squashSurroundedDoubleOnesIntervals(intervals_noedges, candidate.Bounds)
	var runlengths []int = alternatingCandidatesRunlengths(
		intervals_squashed,
		[2][2]int{candidate.Bounds, {1,2}},
	)
	score := sumLargerThan(runlengths, 1)
	// test functionality, considers edges in 2-item sequences with 50% weight
	score_with_ones := sumLargerThan(runlengths, 0)
	score += ( score_with_ones - score) / 2 
	return score

}

/*
	Calculates special case for alternating interval items score, where width items are only 1 are absent completely
	This means that for candidate in range [4,5], the following sequence is correct:
	[5,5,4,1,4,5,1,4,5,1,5]

*/ 
func alternatingCandidatesRunlengthScore_0_1(intervals types.IntervalList, candidate types.IntervalRangeEntry) int {
	intervals_noedges := intervals.Intervals[1:intervals.TotalCount-1]
	var runlengths []int = singleCandidateWithOneRunlengths(intervals_noedges, candidate.Bounds)

	score := sumLargerThan(runlengths, 1)
	// test functionality, considers edges in 2-item sequences with 50% weight
	score_with_ones := sumLargerThan(runlengths, 0)
	score += ( score_with_ones - score) / 2 
	return score

}


/*
	Calculates score for comparing selected interval range entry candidate with other candidates and arrangements
*/
func singleCandidateRunlengthScore(intervals types.IntervalList, candidate types.IntervalRangeEntry, min_runglength int) int {
	intervals_noedges := intervals.Intervals[1:intervals.TotalCount-1]

	var runlengths []int = singleCandidateRunlengths(intervals_noedges, candidate.Bounds)
	score := sumLargerThan(runlengths, min_runglength)

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
	Given list of intervals without first and last elements, and bounds of an interval range entry,
	Find all pairs of intervals such that:
		 each pair consists of either:
		 1. one interval conforming to interval range entry and one 1-length interval
		 2.	Both intervals conforming to interval range entry
	Returned slice contains lenghts of consecutive runlengths of such (overlapping) pairs 
	Example : 
		bounds = [4,5]
		intervals = [2,5,1,7,4,4,5,4,1,5,7,1,5]

	Then the substrings with all pairs conforming would be the following:
		[5,1]
		[4,4,5,4,1,5]
		[1,5]
	Hence the result slice would be: [1,5,1]

	Made for spotting edge case between gridless and 1-width griddy image

*/ 

func singleCandidateWithOneRunlengths(intervals_noedges []uint, candidate_bounds [2]int) []int {
	lookup_size := len(intervals_noedges) - 1
	lookup := make([]bool, lookup_size)

	for i := 0; i<lookup_size; i++ {
		left := int(intervals_noedges[i])
		right := int(intervals_noedges[i+1])

		var left_belongs_to_candidate bool  = candidate_bounds[0] <= left && left <= candidate_bounds[1]
		var right_belongs_to_candidate bool = candidate_bounds[0] <= right && right <= candidate_bounds[1]
		var left_is_1 bool = left == 1
		var right_is_1 bool = right == 1


		var ok_edge_found bool =  left_belongs_to_candidate && right_belongs_to_candidate || 
						   						  left_is_1 && right_belongs_to_candidate ||
						   						  right_is_1 && left_belongs_to_candidate
		lookup[i] = ok_edge_found

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
	lookup := candidateBelongsLookup(intervals_noedges, cadidate_bounds)

	runlengths := consecutiveTrueRunlengths(lookup)


	return runlengths
}


/*
	Given a slice of intervals (without edges) and bounds of interval candidate entry,
	return a same-lengthed slice of bools 
	position is set to true if interval under the same index belongs to candidate bounds, otherwise false
*/
func candidateBelongsLookup(intervals_noedges[]uint, cadidate_bounds [2]int) []bool {
	lookup := make([]bool, len(intervals_noedges))
	for i, interval := range intervals_noedges {
		var is_in_range bool = int(interval) >= cadidate_bounds[0] && int(interval) <= cadidate_bounds[1]
		lookup[i] = is_in_range
	}
	return lookup
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

func sliceSumU(slice []uint) uint {
	sum := uint(0) 
	for _, value := range slice {
		sum += value
	}
	return sum
}

/*
	Given slice of bools, return slice of ints that specifies all indexes in bool slice where there is true
*/
func indexesOfTrue(slice []bool) []int {
	result :=  make([]int, 0, len(slice))
	for i, isTrue := range slice {
		if isTrue{
			result = append(result, i)
		}
	}

	return result
}

/*
    Assuming image is gridless, check if provided interval candidate has alignment problems.

    With input list of [8,4,8,8] and candidate 8, the return value will be False, because
    first two 8s are not aligned with each other. (there is 4 units between them, but should be around 8)

    With input list of 8,4,4,3,4,8 and candidate 8, the return value will be True, because
    4,4,3,4 add up to 15, which is very close to 16 (2 * 8)
 */
func isDoubleSizedIntervalAligned(intervals types.IntervalList, candidate types.IntervalRangeEntry) bool{
	var intervals_noedges []uint = intervals.Intervals[1:intervals.TotalCount - 1]
	var belongment_array []bool = candidateBelongsLookup(intervals_noedges, candidate.Bounds)
	var indexes_true []int = indexesOfTrue(belongment_array)

	// any <gap/mean - floor(gap/mean)> between failure range bounds means that intervals are incorrectly aligned
	failure_range:= [2]float64{0.2, 0.8}

	for i:=1; i< len(indexes_true); i++{
		left_id, right_id := indexes_true[i-1], indexes_true[i]
		segment_inbetween := intervals_noedges[left_id+1:right_id]
		total_segment_length := sliceSumU(segment_inbetween)

		alignment_score_exact := (float64(total_segment_length) / candidate.Mean ) 
		alignment_score_truncated := alignment_score_exact - math.Floor(alignment_score_exact)
		if alignment_score_truncated >= failure_range[0] && alignment_score_truncated <= failure_range[1]{
			return false
		}

	}
	return true
}


/*
	Given an interval slice (no edges) and bounds of interval candidate entry
	Make a new slice, where each [1,1] sequence surrounded by items belonging to the entry bounds 
	is squashed to a single item '2'

	Example: 
		bounds = [6,8]
		intervals = [1,3,6,1,1,7,1,6,1,1,7,1,1,1,6]
	Then result will be:
		[1,3,6,2,7,1,6,2,7,1,1,1,6]

*/
func squashSurroundedDoubleOnesIntervals(intervals_noedges []uint, cadidate_bounds [2]int) []uint {
	result := make([]uint, 0 , len(intervals_noedges))
	if len(intervals_noedges) < 4 {
		result = result[0:len(intervals_noedges)]
		copy(result, intervals_noedges)
		return result
	}

	//append first item since sliding window won't interact with it
	result = append(result, intervals_noedges[0])
	var i int
	for i = 1; i < len(intervals_noedges) - 2; i++ {
		left,    right    := intervals_noedges[i - 1], intervals_noedges[i + 2]
		midleft, midright := intervals_noedges[i + 0], intervals_noedges[i + 1]

		var middle_ones bool = midleft == 1 && midright == 1
		var left_in_range bool  = int(left ) >= cadidate_bounds[0] && int(left ) <= cadidate_bounds[1]
		var right_in_range bool = int(right) >= cadidate_bounds[0] && int(right) <= cadidate_bounds[1]	

		all_conditions_met := middle_ones && left_in_range && right_in_range
		if all_conditions_met {      
			result = append(result, 2)
			i += 1
		}else{
			result = append(result, midleft)
		}

	}

	// append remaining items that werent covered by sliding window
	// cant just add 2 last items, because if last window hits, there would be too much elements in result
	for ;i < len(intervals_noedges); i++ {
		result = append(result, intervals_noedges[i])
	}

	return result
}