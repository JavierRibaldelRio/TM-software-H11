package sensors

// Stores the data before to be sended
type Record struct {
	Timestamp string // Time
	Magnitude string
	Min       float64
	Max       float64
	Avg       float64
}
