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
type message struct {
    entries[] Entry
}

type Entry struct {
    id string,
    time int,
    messages[]
}

/*
 * eliza(w http.ResponseWriter, r *http.Request) 
 *
 * Handler that lets the messenger API interface
 * with the elizabot.
 */
func eliza(wr http.ResponseWriter, req *http.Request) 
{
	// parse the request json format, see stack overflow for reference
	// http://stackoverflow.com/questions/15672556/handling-json-post-request-in-go
	decoder := json.NewDecoder(req.Body);
	
	messagingEvents := req.Body.entry[0].messaging;
	for i := 0; i < messagingEvents
}
