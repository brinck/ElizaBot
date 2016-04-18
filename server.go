package elizabot

import (
    "fmt"
    "net/http"
    "encoding/json"
)

/* 
 * init()
 *
 * Supposedly this function is run by the google containers
 * rather than the main() function to setup the server 
 * bindings
 */
func init() {
    http.HandleFunc("/", handler)
    http.HandleFunc("/webhook/", eliza)
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello, world!")
}



/*
 * JSON webhook structs
 */
type Webhook struct {
    entry []Entry
}

type Entry struct {
    id string,
    time int,
    messaging []Messaging
}

type Messaging {
	sender Sender,
	recipient Recipient,
	timestamp int,
	message Message;
}

type Sender {
	id string
}

type Recipient {
	id string
}

type Message {
	mid string,
	seq int, 
	text string
}



/*
 * eliza(w http.ResponseWriter, r *http.Request) 
 *
 * Handler that lets the messenger API interface
 * with the elizabot.
 */
func eliza(wr http.ResponseWriter, req *http.Request) 
{
	// parse the request in json format
	var data Webhook;
	err := json.Unmarshal(req.Body, &data);
	
	if err != nil {

	} else {
		// loop through messages
		messagingEvents := data.entry[0].messaging;
		for i := 0; i < len(messagingEvents); i++ {
			event := messagingEvents[i]
			sender := event.sender.id
			if event.message
		}
	}	
}
