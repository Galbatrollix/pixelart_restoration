package kuwahara

import "math"

/*

	Computes flat gaussian kernel, implementation based on opencv
	https://docs.opencv.org/4.x/d4/d86/group__imgproc__filter.html#gac05a120c1ae92a6060dd0db190a61afa

*/
func gaussianKernel1D(ksize int, sigma float32) []float32 {
	if ksize < 1 {
		panic("Non positive ksize provided to gaussianKernel1D")
	}
	if sigma <= 0{
		panic("Non positive sigman provided to gaussianKernel1D")
	}

	var denominator = - (2.0 * sigma * sigma)

	var total_sum float32 = 0.0
	kernel := make([]float32, ksize)
	for i:=0 ; i<ksize ; i++{
		nominator_sqrt := (float32(i) - (float32(ksize) - 1.0)/2.0)
		nominator := nominator_sqrt * nominator_sqrt
		full_value := float32(math.Exp(float64(nominator / denominator)))
		total_sum += full_value
		kernel[i] = full_value
	}

	// dividing everything by total sum to get sum equal to 1
	for i:=0 ; i<ksize ; i++{
		kernel[i] /= total_sum
	}

	return kernel
}