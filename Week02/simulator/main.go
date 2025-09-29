package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
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

	// For each second
	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()

	// TCP connection
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	// Register the output of the sensors
	runProcessingLoop(ctx, readings, data, ticker, conn)

}

// SendData sends de data to the server
// For each magnitude with values, it computes avg/min/max,
// sends [timestamp, magnitude, avg, min, max] to the backend,
// and clears the slice for the next batch.
func SendData(data map[string][]float64, t time.Time, conn net.Conn) {

	// For each magnitude
	for magnitude, values := range data {

		// Avoid error in case there is not new data
		if len(values) == 0 {
			continue
		}

		// Calculates data
		avg, min, max := stats.CalculateStats(values)

		// Format new line
		record := sensors.Record{Timestamp: t.Format(time.RFC3339), Magnitude: magnitude, Min: min, Max: max, Avg: avg}

		// Transform to json
		msg, err := json.Marshal(record)

		if err != nil {
			fmt.Println("error:", err)
		}

		// send to server
		conn.Write(msg)

		// Remove content
		data[magnitude] = data[magnitude][:0]
	}
}

// runProcessingLoops recives the output of the sensors, stores it and obtains the main stats
func runProcessingLoop(ctx context.Context, readings chan sensors.Data, data map[string][]float64, ticker *time.Ticker, conn net.Conn) {
	// Forever
	for {

		select {

		// Adds the reading to de batch
		case d := <-readings:
			data[d.Name] = append(data[d.Name], d.Measure)

		// Each seconds writes the log
		case t := <-ticker.C:
			SendData(data, t, conn)

		// Stop app
		case <-ctx.Done():
			return

		}

	}
}
