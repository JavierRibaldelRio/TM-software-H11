package sensors

import (
	"context"
	"math/rand"
	"time"
)

type Sensor struct {
	Name   string
	Number int64
	Unit   string
	Period time.Duration
	Avg    float64
	Sdev   float64
}

type Data struct {
	Measure   float64
	Unit      string
	Timestamp time.Time
	Name      string
	Number    int64
}

// Using the properties of a Sensor generates following a normal distribution values
func (s Sensor) Read() Data {
	v := s.Avg + s.Sdev*rand.NormFloat64() // N(Avg, Sdev)
	return Data{
		Timestamp: time.Now(),
		Measure:   v,
		Name:      s.Name,
		Number:    s.Number,
		Unit:      s.Unit,
	}
}

// StartMeasuring starts a goroutine that periodically sends readings to out.
func (s Sensor) StartMeasuring(ctx context.Context, out chan<- Data) {

	// Starts the gororutine, uses anonymous function
	go func() {

		ticker := time.NewTicker(s.Period) // Creates a new ticker for each period
		defer ticker.Stop()                // Ensures that the ticker is properly removed

		// Un
		for {

			// Wait until a case ocurrs
			select {

			// Tick
			case <-ticker.C:
				out <- s.Read()

			// Stop app
			case <-ctx.Done():
				return

			}

		}

	}()
}
