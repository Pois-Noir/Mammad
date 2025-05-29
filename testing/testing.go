package testing

import (
	"net"
)

type Chunk struct {
	chunkID int
	Payload []byte
}

func ListenConnections(listener net.Listener) {
	// do some error checking

	// cycle of life
	for {
		conn, err := listener.Accept()
		if err != nil {
			// log the error
			continue
		}
		// handle the connection
		go handleConnection(conn)

	}

}

func handleConnection(conn net.Conn) {

	// create a buffer for the buffered reader

	go func() {
		// create a buffered reader

		// create a for loop

\			//io.ReadFull(//mssgHeader)
			// find out message length
			//io.ReadFull(//Payload)
			// make a message struct
			// pass it to a channel

	}()

	go func() {
		// read from the same channel 
		// start decoding baby

	}()

}

// I need two different  states read header and then read payload
