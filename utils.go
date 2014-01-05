package speedtest

import (
	"math"
	"sort"
)

const (
	earthRadius = 6371
	degToRad    = math.Pi / 180.0
)

// MaximalSumWindow finds the the largest sum of sequential values for the given
// window size.
//
// e.g.
// f([1, 3, 5, 2, 4], 1) = 5 as 5 is the largest single value in the list.
// f([1, 3, 5, 2, 4], 2) = 8 as (3 + 5) is 8, and no two other consecutive
// numbers can produce a larger sum.
//
// This function is used to get a reasonable best estimate of a sustained
// maximum within a set of data.
func MaximalSumWindow(data []int, size int) int {
	// Reduce window size to that of the input data
	if size > len(data) {
		size = len(data)
	}

	best := 0
	for i := 0; i < size; i++ {
		best += data[i]
	}

	curr := best
	for i := 0; i < len(data)-size; i++ {
		curr = curr - data[i] + data[i+size]
		if curr > best {
			best = curr
		}
	}

	return best
}

// MedianSumWindow calculates a median sum of the given window size within
// the given data.
func MedianSumWindow(data []int, size int) int {
	sorted := make([]int, len(data))
	copy(sorted, data)
	sort.Ints(sorted)
	sum := 0
	for i := 0; i < size; i++ {
		pos := (len(data)-size)/2 + i
		sum += sorted[pos]
	}
	return sum
}

// Distance calculates the distance between two geographical positions
func Distance(lat1 float64, lon1 float64, lat2 float64, lon2 float64) float64 {
	return earthRadius * math.Acos(
		math.Sin(lat1)*math.Sin(lat2)+
			math.Cos(lat1)*math.Cos(lat2)*math.Cos(lon2-lon1))
}
