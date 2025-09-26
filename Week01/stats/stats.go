package stats

// CalculateStats computes the average, minimum, and maximum values
// from a slice of float64 numbers.
func CalculateStats(data []float64) (avg, min, max float64) {

	sum := 0.0
	max = data[0]
	min = data[0]

	for _, d := range data {
		sum += d

		if d > max {
			max = d
		}

		if d < min {
			min = d
		}
	}

	avg = sum / float64(len(data))

	// Return average, minimum, and maximum
	return avg, min, max
}
