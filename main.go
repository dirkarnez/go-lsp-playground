package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// LSP message structure
type Message struct {
	Jsonrpc string          `json:"jsonrpc"`
	ID      *int            `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// Response structure
type Response struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func main() {
	stdin := os.Stdin
	stdout := os.Stdout

	buf := make([]byte, 4096) // Buffer to read LSP messages
	for {
		n, err := stdin.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			continue
		}

		// Parse LSP message
		var msg Message
		err = json.Unmarshal(buf[:n], &msg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid JSON: %v\n", err)
			continue
		}

		// Handle LSP methods
		switch msg.Method {
		case "initialize":
			handleInitialize(stdout, msg)
		case "shutdown":
			handleShutdown(stdout, msg)
			return // Exit after shutdown
		default:
			fmt.Fprintf(os.Stderr, "Unhandled method: %s\n", msg.Method)
		}
	}
}

func handleInitialize(stdout *os.File, msg Message) {
	response := Response{
		Jsonrpc: "2.0",
		ID:      *msg.ID,
		Result: map[string]interface{}{
			"capabilities": map[string]interface{}{}, // No capabilities supported in this minimal example
		},
	}
	sendResponse(stdout, response)
}

func handleShutdown(stdout *os.File, msg Message) {
	response := Response{
		Jsonrpc: "2.0",
		ID:      *msg.ID,
		Result:  nil,
	}
	sendResponse(stdout, response)
}

func sendResponse(stdout *os.File, response Response) {
	respBytes, err := json.Marshal(response)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling response: %v\n", err)
		return
	}

	// Write response header and body
	fmt.Fprintf(stdout, "Content-Length: %d\r\n\r\n", len(respBytes))
	stdout.Write(respBytes)
}
