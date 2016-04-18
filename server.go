package elizabot

import (
    "fmt"
    "log"
    "bytes"
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
 * JSON webhook structs that handle 
 * the parsing of the JSON data
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

type Reply struct {
	recipient Recipient
	message Message
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

            // Reply to user
            webhookReply(recipient, message)
        }
    }
}


/*
 * reply(wr http.ResponseWriter)
 *
 * Function for replying to facebook.
 */
func webhookReply(recipient Recipient, message Message) {

	// Define client and url
	client := http.Client{}
	url := "https://graph.facebook.com/v2.6/me/messages?access_token=" + SecretToken

	// Prepare payload, and encode
	// the payload correctly
	reply := Reply{recipient: recipient, message: message}
	payload, errMarshal := json.Marshal(reply)
	if errMarshal != nil {
		log.Println("Unable to serialise JSON message ", errMarshal)
		return
	}

	// Create stream, set header and
	// create request object
	req, errPost := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/javascript")
	if errPost != nil {
         log.Println("Unable to create post request, ", errPost)
         return
    }

	// Execute request
	_, errSend := client.Do(req)
    if errSend != nil {
         log.Println("Unable to reach the server, ", errSend)
         return
    }
}
