package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
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

	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Handle client connection in a goroutine
		go handleClient(conn, writer)
	}
}

func handleClient(conn net.Conn, writer *csv.Writer) {
	defer conn.Close()

	for {
		// Read data from the client
		d := json.NewDecoder(conn)

		var rec sensors.Record
		err := d.Decode(&rec)
		if err != nil {
			fmt.Println("error:", err)
		}

		// Process and use the data (here, we'll just print it)
		csvutil.AddToCSV(*writer, rec)

	}
}
