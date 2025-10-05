package csvutil

import (
	"encoding/csv"
	"fmt"
	"os"
	"ribal-backend-receiver/logger"
	"ribal-backend-receiver/sensors"
	"strconv"
	"time"
)

var (
	writer   *csv.Writer
	file     *os.File
	filename string
)

// creates the loger file
func init() {

	// Create file
	filename = fmt.Sprintf("logs/data/data-output-%s.csv", time.Now().Format("2006-01-02_15-04-05"))

	createCSV()

	logger.Info("CSV file was created at " + filename)

}

// creates the csv file
func createCSV() {
	f, err := os.Create(filename)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	file = f

	// Create writer
	writer = csv.NewWriter(f)

	// write header in log
	writer.Write([]string{"timestamp", "magnitude", "avg", "min", "max"})
}

// adds a record to the csv
func AddToCSV(data sensors.Record) {

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

// Overwrites (cleans) the current CSV file, preserving the same filename
func ClearCSV() {
	if file != nil {
		file.Close()
	}
	createCSV()

	logger.Info("CSV file was cleared out")
}
