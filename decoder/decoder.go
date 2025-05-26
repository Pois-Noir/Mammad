package decoder

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

// Type markers must match your encoder.
const (
	TypeString  = 0x01
	TypeInt64   = 0x02
	TypeFloat64 = 0x03
	TypeBool    = 0x04
	TypeMap     = 0x05
	TypeSlice   = 0x06
)

// Decoder wraps any io.Reader and produces a map[string]interface{}.
type Decoder struct {
	r io.Reader
}

// NewDecoder returns a Decoder over the supplied reader.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// DecodeBytes is a convenience for decoding a single []byte payload.
func DecodeBytes(data []byte) (map[string]interface{}, error) {
	return NewDecoder(bytes.NewReader(data)).Decode()
}

// Decode reads key/value pairs until EOF, and returns the assembled map.
func (d *Decoder) Decode() (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for {
		key, err := d.readString()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		val, err := d.readValue()
		if err != nil {
			return nil, err
		}
		result[key] = val
	}
	return result, nil
}

// readType reads a single type-marker byte.
func (d *Decoder) readType() (byte, error) {
	var t [1]byte
	if _, err := io.ReadFull(d.r, t[:]); err != nil {
		return 0, err
	}
	return t[0], nil
}

// readUint16 reads a big-endian 2-byte length.
func (d *Decoder) readUint16() (uint16, error) {
	var buf [2]byte
	if _, err := io.ReadFull(d.r, buf[:]); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(buf[:]), nil
}

// readBytes reads exactly n bytes or returns an error.
func (d *Decoder) readBytes(n uint16) ([]byte, error) {
	b := make([]byte, n)
	if _, err := io.ReadFull(d.r, b); err != nil {
		return nil, err
	}
	return b, nil
}

// readString expects a string marker, then length, then UTF-8 bytes.
func (d *Decoder) readString() (string, error) {
	t, err := d.readType()
	if err != nil {
		return "", err
	}
	if t != TypeString {
		return "", fmt.Errorf("decoder: expected string type, got 0x%02x", t)
	}
	length, err := d.readUint16()
	if err != nil {
		return "", err
	}
	data, err := d.readBytes(length)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// readValue handles all supported value types.
func (d *Decoder) readValue() (interface{}, error) {
	t, err := d.readType()
	if err != nil {
		return nil, err
	}
	length, err := d.readUint16()
	if err != nil {
		return nil, err
	}
	payload, err := d.readBytes(length)
	if err != nil {
		return nil, err
	}

	switch t {
	case TypeString:
		return string(payload), nil

	case TypeInt64:
		if len(payload) != 8 {
			return nil, fmt.Errorf("decoder: invalid int64 length %d", len(payload))
		}
		return int64(binary.BigEndian.Uint64(payload)), nil

	case TypeFloat64:
		if len(payload) != 8 {
			return nil, fmt.Errorf("decoder: invalid float64 length %d", len(payload))
		}
		bits := binary.BigEndian.Uint64(payload)
		return math.Float64frombits(bits), nil

	case TypeBool:
		if len(payload) != 1 {
			return nil, fmt.Errorf("decoder: invalid bool length %d", len(payload))
		}
		return payload[0] == 1, nil

	case TypeMap:
		// nested map: recurse on its payload
		nested, err := DecodeBytes(payload)
		if err != nil {
			return nil, err
		}
		return nested, nil

	case TypeSlice:
		// slices are just concatenated values: keep reading until EOF
		subDec := NewDecoder(bytes.NewReader(payload))
		var list []interface{}
		for {
			v, err := subDec.readValue()
			if err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}
			list = append(list, v)
		}
		return list, nil

	default:
		return nil, fmt.Errorf("decoder: unsupported type marker 0x%02x", t)
	}
}
