package testing

import (
	"bufio"
	"io"
	"net"
	"encoding/binary"
)

type Payload struct{
	PayloadID int
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
	headerBuffer := make([]byte , 4)
	payloadID := -1
	var payloadLen uint32 
	payloadChannel := make(chan *Payload , 100)

	go func() {
		// create a buffered reader
		buffReader := bufio.NewReader(conn)

		for{

			_, err := io.ReadFull(buffReader , headerBuffer)
			if err != nil {
				// do smth sensible
			}
			payloadLen = binary.BigEndian.Uint32(headerBuffer)
			payloadID++
			payload := &Payload{
				PayloadID: (payloadID),
				Payload: make([]byte , payloadLen),
			}
			_, err = io.ReadFull(buffReader , payload.Payload)
			if err != nil{
				// do smth sensible
			}
			payloadChannel <- payload

		}
		// create a for loop

			//io.ReadFull(//mssgHeader)
			// find out message length
			//io.ReadFull(//Payload)
			// make a message struct
			// pass it to a cha
	}()

	go func() {
		

	}()

}

// I need two different  states read header and then read payload
// have a utility package 

