package contrast

import (
	"slices"
)
import (
	"pixel_restoration/params"	
)


/*
	CalculateMinPeakHeight approximates minimum peak height for optimal gridline detection
    by analyzing all found distances.

    Around 75 is generally good peak height value.
    This function will return 100 if max possible distance was detected
    Upper bound on returned value can be configured in params package
*/

func CalculateMinPeakHeight(all_distances []float32) float32{
	const max_possible_color_diff float32 = 443.0
	const base_height float32 = 100.0

	var dist_max float32 = slices.Max(all_distances)
	min_peak_height := base_height * dist_max / max_possible_color_diff
	return min(params.MinPeakHeightLimit, min_peak_height)
}