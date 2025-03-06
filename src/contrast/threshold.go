package contrast

import (
	"image"
)
import (
	"pixel_restoration/images"
)

/*

	Given a grayscale image and threshold value, returns a new grayscale image.

	Resulting image has the same shape and all positions (y, x) 
	where value of original image (y, x) is higher or equal to threshold 
	are set to 255 and positions with values lower than threshold are set to 0

*/
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