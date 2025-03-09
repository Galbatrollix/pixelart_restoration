package types

import "fmt"

/* 
	IntervalRangeEntry is used to describe a category of interval items considered 
	"the same" for gridline parameter guessing algorithm.
	Bounds:
		[MIN, MAX] inclusive range of interval lengths for the entry
	Count:
		Number of detected intervals falling within the bounds
	Mean:
		Average length of interval falling within the bounds

 */

type IntervalRangeEntry struct {
	Bounds [2]int
	Count int
	Mean float64
}

func GetZeroRangeEntry() IntervalRangeEntry {
	return IntervalRangeEntry{[2]int{0,0},0,0}
}

/*
	Pretty prints array of entries by filtering zero and using pretty formatting
*/
func PrintRangeArrayZeroFiltered(entries []IntervalRangeEntry) {
	for _, entry := range entries{
		if entry.Count == 0{
			continue
		}
		items_available := []int{}
		for i:= entry.Bounds[0]; i<= entry.Bounds[1]; i++{
			items_available = append(items_available, i)
		}

		fmt.Printf("%v:\n", entry.Bounds)
		fmt.Printf("    Elements: %v\n", items_available)
		fmt.Printf("    Mean: %f\n", entry.Mean)
		fmt.Printf("    Count: %d\n", entry.Count)
	}
}