package contrast

import (
	"slices"
)


/*
	CalculateMinPeakHeight approximates minimum peak height for optimal gridline detection
    by analyzing all found distances.

    Around 50 is generally good peak height value.
    This function will return 58 if max possible distance was detected
    Upper bound on returned value can be provided in params struct
*/

func CalculateMinPeakHeight(all_distances []uint8, params PeakHeightParams) uint8{
	const max_possible_dist float64 = 255.0
	const base_height float64 = 58.0

	var dist_max uint8 = slices.Max(all_distances)
	min_peak_height := base_height * float64(dist_max) / max_possible_dist
	return uint8(min(params.MinPeakHeightLimit, min_peak_height) + 0.5)
}


/*
	Mean peak height limit allows to set the limit on the return value of min peak height function
	If function is about to return a value higher than this parameter, it gets truncated to this value.
*/
type PeakHeightParams struct {
	MinPeakHeightLimit float64
}

func GetBasePeakHeightParams() PeakHeightParams{
	return PeakHeightParams{
		MinPeakHeightLimit : 58.0,
	}
}