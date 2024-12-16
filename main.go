package main

import (
	"fmt"
	"net/http"
	"time"
)

func handleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	_, err := fmt.Fprintf(w, "data: Hello, client! Time: %s\n\n", time.Now().Format(time.UnixDate))
	if err != nil {
		http.Error(w, "Failed to send initial message", http.StatusInternalServerError)
		return
	}
	flusher.Flush()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_, err := fmt.Fprintf(w, "data: The time is %s\n\n", time.Now().Format(time.UnixDate))
			if err != nil {
				fmt.Println("Error writing data:", err)
				return
			}
			flusher.Flush()
		case <-r.Context().Done():
			fmt.Println("Client disconnected")
			return
		}
	}
}

func main() {
	http.HandleFunc("/events", handleSSE)

	fmt.Println("Server started at http://localhost:9000")
	if err := http.ListenAndServe(":9000", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
