package kuwahara



func sliceSumFloat32(slice []float32) float32{
    var sum float32 = 0.0
    for _, value := range slice{
        sum += value
    }
    return sum
}

func sliceDivbyFloat32(slice []float32, divisor float32){
    for id := range slice{
        slice[id] /= divisor
    }
}

func sliceGetSquared(slice []float32) []float32{
	result := make([]float32, len(slice))
	for id := range slice {
		result[id] = slice[id] * slice[id]
	}
	return result
}

func sliceSubtractSquared(slice []float32, other []float32){
	for id := range slice {
		slice[id] -= other[id] * other[id]
	}
}


func sliceFloat32ToUint8(floats []float32, uints []uint8){
	for id := range floats {
		uints[id] = uint8(floats[id] + 0.5)
	}
}

func sliceUint8ToFloat32(uints []uint8, floats []float32){
	for id := range floats {
		floats[id] = float32(uints[id])
	}
}