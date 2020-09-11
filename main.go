package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Channels
type globalChannel struct {
	Recepient chan string
	Sender    chan string
}

// POST
type postResponse struct {
	Author string
}

// Game is ...
type Game struct {
	LastAuthor           string
	registeredMessage    string
	LastAuthorMessage    string
	AskLastAuthorMessage string
}

// CreateGame is ...
func CreateGame(LastAuthor string, registeredMessage string, LastAuthorMessage string, AskLastAuthorMessage string) Game {
	return Game{LastAuthor, registeredMessage, LastAuthorMessage, AskLastAuthorMessage}
}

func handleGame(gChannel *globalChannel, currentGame *Game) {
	for {
		// Storing author
		author := <-gChannel.Recepient
		// Settings last author
		currentGame.LastAuthor = author

		// If no one is last auhor
		if currentGame.LastAuthor == "Aucun author" {
			author = "personne !"
		}

		// Formatting string for concatenation
		var buffer bytes.Buffer
		buffer.WriteString(currentGame.LastAuthorMessage)
		buffer.WriteString(author)

		// Print texr
		fmt.Println(<-gChannel.Sender)
		fmt.Println(buffer.String())
	}
}

// POST
func handleRequest(postResponse *postResponse, gChannel *globalChannel, currentGame *Game) {
	// Sending message into channel
	gChannel.Recepient <- postResponse.Author
	gChannel.Sender <- currentGame.registeredMessage
}

// GET
func getAuthor(gChannel *globalChannel, currentGame *Game) {
	// Sending message into channel
	gChannel.Recepient <- currentGame.LastAuthor
	gChannel.Sender <- currentGame.AskLastAuthorMessage
}

func generateHandler(gChannel *globalChannel, currentGame *Game) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		// Settings headers
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {

		case "GET":
			// Go routines which use new threads for each get request
			go getAuthor(gChannel, currentGame)
			// Sending back to client header and response
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(currentGame.LastAuthor))

		case "POST":
			// Structure declaration
			var postResponse postResponse
			// Decode body -> store body in structure & store errors if exist
			err := json.NewDecoder(r.Body).Decode(&postResponse)
			if err != nil {
				// Sending bad request status if exists
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			// Go routines which use new threads for each post request
			go handleRequest(&postResponse, gChannel, currentGame)
			// Sending back to client header and response
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(postResponse.Author))
		}
	}
}

func main() {
	// Creating global channel struct composed with 2 chans
	gChannel := globalChannel{make(chan string), make(chan string)}
	// Creating game struct
	currentGame := CreateGame("Aucun author", "Votre valeur a bien été enregistrée", "Le dernier auteur est ", "Vous avez demandé le dernier auteur")

	// Go routine which is chans listenner
	go handleGame(&gChannel, &currentGame)

	// Route wich listen on GET & POST
	http.HandleFunc("/game", generateHandler(&gChannel, &currentGame))
	// HTTP serve creation & Logging stuff
	log.Fatal(http.ListenAndServe(":8080", nil))
}
