package sensors

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

type Sensor struct {
	Name   string
	Number int64
	Unit   string
	Period time.Duration // on ms
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

		ticker := time.NewTicker(s.Period * time.Millisecond) // Creates a new ticker for each period
		defer ticker.Stop()                                   // Ensures that the ticker is properly removed

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

// StartSensors starts sensor goroutines that write to 'readings'.
// needs the ctx, reading chan and the configuration path
// They stop when ctx is canceled. This function does not close 'readings'.
func StartSensors(ctx context.Context, readings chan<- Data, jsonPath string) {

	// Open json file
	jsonFile, err := os.Open(jsonPath)

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	// Reads the file
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Parses the json as an slice of Sensors
	var ss []Sensor
	json.Unmarshal([]byte(byteValue), &ss)

	// Activates each sensor
	for i := range ss {
		ss[i].StartMeasuring(ctx, readings)
	}

}
