package elizabot

import (
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    "github.com/kennysong/goeliza"
)

/* 
 * init()
 *
 * Supposedly this function is run by the google containers
 * rather than the main() function to setup the server 
 * bindings
 */
func init() {
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/webhook/", webhookHandler)
}

/*
 * homeHandler(wr http.ResponseWriter, req *http.Request)
 *
 * Renders the page at "/"
 */
func homeHandler(wr http.ResponseWriter, req *http.Request) {
    fmt.Fprint(wr, goeliza.ElizaHi())
}

/*
 * JSON webhook structs
 */
type Webhook struct {
    entry []Entry
}

type Entry struct {
    id string
    time int
    messaging []Messaging
}

type Messaging struct {
	sender Sender
	recipient Recipient
	timestamp int
	message Message
}

type Sender struct {
	id string
}

type Recipient struct {
	id string
}

type Message struct {
	mid string
	seq int
	text string
}

/*
 * webhookHandler(wr http.ResponseWriter, req *http.Request) 
 *
 * Handler that lets the FB Messenger API interface
 * with the elizabot at "/webhook" using POST messages
 */
func webhookHandler(wr http.ResponseWriter, req *http.Request) {
	// Verify Facebook validation token
	token := req.URL.Query().Get("hub.verify_token")
	if (token == "quanfucius") {
		fmt.Fprint(wr, req.URL.Query().Get("hub.challenge"))
	}

	// Parse the request in JSON format
	var data Webhook
	dec := json.NewDecoder(req.Body)
	err := dec.Decode(&data);
	
	if err != nil {
		log.Println(err)
		return
	} 
	
	// Loop through messages
	messagingEvents := data.entry[0].messaging;
    for _, event := range messagingEvents {
        if event.message != (Message{}) && event.message.text != "" {
            // Get reply to input message from goeliza
            input := event.message.text
            output := goeliza.ReplyTo(input)

            // Construct Recipient and Message structs
            recipient := Recipient{event.sender.id}
            message := Message{"", 0, output}

            // Reply here
        }
    }
}
