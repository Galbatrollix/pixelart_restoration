package contrast

import "image"

/*
	Given binary (thresholded) grayscale image of edges detected,
	create a slice where slice[x] = count of edges detected at position x.
*/

func EdgesToEdgeCounts(edges_binary *image.Gray) []uint {
	height, width := edges_binary.Rect.Dy(), edges_binary.Rect.Dx()

	edge_positions_count := width - 1
	edge_counts := make([]uint, edge_positions_count)

	for y:=0; y<height; y++{
		for x:=0 ;x<edge_positions_count; x++{
			flat_id := edges_binary.PixOffset(x + edges_binary.Rect.Min.X, y + edges_binary.Rect.Min.Y)
			if edges_binary.Pix[flat_id] != 0 {
				edge_counts[x] += 1
			}
		}
	}
	return edge_counts

}