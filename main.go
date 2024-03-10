/*
TOKI

Service for WikiForum
*/
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

  "github.com/rs/cors"
)

// Message struct
type Message struct {
	Name  string `json:"name"`
	Value string `json:"message"`
	Time  string `json:"date"`
}

// Function to get last ten messages from file
func getLastTenMessages() []Message {
	var messages []Message

	// Open file
	file, err := os.Open("all.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Scan all messages to array
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		messageJSON := scanner.Text()
		var message Message
		err := json.Unmarshal([]byte(messageJSON), &message)
		if err != nil {
			log.Fatal(err)
		}
		messages = append(messages, message)
	}

	// Get last 10 messages from array
	lastTenMessages := make([]Message, 0, 100)
	j := len(messages) - 1
	for i := 1; i <= 100; i++ {
		if j >= 0 {
			lastTenMessages = append(lastTenMessages, messages[j])
			j--
		}
	}
	return lastTenMessages
}

func createMessage(w http.ResponseWriter, r *http.Request) {
	// Fields
	nameForm := r.FormValue("name")
	messageForm := r.FormValue("message")
	date := time.Now().Format("2006-01-02")

	// Name and Message can't be blank
	if nameForm == "" || messageForm == "" {
		w.Write([]byte("Name and/or Message can't be blank"))
		return
	}

	message := Message{nameForm, messageForm, date}

	// Convert message to JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}

	// Create file
	file, err := os.OpenFile("all.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// It should be easier, but I'm really stupid
	messageText := string(messageJSON) + "\n"

	// Write to file
	_, err = file.Write([]byte(messageText))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(messageText)

	w.Write([]byte(messageJSON))
}

func main() {
	// Create router
	mux := http.NewServeMux()

	// Bind /get to function
	mux.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		lastTenMessages := getLastTenMessages()
		json.NewEncoder(w).Encode(lastTenMessages)
	})

	// Bind /create to function
	mux.HandleFunc("/create", createMessage)

	// Bind / to show text
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// For CORS policy. If not do, we can't do requests from JS
	handler := cors.Default().Handler(mux)

	// Running server
	fmt.Println("Server running on 8080")
	http.ListenAndServe(":8080", handler)
}
