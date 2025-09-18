package main

import (
	"context"
	"ribal-backend/sensors"
	"time"
)

func main() {

	// Context ensures that when the program is closed
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// channel of reading
	readings := make(chan sensors.Data, 4096)

	// Stores all the data
	batch := make([]sensors.Data, 0, 4096)

	// Starts all the sensors
	StartSensors(ctx, readings)

	// For each second
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {

		select {

		// Adds the reading to de batch
		case d := <-readings:
			batch = append(batch, d)

		case <-ticker.C:
			if len(batch) > 0 {

				//TODO: stats and save data

				// Empties the batch
				batch = batch[:0]
			}

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
		{Name: "voltaje", Number: 1, Unit: "V", Period: 50 * time.Millisecond, Avg: 12, Sdev: 0.25},
		{Name: "voltaje", Number: 2, Unit: "V", Period: 50 * time.Millisecond, Avg: 12, Sdev: 0.50},
		{Name: "voltaje", Number: 3, Unit: "V", Period: 55 * time.Millisecond, Avg: 12, Sdev: 1.00},
		{Name: "voltaje", Number: 4, Unit: "V", Period: 60 * time.Millisecond, Avg: 12, Sdev: 2.00},
		{Name: "distance", Number: 1, Unit: "cm", Period: 20 * time.Millisecond, Avg: 12, Sdev: 0.15},
		{Name: "distance", Number: 2, Unit: "cm", Period: 23 * time.Millisecond, Avg: 12, Sdev: 0.50},
		{Name: "distance", Number: 3, Unit: "cm", Period: 23 * time.Millisecond, Avg: 12, Sdev: 0.50},
	}

	// Launch all sensors.
	for i := range ss {
		ss[i].StartMeasuring(ctx, readings)
	}
}
