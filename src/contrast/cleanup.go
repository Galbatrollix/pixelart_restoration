package contrast

import "image"

import "pixel_restoration/images"


/* 
	Given a binarized (tresholded with min peak height) gray image od edges,
	Return a new image where random single pixels, diagonal and against-the-grain lines are removed.
	Second return is a count of manipulated pixels. (High may sugest noisy and watermarkey image)
*/
func CleanupEdgeArtifacts(edges_binary *image.Gray) (*image.Gray, int){
	result := images.GrayscaleGetNormalized(edges_binary)
	width, height := result.Rect.Dx(), result.Rect.Dy()

	count_changed := 0

	for y:= 1; y<height - 1;y++{
		for x:=0; x<width; x++{
			above_id := (y - 1) * width + x
			curr_id := y * width + x 
			below_id := (y + 1) * width + x

			if(result.Pix[above_id] == 0 && result.Pix[below_id] == 0 && result.Pix[curr_id] == 255){
				result.Pix[curr_id] = 0
				count_changed += 1
			}

		}
	}

	return result, count_changed
}	