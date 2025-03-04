package kuwahara

import "image"

/*
getGreyscaledChannel makes a float32 slice. Slice values represent a greyscaled version of input image. 
The returned slice assumes row major order.
Greyscale algorithm used: Y = 0.299 * R + 0.587 * G + 0.114 *B
*/
func getGreyscaledChannel(img *image.RGBA) []float32 {
	var pixel_count int = img.Rect.Dy() * img.Rect.Dx()
	greyscaled := make([]float32, pixel_count)

	for y := 0 ; y < img.Rect.Dy() ; y++ {
		for x := 0 ; x < img.Rect.Dx() ; x++ {
			flat_id := img.PixOffset(x + img.Rect.Min.X, y + img.Rect.Min.Y)
			greyscale_id := y * img.Rect.Dx() + x

			greyscaled[greyscale_id] = (
				0.299 * float32(img.Pix[flat_id + 0]) +
				0.587 * float32(img.Pix[flat_id + 1]) +
				0.114 * float32(img.Pix[flat_id + 2]))
		}
	}

	return greyscaled

}

/* 
	getSplitChannels seturns array of 3 slices of uint8, each slice has different memory block
	each uint8 slice corresponds to flattened R, G, B channel of the image, alpha is discarded
*/

func getSplitChannels(img *image.RGBA) [3][]uint8 {
	var channel_size int = img.Rect.Dy() * img.Rect.Dx()

	channels := [3][]uint8{ 
		make([]uint8, channel_size),
		make([]uint8, channel_size),
		make([]uint8, channel_size),
	}

	for y := 0 ; y < img.Rect.Dy() ; y++ {
		for x := 0 ; x < img.Rect.Dx() ; x++ {
			flat_id := img.PixOffset(x + img.Rect.Min.X, y + img.Rect.Min.Y)
			id_channel := y * img.Rect.Dx() + x

			channels[0][id_channel] = img.Pix[flat_id + 0]
			channels[1][id_channel] = img.Pix[flat_id + 1]
			channels[2][id_channel] = img.Pix[flat_id + 2]
		}

	}
	return channels
}