package encoder

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"reflect"

	g_errors "github.com/Pois-Noir/Mammad/errors"
)

const (
	TypeString  = 0x01
	TypeInt64   = 0x02
	TypeFloat64 = 0x03
	TypeBool    = 0x04
	TypeMap     = 0x05
	TypeSlice   = 0x06
)

// struct responsible for encoding a map[string]interface{}
type Encoder struct {
	// pointer to a buffer
	// will store the byte stream
	buf *bytes.Buffer
}

// Constructor to create a new Encoder
func NewEncoder() *Encoder {
	return &Encoder{buf: new(bytes.Buffer)}
}

// // public function exposed to the user
// // we need to pass the target map
// // it will encode it in the form | value type (1 byte) | value len (2 bytes) | payload |
// func (e *Encoder) EncodeMap(m map[string]interface{}) ([]byte, error) {
// 	// loop through the map
// 	for key, val := range m {
// 		// encode the key, val entry
// 		if err := e.writeEntry(key, val); err != nil {
// 			return nil, err
// 		}
// 	}
// 	// return the byte slice
// 	return e.buf.Bytes(), nil
// }

// temp for testing
// got this from gpt
func (e *Encoder) EncodeMap(m map[string]interface{}) ([]byte, error) {
	for key, val := range m {
		if err := e.writeEntry(key, val); err != nil {
			return nil, err
		}
	}

	// calculate message size
	payload := e.buf.Bytes()
	return payload, nil
}

// writeEntry = [ keyEntry | valueEntry ]
func (e *Encoder) writeEntry(key string, val interface{}) error {
	// 1) write key as string
	if err := e.writePrimitive(TypeString, []byte(key)); err != nil {
		return err
	}
	// 2) write the value
	return e.writeValue(val)
}

// func for encoding the values
func (e *Encoder) writePrimitive(typeByte byte, payload []byte) error {
	// type
	if err := e.buf.WriteByte(typeByte); err != nil {
		return err
	}
	// we need to make sure the length is not above 2 ^ 16 (highly unlikely but better to be save)
	if len(payload) >= 1<<16 {
		return g_errors.ErrByteOverFlow
	}
	// we write to the buffer the length of the data we use 2 bytes
	length := uint16(len(payload))
	if err := binary.Write(e.buf, binary.BigEndian, length); err != nil {
		return err
	}
	// we are writing the actual data
	_, err := e.buf.Write(payload)
	return err
}

func (e *Encoder) writeValue(v interface{}) error {
	// Use a type‐switch on the interface value to pick the correct encoding path.
	switch val := v.(type) {

	// STRING
	case string:
		// Convert the Go string to its raw UTF-8 bytes, then write with the string type marker.
		return e.writePrimitive(TypeString, []byte(val))

	// INTEGER VARIANTS
	case int, int8, int16, int32, int64:
		// We want a uniform 8-byte representation, so:
		//  1) Make an 8-byte buffer.
		//  2) Use reflection to extract the signed integer as int64.
		//  3) Cast that to uint64 (binary.BigEndian.PutUint64 expects unsigned).
		buf := make([]byte, 8)
		// important to use reflect.ValueOf(val)
		// or else i would have to seperately check int, int8, int16 ....
		// it returns int64 no matter the type
		binary.BigEndian.PutUint64(buf, uint64(reflect.ValueOf(val).Int()))
		// Write with the int64 type marker.
		return e.writePrimitive(TypeInt64, buf)

	// FLOAT VARIANTS
	case float32, float64:
		// reflection.Float always returns float64, even if the original was float32.
		f := reflect.ValueOf(val).Float()
		// Prepare an 8-byte buffer for the IEEE-754 bits.
		buf := make([]byte, 8)
		// math.Float64bits converts a float64 into its IEEE-754 binary layout.
		binary.BigEndian.PutUint64(buf, math.Float64bits(f))
		// Write with the float64 type marker.
		return e.writePrimitive(TypeFloat64, buf)

	// BOOLEAN
	case bool:
		// Encode false as 0x00, true as 0x01.
		b := byte(0)
		if val {
			b = 1
		}
		// Single‐byte payload.
		return e.writePrimitive(TypeBool, []byte{b})

	// NESTED MAP
	case map[string]interface{}:
		// if the value is a map
		// we can recursively encode, using the custom encoder
		tmp := NewEncoder()
		// EncodeMap writes each key/value pair in the map into tmp.buf.
		if _, err := tmp.EncodeMap(val); err != nil {
			return err
		}
		// Then write the entire sub-buffer as a single "map" primitive.
		return e.writePrimitive(TypeMap, tmp.buf.Bytes())

	// SLICES OF INTERFACE{}
	case []interface{}:
		// for lists it is a bit different compared to maps
		// we allow the list to be of holding value of type interface{}
		// to encode all of them we need to use the Encoder.writeValue function
		// we have to manually call the .writeValue function on each of the values of the list
		tmp := NewEncoder()
		for _, elem := range val {
			if err := tmp.writeValue(elem); err != nil {
				return err
			}
		}
		// Emit the concatenated slice bytes as one "slice" primitive.
		return e.writePrimitive(TypeSlice, tmp.buf.Bytes())

	// UNSUPPORTED TYPES
	default:
		// If we hit anything else (structs, pointers, custom types…), error out.
		return fmt.Errorf("unsupported type: %T", v)
	}
}
