package main

import "testing"

const lastAuthor = "Aucun author"
const registeredMessage = "Votre valeur a bien été enregistrée"
const lastAuthorMessage = "Le dernier auteur est "
const askLastAuthorMessage = "Vous avez demandé le dernier auteur"

func TestGameCreation(t *testing.T) {
	testGame := CreateGame(lastAuthor, registeredMessage, lastAuthorMessage, askLastAuthorMessage)
	if testGame.LastAuthor != lastAuthor || testGame.registeredMessage != registeredMessage || testGame.LastAuthorMessage != lastAuthorMessage || testGame.AskLastAuthorMessage != askLastAuthorMessage {
		t.Error("Les valuers passées en paramètres sont erronées")
	}
}

// func TestGchan(t *testing.T) {
// 	gChannel := globalChannel{make(chan string), make(chan string)}
// 	gGame := Game{lastAuthor, registeredMessage, lastAuthorMessage, askLastAuthorMessage}
// 	noAuthor := "Aucun author"
// 	gChannel.Recepient <- noAuthor
// 	handleGame(&gChannel, &gGame)
// }
