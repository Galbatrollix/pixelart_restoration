package types


/*

	Combined list holds three attributes: Intevals, IntervalTypes and totalCount
	Intervals:
		slice of N intervals of consecutive pixels without edges between them.
		Intervals slice can contain 0 - length items!
	IntervalTypes:
		slice of N uint8 items denoting type of interval item of the same index - Pixel, Gridline or Unknown
		definition of values of this slice can be seen below in enum-like constant group

	For example:
		Intervals: 		[4,1,5,13,4,1,5]
		IntervalTypes: 	[0,1,0,2 ,0,1,0]

		This means that there are 33 pixels in the image, split into 7 consecutive segments, which represents the following:
		4 pixels 'pixel' type -> 1 pixel 'gridline' type -> 5 pixels 'pixel' type, -> 13 pixels 'unknown' type (...)
*/
type CombinedList struct{
	Intervals []uint

}

/*
	
*/
const (  
        INTERVAL_PIXEL uint8 = iota   // 0
       	INTERVAL_GRID uint8 = iota    // 1
        INTERVAL_UNKNOWN uint8 = iota // 2
)

