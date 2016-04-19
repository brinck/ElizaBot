package elizabot

/*
 * JSON webhook structs that handle parsing the JSON data
 */
type Webhook struct {
    Entry []Entry `json:"entry"`
}

type Entry struct {
    Id int64 `json:"id"`
    Time int64 `json:"time"`
    Messaging []Messaging `json:"messaging"`
}

type Messaging struct {
    Sender Sender `json:"sender"`
    Recipient Recipient `json:"recipient"`
    Timestamp int64 `json:"timestamp"`
    Message Message `json:"message"`
}

type Sender struct {
    Id int64 `json:"id"`
}

type Reply struct {
    Recipient Recipient `json:"recipient"`
    Message Message `json:"message"`
}

type Recipient struct {
    Id int64 `json:"id"`
}

type Message struct {
    Mid string `json:"mid,omitempty"` 
    Seq int64 `json:"seq,omitempty"`
    Text string `json:"text"`
}
