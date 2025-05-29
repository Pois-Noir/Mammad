package buffered_reader

import (
	payload "botzila/parser/payload"
	"bufio"
	"encoding/binary"
	"io"
	"net"
)

// StartTCPStreamReader is the entry point to begin reading a TCP stream
// from a given net.Conn. It wraps the connection in a buffered reader
// and pushes framed payloads into the provided channel for decoding.
func StartTCPStreamReader(ioConn net.Conn, payloadChannel chan<- *payload.Payload) {
	// Launch the internal reader with a 4KB buffer and a 4-byte header length.
	// NOTE: 4-byte header implies the first 4 bytes of each message encode the payload length.
	startTCPStreamReader(ioConn, payloadChannel, 4*1024, 4)
}

// startTCPStreamReader is the internal implementation of the TCP reader.
// It reads framed messages from the TCP stream, each consisting of:
//   - A fixed-size header (msgHeaderLen bytes) indicating payload length
//   - A variable-length payload based on that length
//
// It wraps the TCP connection in a bufio.Reader for efficient buffered reads
// and forwards each complete Payload struct to the decoder goroutine via a channel.
func startTCPStreamReader(ioConn net.Conn, payloadChannel chan<- *payload.Payload, bufferSize int, msgHeaderLen int) {
	bufferedReader := bufio.NewReader(ioConn) // Buffered reader to minimize syscalls
	buffHeader := make([]byte, msgHeaderLen)  // Buffer to hold the header (e.g., 4 bytes for uint32 length)
	payloadID := -1                           // Counter to assign unique or sequential IDs to payloads
	var payloadLength uint32 = 0              // Length of the upcoming payload

	for {
		// Read exactly msgHeaderLen bytes to determine the payload length
		_, err := io.ReadFull(bufferedReader, buffHeader)
		if err != nil {
			// TODO: Discuss with Amir and GPT — what’s the appropriate failure strategy?
			// Options:
			//   - Log and continue (e.g., for malformed message)
			//   - Retry with backoff
			//   - Close the connection and signal shutdown
			// Current behavior: break the loop and exit the goroutine
			break
		}

		// Parse the payload length from the header (assuming BigEndian encoding)
		payloadLength = binary.BigEndian.Uint32(buffHeader)
		payloadID++ // Increment ID (you could also use a UUID or hash if preferred)

		// Allocate a new Payload with the given ID and length
		payloadPtr, err := payload.NewPayload(payloadID, int(payloadLength))
		if err != nil {
			// TODO: Decide on recovery strategy with Amir
			// For now: assume unrecoverable and break
			break
		}

		// Send the payload to the decoder goroutine via the channel
		payloadChannel <- payloadPtr
	}
}
