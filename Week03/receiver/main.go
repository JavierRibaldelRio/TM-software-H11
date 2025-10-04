package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"ribal-backend-receiver/csvutil"
	"ribal-backend-receiver/httpws"
	"ribal-backend-receiver/sensors"
)

func main() {

	// buffer to write into the csv
	writeBuffer := make(chan sensors.Record, 4096)

	// TCP connction with backends
	go acceptIncomeConn(writeBuffer)

	go writeBufferToCSV(writeBuffer)

	httpws.StartHttpWSServer()

}

// Accept new tpconnections and register the data into de buffer after each msg
func acceptIncomeConn(buffer chan<- sensors.Record) {

	// open server
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

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

		buffer <- rec
	}
}

// writes buffer into the csv
func writeBufferToCSV(buffer chan sensors.Record) {
	// Open csv
	// CSV writer

	writer, closeFile := csvutil.SetUpCSVWriter()
	defer closeFile()
	for rec := range buffer {

		csvutil.AddToCSV(*writer, rec)

		httpws.BroadcastJSON(rec)

	}
}
