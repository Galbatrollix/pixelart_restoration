package contrast

/*


*/
func SelectMostFrequent(edge_counts []uint, params MostFrequentParams) []uint {
	return nil
}

/*
	Clip top: float32
		(After filtering out positions with 0 edges)
		removes <ClipTop * 100%> most common edge positions 
	Clip bottom: float32
		(After filtering out positions with 0 edges)
		removes <ClipBottom * 100%> least common (but non-0) edge position
	CutoffTreshold: float32
		(After applying both clip top and clip bottom)
		The lowest permissable edge count for a position is set to CutoffTreshold * 100% *<most common edge position count>)

	Constraints:
	0 <= ClipTop < 1.0
	0 <= ClipBottom < 1.0
	ClipTop + ClipBottom < 1.0
	CutoffTreshold >= 0
	
*/
type MostFrequentParams struct {
	ClipTop float32
	ClipBottom float32
	CutoffTreshold float32
}


func GetBaseMostFreuentParams() MostFrequentParams {
	return MostFrequentParams{
		ClipTop: 0.05,
		ClipBottom: 0.0,
		CutoffTreshold: 0.3,
	}
}