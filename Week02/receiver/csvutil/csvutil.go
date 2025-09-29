package csvutil

import (
	"encoding/csv"
	"fmt"
	"os"
	"ribal-backend-receiver/sensors"
	"strconv"
	"time"
)

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

// given a records an a csv adds it to the csv

func AddToCSV(writer csv.Writer, data sensors.Record) {

	// Format new line
	record := []string{
		data.Timestamp,
		data.Magnitude,
		strconv.FormatFloat(data.Avg, 'f', -1, 64),
		strconv.FormatFloat(data.Min, 'f', -1, 64),
		strconv.FormatFloat(data.Max, 'f', -1, 64),
	}

	// Add new line
	writer.Write(record)
	writer.Flush()

}
