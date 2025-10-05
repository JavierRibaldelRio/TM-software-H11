package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"ribal-backend-receiver/csvutil"
	"ribal-backend-receiver/httpws"
	"ribal-backend-receiver/logger"
	"ribal-backend-receiver/ringbuffer"
	"ribal-backend-receiver/sensors"
	"ribal-backend-receiver/state"
)

func main() {

	// buffer to write into the csv
	writeBuffer := make(chan sensors.Record, 4096)

	// ring buffer
	ring := ringbuffer.NewRing[sensors.Record](10)

	// TCP connction with
	go acceptIncomeConn(writeBuffer)

	go recieveData(writeBuffer, ring)

	// Start http server
	httpws.StartHttpWSServer(ring)

}

// Accept new tpconnections and register the data into de buffer after each msg
func acceptIncomeConn(buffer chan<- sensors.Record) {

	// open server
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error:", err)
		logger.Error(err.Error())

		return
	}
	defer listener.Close()

	logger.Info("TCP server listening on port 8080 (sensors)")

	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			logger.Error(err.Error())
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
		// Parses data
		var rec sensors.Record
		if err := dec.Decode(&rec); err != nil {
			if err == io.EOF {
				logger.Info("Client disconnected to TCP 8080 server")
				return
			}
			fmt.Println("decode error:", err)
			logger.Error("decode error" + err.Error())
			return
		}

		// Adds it to chanel
		buffer <- rec
	}
}

// Recieve the data from sensors, writes it down on the CSV, sends it through webscokets and adds it to the round buffer
func recieveData(buffer chan sensors.Record, ring *ringbuffer.RingBuffer[sensors.Record]) {

	for rec := range buffer {

		// If it is power on avoids scripture
		if !state.IsPowerOn() {
			continue
		}

		// CSV
		csvutil.AddToCSV(rec)

		// Broadcast
		httpws.BroadcastJSON(rec)

		// Ring buffer
		ring.Add(rec)
	}
}
