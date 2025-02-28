package common


// MemCopy: this function will copy all elements from source slice to destination slice. 
// Will panic if source slice is larger than destination
// UB if memory of destination and source overlaps
func MemCopy[T any](to []T, from []T){
    for i := range from {
        to[i] = from[i]
    }

}

// MemSet: this function will set all elements of target slice to given value.
func MemSet[T any](slice []T , val T){
    for i := range slice {
        slice[i] = val
    }
}

// MemRepeat: this function will repeat elements from sequence over and over an write them to target 
// until target's length is reached. Will panic if provided sequence slice doesn't contain any elements
// UB if inputs overlap
func MemRepeat[T any](target []T, sequence []T){
    for i := range target {
        sequence_id := i % len(sequence)
        target[i] = sequence[sequence_id]
    }
}

// MemSwap: this function will set swap all elements of left slice with all elements of right slice. 
// Will panic on call with different lengthed slices
// UB if inputs overlap
func MemSwap[T any](left []T, right []T){
    for i := range left {
        left[i], right[i] = right[i], left[i]
    }
}


// MemReverse: this function will reverse ordering of all elements in the given slice
func MemReverse[T any](slice []T) {
    count := len(slice)

    for i := 0 ; i < count/2 ;i++{
        slice[i], slice[count - 1 - i] = slice[count - 1 - i], slice[i]
    }

}

// MemReverseExtended: this function will reverse ordering of all elements in batch_sized batches
// Will panic if length of slice is not divisible by batch_size, and if not (0 < batch_size <= len(slice))
func MemReverseExtended[T any](slice []T, batch_size int){
    if len(slice) % batch_size != 0 {
        panic("Slice length not divisible by batch size")
    }
    if batch_size > len(slice) || batch_size < 1{
        panic("Batch size not in range. Required: (0 < batch_size <= len(slice))")
    }

    count_total := len(slice)
    count_items := count_total / batch_size
    for i := 0 ; i < count_items / 2; i++ {
        id_absolute_front := i * batch_size
        id_absolute_back := count_total - batch_size - id_absolute_front
        MemSwap(
            slice[id_absolute_front: id_absolute_front + batch_size],
            slice[id_absolute_back: id_absolute_back + batch_size],
        )
    }


}