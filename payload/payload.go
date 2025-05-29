package payload

// Payload represents a fully framed and raw message that has been
// read from the TCP stream by the buffered reader goroutine.
// This struct is passed through a channel to the decoder goroutine,
// which is responsible for interpreting the byte slice according to
// your application-level protocol.
type Payload struct {
	PayloadID int    // A unique or protocol-defined identifier for the message type or source
	Payload   []byte // Raw byte content of the message; not yet decoded
}

// NewPayload constructs a new Payload instance with a given payload ID
// and raw byte slice. This is typically called by the buffered reader
// once a full message has been read and framed.
//
// TODO: Discuss with Amir what additional validation might be useful.
// For example:
//   - Should PayloadID be within a known range?
//   - Should Payload have a minimum or maximum length?
//   - Should certain reserved IDs trigger errors?
func NewPayload(payloadID int, payloadLength int) (*Payload, error) {
	// Basic construction — you can add checks here as needed
	return &Payload{
		PayloadID: payloadID,
		Payload:   make([]byte, payloadLength),
	}, nil
}
