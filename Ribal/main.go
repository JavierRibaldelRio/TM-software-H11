package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"ribal-backend/sensors"
	"ribal-backend/stats"
	"time"
)

func main() {

	// Context ensures that when the program is closed
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// channel of reading
	readings := make(chan sensors.Data, 4096)

	// Stores all the data of the sensors
	data := make(map[string][]float64)

	// Starts all the sensors
	StartSensors(ctx, readings)

	// Create file
	filename := fmt.Sprintf("output_%s.csv", time.Now().Format("20060102_150405"))

	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	// Create writer
	writer := csv.NewWriter(f)
	defer writer.Flush()

	// write header in log
	writer.Write([]string{"timestamp", "magnitude", "avg", "min", "max"})

	// For each second
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// Forever
	for {

		select {

		// Adds the reading to de batch
		case d := <-readings:
			data[d.Name] = append(data[d.Name], d.Measure)

		// Each seconds writes the log
		case t := <-ticker.C:
			ParseDataWriteLog(data, t, writer)

		// Stop app
		case <-ctx.Done():
			return

		}

	}

}

// StartSensors starts sensor goroutines that write to 'readings'.
// They stop when ctx is canceled. This function does not close 'readings'.
func StartSensors(ctx context.Context, readings chan<- sensors.Data) {
	ss := []sensors.Sensor{
		{Name: "voltaje", Number: 1, Unit: "V", Period: 50 * time.Millisecond, Avg: 12, Sdev: 1},
		{Name: "voltaje", Number: 2, Unit: "V", Period: 50 * time.Millisecond, Avg: 12, Sdev: 2},
		{Name: "voltaje", Number: 3, Unit: "V", Period: 55 * time.Millisecond, Avg: 12, Sdev: 2.5},
		{Name: "voltaje", Number: 4, Unit: "V", Period: 60 * time.Millisecond, Avg: 12, Sdev: 3},

		{Name: "distance", Number: 1, Unit: "cm", Period: 20 * time.Millisecond, Avg: 20, Sdev: 1},
		{Name: "distance", Number: 2, Unit: "cm", Period: 23 * time.Millisecond, Avg: 20, Sdev: 0.50},
		{Name: "distance", Number: 3, Unit: "cm", Period: 23 * time.Millisecond, Avg: 20, Sdev: 2},
	}

	// Launch all sensors.
	for i := range ss {
		ss[i].StartMeasuring(ctx, readings)
	}
}

// ParseDataWriteLog writes one CSV row per magnitude in `data`.
// For each magnitude with values, it computes avg/min/max,
// writes [timestamp, magnitude, avg, min, max] to `writer`,
// and clears the slice for the next batch.
func ParseDataWriteLog(data map[string][]float64, t time.Time, writer *csv.Writer) {

	// For each magnitude
	for magnitude, values := range data {

		// Avoid error in case there is not new data
		if len(values) == 0 {
			continue
		}

		// Calculates data
		avg, min, max := stats.CalculateStats(values)

		// Format new line
		record := []string{
			t.Format(time.RFC3339),
			magnitude,
			fmt.Sprintf("%f", avg),
			fmt.Sprintf("%f", min),
			fmt.Sprintf("%f", max),
		}

		// Add new line
		writer.Write(record)
		writer.Flush()

		// Remove content
		data[magnitude] = data[magnitude][:0]
	}
}
