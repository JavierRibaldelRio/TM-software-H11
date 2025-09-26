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
	sensors.StartSensors(ctx, readings, "sensors.json")

	// CSV writer
	writer, closeFile := SetUpCSVWriter()
	defer closeFile()

	// For each second
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// Register the output of the sensors
	runProcessingLoop(ctx, readings, data, ticker, writer)

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

// setCSVWriter configures and return the csv file writer and a function
// to close the writer
func SetUpCSVWriter() (*csv.Writer, func()) {

	// Create file
	filename := fmt.Sprintf("output_%s.csv", time.Now().Format("20060102_150405"))

	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	// Create writer
	writer := csv.NewWriter(f)

	// write header in log
	writer.Write([]string{"timestamp", "magnitude", "avg", "min", "max"})

	// Retunrs thew writer and a function to close it
	return writer, func() {
		writer.Flush()
		f.Close()
	}
}

// runProcessingLoops recives the output of the sensors, stores it and obtains the main stats
func runProcessingLoop(ctx context.Context, readings chan sensors.Data, data map[string][]float64, ticker *time.Ticker, writer *csv.Writer) {
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
