package types


/*

	Interval list holds two attributes: data and totalCount
	Intervals:
		slice of N intervals of consecutive pixels without edges between them.
	TotalCount:
		sum of all interval lengths in data slice - equal to dimension of the origin image in pixels

	For example:
		Intervals: [4,3,5,4], TotalCount: 16
		This means that there are 16 pixels in the image, split into 4 consecutive segments, which represents the following:
		4 pixels > edge > 3 pixels > edge > 5 pixels > edge 4 pixels
*/

type IntervalList struct{
	Intervals []uint
	TotalCount uint
}

