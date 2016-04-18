# ElizaBot

ElizaBot is a Facebook Messenger bot that simulates a Rogerian psychotherapist. The chatbot server is written in Go and runs on Google App Engine.

Based on Weizenbaum's 1966 [ELIZA chatbot](https://en.wikipedia.org/wiki/ELIZA), and uses the [goeliza package](https://github.com/kennysong/goeliza).

### Usage

To install dependencies:
```
go get github.com/kennysong/goeliza
go get google.golang.org/appengine
```

To run the webapp on the local GAE development server:
```
goapp serve
```

To deploy the code to GAE:
```
appcfg.py -A elizabot-fb update .
```
