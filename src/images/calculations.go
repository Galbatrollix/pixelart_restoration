package images
import (
	"image"
	"math"
	"sort"
)


func CalculateAverageColor(img *image.RGBA) [3]uint8 {
	var result = [3]uint8{}
	channels := ImageGetSplitChannels(img)

	for i := 0; i<3; i++ {
		result[i] = meanOfSlice(channels[i])
	}

	return result
}

func CalculateMedianColor(img *image.RGBA) [3]uint8{
	var result = [3]uint8{}
	channels := ImageGetSplitChannels(img)

	for i := 0; i<3; i++ {
		result[i] = medianOfSlice(channels[i])
	}

	return result
}



// todo make median of slice and median of array generic on all number types and move to common-slices if necessary

// todo make median with implementation better than sort
func medianOfSlice(slice []uint8) uint8 {
	comparator :=  func(i, j int) bool {
		return slice[i] < slice[j]
	}
	sort.Slice(slice, comparator)

	length := len(slice)
 	bigger := slice[length / 2]    
    smaller:= slice[(length - 1 ) / 2]

    // calculates average of 2 without running into overflow
    result := (bigger - smaller) / 2 + smaller

	return result;
}

//https://stackoverflow.com/questions/1930454/what-is-a-good-solution-for-calculating-an-average-where-the-sum-of-all-values-e
func meanOfSlice(slice []uint8) uint8 {
	var average float64 = 0

	for index, value := range slice {
		f64val := float64(value)
		average += (f64val - average) / float64(index + 1)
	} 

	return uint8(math.Round(average))
}