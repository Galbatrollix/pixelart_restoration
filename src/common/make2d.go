package common

// // Make2D will create a slice of slices of T, where all inner slices share common memory which can be accessed flat
// // To access flat the first inner slice must be resliced to size of (height * width)
// func Make2D[T any](height, width int) [][]T{
//     buffer := make([]T, height*width)
//     return Make2DFromFlat(buffer, height, width)
// }

// // Make2DFromFlat will create a slice of slices of T, where all inner slices share common memory which can be accessed flat
// // The flat slice is to make 2d slice from is provided as parameter
// func Make2DFromFlat[T any](buffer []T, height, width int) [][]T {
// 	vector := make([][]T, height)
//     for i := range vector {
//         start := i * width
//         end := (i+1) * width
//         vector[i] = buffer[start:end]
//     }
//     return vector
// }