package images

import (
	"image"
    "image/png"
    _ "image/jpeg"
    "image/draw"
	"os"
)


func ImageLoadFromFile(filepath string) (*image.RGBA, error){
	// Read image from file that already exists
	infile, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer infile.Close()

	imageData, _, err := image.Decode(infile)
	if err != nil {
		return nil, err
	}

	return toRGBA(imageData) , nil
}

func ImageSaveToFile(filepath string, img *image.RGBA) error {
	outfile, err := os.Create(filepath)
 	if err != nil {
   		return err
    }
    defer outfile.Close()

	err = png.Encode(outfile, img)
	if err != nil {
   		return err
	}

	return nil
}


func toRGBA(src image.Image) (*image.RGBA){
	bounds := src.Bounds()
	converted := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(converted, bounds, src, bounds.Min, draw.Src)
	return converted
}