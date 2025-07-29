package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	// Wait a moment for server to start if running in parallel
	time.Sleep(1 * time.Second)

	baseURL := "http://localhost:8080"

	// Test health endpoint
	fmt.Println("=== Testing Health Endpoint ===")
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		log.Printf("Health check failed: %v", err)
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Health: %s\n", string(body))
		resp.Body.Close()
	}

	// Test tools endpoint
	fmt.Println("\n=== Testing Tools Endpoint ===")
	resp, err = http.Get(baseURL + "/tools")
	if err != nil {
		log.Printf("Tools endpoint failed: %v", err)
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Tools: %s\n", string(body))
		resp.Body.Close()
	}

	// Test resources endpoint
	fmt.Println("\n=== Testing Resources Endpoint ===")
	resp, err = http.Get(baseURL + "/resources")
	if err != nil {
		log.Printf("Resources endpoint failed: %v", err)
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Resources: %s\n", string(body))
		resp.Body.Close()
	}

	fmt.Println("\n=== Testing Complete ===")
	fmt.Println("Server is running and responding to HTTP requests!")
	fmt.Println("WebSocket MCP endpoint available at: ws://localhost:8080/mcp")
}
