package elizabot

import (
    "bytes"
    "encoding/json"
    "fmt"
    "html/template"
    "github.com/kennysong/goeliza"
    "github.com/urakozz/go-emoji"
    "google.golang.org/appengine"
    "google.golang.org/appengine/log"
    "google.golang.org/appengine/urlfetch"
    "net/http"
    "strings"
)


// TODO
// set global variable that holds information
// on the senders that are communicating
// with Eliza, and the last time they communicated
// TODO
// set a timeout constant

/* 
 * init()
 *
 * Supposedly this function is run by the google containers rather than the 
 * main() function to setup the server bindings
 */
func init() {
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/webhook/", webhookHandler)
    http.HandleFunc("/privacy/", privacyHandler)
}

/*
 * homeHandler(wr http.ResponseWriter, req *http.Request)
 *
 * Handler that renders the page at "/"
 * TODO: Message Us button doesn't seem to be working
 */
func homeHandler(wr http.ResponseWriter, req *http.Request) {
    ctx := appengine.NewContext(req)
    tpl := template.Must(template.ParseGlob("templates/*.html"))
    wr.Header().Set("Content-Type", "text/html; charset=utf-8")
    if err := tpl.ExecuteTemplate(wr, "index.html", nil); err != nil {
        log.Errorf(ctx, "%v", err)
    }
}

/*
 * privacyHandler(wr http.ResponseWriter, req *http.Request)
 *
 * Handler that renders the Privacy Policy at "/privacy/"
 */
func privacyHandler(wr http.ResponseWriter, req *http.Request) {
    ctx := appengine.NewContext(req)
    tpl := template.Must(template.ParseGlob("templates/*.html"))
    wr.Header().Set("Content-Type", "text/html; charset=utf-8")
    if err := tpl.ExecuteTemplate(wr, "privacy.html", nil); err != nil {
        log.Errorf(ctx, "%v", err)
    }
}

/*
 * webhookHandler(wr http.ResponseWriter, req *http.Request) 
 *
 * Handler that lets the FB Messenger API interface
 * with the elizabot at "/webhook/" using POST messages
 */
func webhookHandler(wr http.ResponseWriter, req *http.Request) {
    // Create a GAE Context for this request
    ctx := appengine.NewContext(req)

    // Verify Facebook validation token
    token := req.URL.Query().Get("hub.verify_token")
    if (token == "quanfucius") {
        fmt.Fprint(wr, req.URL.Query().Get("hub.challenge"))
    }

    // Parse the request in JSON format
    var webhookData Webhook
    decoder := json.NewDecoder(req.Body)
    if err := decoder.Decode(&webhookData); err != nil {
        log.Errorf(ctx, "JSON decoding error:\nerr: %v\nreq.Body: %v", err, req.Body)
    }

    // Loop through messages
    messagingEvents := webhookData.Entry[0].Messaging;
    for _, event := range messagingEvents {
        if event.Message != (Message{}) && event.Message.Text != "" {
            // Preprocess and validate input message from goeliza
            input := preprocess(event.Message.Text)
            if notValidInput(input) {
                log.Errorf(ctx, "Not valid input: %v", event.Message.Text)
                return
            }

            // Get reply from ELIZA
            output := goeliza.ReplyTo(input)

            // TODO 
            // get sender id and store it in 
            // the global variable along with a time stamp 
            // TODO
            // if this is the first time the sender is 
            // communicating with Eliza, send a hello message
            // TODO 
            // loop over sender objects and check for any
            // senders that have been idle for more than
            // the time out variable. If so, send them a
            // "goodbye" message and delete them from
            // the object

            log.Debugf(ctx, "Input: \"%s\"\nOutput: \"%s\"", input, output)

            // Construct Recipient and Message structs
            recipient := Recipient{event.Sender.Id}
            message := Message{"", 0, output}

            // Reply to user
            webhookReply(recipient, message, req)
        }
    }
}

/*
 * preprocess(statement string) string
 *
 * Preprocesses a string for ELIZA
 */
func preprocess(statement string) string {
    // Remove all emoji
    parser := emoji.NewEmojiParser()
    statement = parser.ReplaceAllStringFunc(statement, func(s string) string {
        return ""
    })

    // Trim whitespace
    statement = strings.TrimSpace(statement)
    
    return statement
}

/*
 * notValidInput(statement string) bool
 *
 * Checks if the statement is a valid input for Eliza (not empty)
 */
func notValidInput(statement string) bool {
    return strings.TrimSpace(statement) == ""
}

/*
 * webhookReply(recipient Recipient, message Message)
 *
 * Function for replying to facebook.
 */
func webhookReply(recipient Recipient, message Message, req *http.Request) {
    // Create a GAE Context and urlfetch client for this request
    ctx := appengine.NewContext(req)
    client := urlfetch.Client(ctx)
    url := "https://graph.facebook.com/v2.6/me/messages?access_token=" + SecretToken

    // Prepare payload, and encode
    // the payload correctly
    reply := Reply{recipient, message}
    payload, errMarshal := json.Marshal(reply); if errMarshal != nil {
        log.Errorf(ctx, "Serializing JSON error: %s", errMarshal)
        return
    }
    
    log.Debugf(ctx, "payload = %v", string(payload))

    // Create stream, set header and
    // create request object
    req, errPost := http.NewRequest("POST", url, bytes.NewBuffer(payload)); if errPost != nil {
         log.Errorf(ctx, "Unable to create post request: %s", errPost)
         return
    }
    req.Header.Set("Content-Type", "application/json")
    
    // Execute request
    _, errSend := client.Do(req); if errSend != nil {
         log.Errorf(ctx, "Unable to reach the server: %s", errSend)
         return
    }
}
