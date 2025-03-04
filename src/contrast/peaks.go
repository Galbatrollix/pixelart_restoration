package contrast

import (
	"slices"
	"image"
)
import (
	"pixel_restoration/images"
)


/*
	CalculateMinPeakHeight approximates minimum peak height for optimal gridline detection
    by analyzing all found distances.

    Around 50 is generally good peak height value.
    This function will return 58 if max possible distance was detected
    Upper bound on returned value can be configured in params package
*/

func CalculateMinPeakHeight(all_distances []uint8, min_peak_height_limit float64) uint8{
	const max_possible_dist float64 = 255.0
	const base_height float64 = 58.0

	var dist_max uint8 = slices.Max(all_distances)
	min_peak_height := base_height * float64(dist_max) / max_possible_dist
	return uint8(min(min_peak_height_limit, min_peak_height) + 0.5)
}



func ThresholdWithMinHeight(distances *image.Gray, min_peak_height uint8) *image.Gray{
	distances_new := images.GrayscaleGetNormalized(distances)

	item_count := distances_new.Rect.Dx() * distances_new.Rect.Dy()

	for i:= 0; i<item_count; i++{
		if distances_new.Pix[i] >= min_peak_height {
			distances_new.Pix[i] = 255
		}else{
			distances_new.Pix[i] = 0
		}
	}
	return distances_new
}