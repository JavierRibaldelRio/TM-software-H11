package httpws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ribal-backend-receiver/sensors"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WS updatesr
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Set of clients
var (
	clients   = make(map[*websocket.Conn]struct{}) // SET
	clientsMu sync.RWMutex                         // RWMutex used to avoid reading and writing at the same time
)

func StartHttpWSServer() {

	mux := http.NewServeMux()

	mux.HandleFunc("/api/stream", wsStreamHandler)

	addr := ":8081"
	fmt.Println("Servidor HTTP escuchando en", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		panic(err)
	}

}

/**
* Web Sockets
 */

// Conection to /api/stream
func wsStreamHandler(w http.ResponseWriter, r *http.Request) {

	// Tries to upgrade connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("upgrade error: %v", err)
		return
	}
	addClient(conn) // alta

	// Prints the number of clients
	fmt.Printf("WS CONNECT (%d activos)\n", func() int { clientsMu.RLock(); defer clientsMu.RUnlock(); return len(clients) }())

	// Discard readings
	go func(c *websocket.Conn) {
		defer removeClient(c)
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				fmt.Printf("WS CLOSE: %v\n", err)
				return
			}
		}
	}(conn)
}

// Adds a client to set of clients
func addClient(c *websocket.Conn) {

	clientsMu.Lock()        // Blocks new access to read clients
	clients[c] = struct{}{} // Adds the key to the Set
	clientsMu.Unlock()      // Unblocks the access
}

// Remove a client from the set of clients
func removeClient(c *websocket.Conn) {
	clientsMu.Lock()   // Blocks acces
	delete(clients, c) // Delete the client from the server
	clientsMu.Unlock() // Unblocks the st
	_ = c.Close()      // Closes the connection with that specific client
}

func BroadcastJSON(dataSensor sensors.Record) {

	// Transforms it to a JSON
	jsonData, err := json.Marshal(dataSensor)
	if err != nil {
		return
	}

	// Locks the client to read it
	clientsMu.RLock()
	defer clientsMu.RUnlock()

	for c := range clients {
		// Limit to write
		_ = c.SetWriteDeadline(time.Now().Add(10 * time.Second))

		// writes message
		if err := c.WriteMessage(websocket.TextMessage, jsonData); err != nil {
			// if the brodcast fails shutdown the connection
			go removeClient(c)
		}
	}
}
