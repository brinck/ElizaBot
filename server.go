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
    Entry []Entry
}

type Entry struct {
    Id string
    Time int
    Messaging []Messaging
}

type Messaging struct {
	Sender Sender
	Recipient Recipient
	Timestamp int
	Message Message
}

type Sender struct {
	Id string
}

type Recipient struct {
	Id string
}

type Message struct {
	Mid string
	Seq int
	Text string
}

type Reply struct {
	Recipient Recipient
	Message Message
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
	var webhookData Webhook
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&webhookData); err != nil {
        log.Println("JSON decoding error:", err, req.Body)
    }

	// Loop through messages
	messagingEvents := webhookData.Entry[0].Messaging;
    for _, event := range messagingEvents {
        if event.Message != (Message{}) && event.Message.Text != "" {
            // Get reply to input message from goeliza
            input := event.Message.Text
            output := goeliza.ReplyTo(input)

            log.Println("Input:", input, "Output:", output)

            // Construct Recipient and Message structs
            recipient := Recipient{event.Sender.Id}
            message := Message{"", 0, output}

            // Reply to user
            webhookReply(recipient, message)
        }
    }
}


/*
 * webhookReply(recipient Recipient, message Message)
 *
 * Function for replying to facebook.
 */
func webhookReply(recipient Recipient, message Message) {
	// Define client and url
	client := http.Client{}
	url := "https://graph.facebook.com/v2.6/me/messages?access_token=" + SecretToken

	// Prepare payload, and encode
	// the payload correctly
	reply := Reply{recipient, message}
	payload, errMarshal := json.Marshal(reply)
	if errMarshal != nil {
		log.Println("Serializing JSON error:", errMarshal)
		return
	}

	// Create stream, set header and
	// create request object
	req, errPost := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	if errPost != nil {
         log.Println("Unable to create post request:", errPost)
         return
    }

	// Execute request
	_, errSend := client.Do(req)
    if errSend != nil {
         log.Println("Unable to reach the server:", errSend)
         return
    }
}
