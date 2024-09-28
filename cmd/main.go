package main

import (
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/hypebeast/go-osc/osc"
)

var (
	addr          = "239.255.0.1:9000" // Multicast address for OSC
	websocketAddr = ":8080"            // WebSocket server address
	clients       = make(map[*websocket.Conn]bool)
	broadcast     = make(chan *osc.Message)
	upgrader      = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	mux sync.Mutex
)

func main() {
	// Start the OSC listener
	go oscListener()

	// Start the WebSocket server
	http.HandleFunc("/ws", handleConnections)
	go func() {
		log.Printf("WebSocket server listening on %s\n", websocketAddr)
		if err := http.ListenAndServe(websocketAddr, nil); err != nil {
			log.Fatalf("WebSocket server error: %v", err)
		}
	}()

	// Start broadcasting OSC messages to WebSocket clients
	for {
		message := <-broadcast
		sendToWebsocketClients(message)
	}
}

// oscListener listens for multicast OSC messages and forwards them to WebSocket clients
func oscListener() {
	addr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %v", err)
	}

	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("Failed to join multicast group: %v", err)
	}
	defer conn.Close()

	log.Printf("Listening for OSC messages on %s\n", addr.String())

	buf := make([]byte, 1024)
	for {
		n, _, err := conn.ReadFrom(buf)
		if err != nil {
			log.Printf("Error receiving OSC message: %v", err)
			continue
		}

		packet, err := osc.ParsePacket(string(buf[:n]))
		if err != nil {
			log.Printf("Error parsing OSC packet: %v", err)
			continue
		}

		switch p := packet.(type) {
		case *osc.Message:
			log.Printf("Received OSC message: %v", p)
			broadcast <- p // Send OSC message to WebSocket clients
		}
	}
}

// handleConnections upgrades HTTP connections to WebSockets and registers clients
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// allow any origin for now
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	mux.Lock()
	clients[conn] = true
	mux.Unlock()

	// Read messages from WebSocket and send them to OSC
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading WebSocket message: %v", err)
			mux.Lock()
			delete(clients, conn)
			mux.Unlock()
			break
		}

		log.Printf("Received WebSocket message: %s", msg)
		sendToOSC(string(msg))
	}
}

// sendToWebsocketClients sends the received OSC message to all WebSocket clients
func sendToWebsocketClients(message *osc.Message) {
	mux.Lock()
	defer mux.Unlock()

	for client := range clients {
		err := client.WriteJSON(message)
		if err != nil {
			log.Printf("Error sending message to WebSocket client: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}

// sendToOSC sends messages received from WebSocket clients back to the OSC multicast group
func sendToOSC(msg string) {
	client := osc.NewClient("239.255.0.1", 9000)
	oscMsg := osc.NewMessage("/websocket")
	oscMsg.Append(msg)
	err := client.Send(oscMsg)
	if err != nil {
		log.Printf("Error sending OSC message: %v", err)
	} else {
		log.Printf("Sent OSC message: %s", msg)
	}
}
