var express = require('express');
var eliza = require('elizabot');
var format = require('util').format;
var app = express();

app.get('/', function (req, res) {
	res.send('Hello World!');
});

// create ElizaBot object to interact
// with facebook users
var elizas = {};



// create wrapper function for replying
// to the sender
function elizaReply(sender, text) {
	// define message object
	reply = {
		text:text
	}; 

	// make request
	request({
		url: 'https://graph.facebook.com/v2.6/me/messages',
		qs: { access_token: token },
		method: 'POST',
		json: {
			recipient: { id: sender },
			message: reply,
		}
	// check callback for errors to log
	}, function(error, response, body) {
		if (error) {
			console.log(format('Error sending message to senderid %s: %s', sender, error));
		} else if (response.body.error) {
			console.log('Error: ', response.body.error);
		}
	});
}



// setup webhook to interact with facebook
// API 
app.post('/webhook/', function (req, res) {
	// get message events and loop through them
	messaging_events = req.body.entry[0].messaging;
  	for (var i = 0; i < messaging_events.length; i++) {
  		// get message event info and 
  		// check for message
    	event = req.body.entry[0].messaging[i];
    	sender = event.sender.id;

		//----
		// check if we need to create an individual 
		// eliza bot for each messenger
		//---- 	
    	if (event.message && event.message.text) {
    		// if the sender has not been registered
    		// make sure to create a bot for them
    		if (!(sender in elizas)) {
    			elizas[sender] = new ElizaBot();
    			var initial = eliza.getInitial();

    			//----
    			// insert logic here
    			//----
    			/*var reply = eliza.transform(inputstring);
				if (eliza.quit) {
					// last user input was a quit phrase
				}

				// method `transform()' returns a final phrase in case of a quit phrase
				// but you can also get a final phrase with:
				var final = eliza.getFinal();

				// other methods: reset memory and internal state
				eliza.reset();*/
				//----
    		}
    	}
    }

  	res.sendStatus(200);
});

app.listen(80, function () {
  console.log('Example app listening on port 80!');
});
