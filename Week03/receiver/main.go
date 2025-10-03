package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"ribal-backend-receiver/csvutil"
	"ribal-backend-receiver/sensors"
)

func main() {

	// Open csv
	// CSV writer
	writer, closeFile := csvutil.SetUpCSVWriter()
	defer closeFile()

	// open server
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	// Inform use that port is ready
	fmt.Println("Server is listening on port 8080")

	// buffer to write into the csv
	writeBuffer := make(chan sensors.Record, 4096)

	go writeBufferToCSV(writeBuffer, writer)

	acceptIncomeConn(listener, writeBuffer)

}

func acceptIncomeConn(listener net.Listener, buffer chan<- sensors.Record) {
	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Handle client connection in a goroutine
		go handleClient(conn, buffer)
	}

}

// Handles client coneccion
func handleClient(conn net.Conn, buffer chan<- sensors.Record) {
	defer conn.Close()

	dec := json.NewDecoder(conn)

	for {
		var rec sensors.Record
		if err := dec.Decode(&rec); err != nil {
			if err == io.EOF {
				fmt.Println("client disconnected")
				return
			}
			fmt.Println("decode error:", err)
			return
		}

		// Política simple: bloquear hasta que haya hueco.
		buffer <- rec
	}
}

// writes buffer into the csv
func writeBufferToCSV(buffer chan sensors.Record, writer *csv.Writer) {

	for rec := range buffer {

		csvutil.AddToCSV(*writer, rec)

	}
}
