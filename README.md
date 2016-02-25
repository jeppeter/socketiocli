# Socket.IO [![Build Status](https://travis-ci.org/jeppeter/socketiocli.png?branch=master)](https://travis-ci.org/jeppeter/socketiocli)

Package socketiocli implements a client for SocketIO protocol in Go language as specified in 
[socket.io-spec](https://github.com/LearnBoost/socket.io-spec)

## Usage

```go
package main

import (
	"fmt"
	"github.com/jeppeter/socketiocli"
)

func main() {
	// Open a new client connection to the given socket.io server
	// Connect to the given channel on the socket.io server
	socket, err := socketiocli.DialAndConnect("socketiocli-server.com:80", "/channel", "key=value")
	if err != nil {
		panic(err)
	}

	for {
		// Receive socketiocli.Message from the server
		msg, err := socket.Receive()
		if err != nil {
			panic(err)
		}

		fmt.Printf("Type: %v, ID: '%s', Endpoint: '%s', Data: '%s' \n", msg.Type, msg.ID, msg.Endpoint, msg.Data)
	}
}
```

## Documentation 

http://godoc.org/github.com/jeppeter/socketiocli

## License

The MIT License (MIT)
