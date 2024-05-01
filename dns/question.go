package dns

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

// ParseError records an error
type ParseError struct {
	Position int   // Position where the error occured
	Length   int   // Length of the question label where the error occured.
	Err      error // The actual error.
}

// These are the errors that can be returned in [ParseError.Err].
var (
	ErrLabelLength = errors.New("length exceed buffer size")
)

func (e *ParseError) Error() string {
	if e.Err == ErrLabelLength {
		return fmt.Sprintf("label at position %d, length %d : %v", e.Position, e.Length, e.Err)
	}
	return fmt.Sprintf("unknown parse error at position %d, length %d: %v", e.Position, e.Length, e.Err)
}

func (e *ParseError) Unwrap() error { return e.Err }

type Question struct {
	Name  string // The domain name, encoded as a sequence of labels.
	Type  uint16 // The record type. (1 for an A record, 5 for a CNAME record etc..., see RFC 1035 for a full list of types.)
	Class uint16 // The class, in practice always set to 1.
}

func (q *Question) Read() {
	fmt.Println("+-------+-------+")
	fmt.Println("     Question    ")
	fmt.Println("+-------+-------+")
	fmt.Println("| Field | Value |")
	fmt.Println("+-------+-------+")
	fmt.Println("| Name  |", q.Name)
	fmt.Println("+-------+-------+")
	fmt.Println("| Type  |", q.Type)
	fmt.Println("+-------+-------+")
	fmt.Println("| Class |", q.Class)
	fmt.Println("+-------+-------+")
}

func readLabel(data []byte, offset int) (string, int, *ParseError) {
	position := offset
	labelLength := int(data[position])
	if labelLength == 0 {
		return "", offset, nil
	}
	if labelLength > len(data) {
		return "", -1, &ParseError{
			Position: position,
			Err:      ErrLabelLength,
		}
	}

	nextOffset := position + labelLength + 1
	label := string(data[position+1 : nextOffset])
	return label, nextOffset, nil
}

func (q *Question) Create(buffer []byte) (*Question, *ParseError) {
	var labels []string
	var domain string

	offset := 0

	for {
		label, nextOffset, err := readLabel(buffer, offset)
		if err != nil {
			return nil, err
		}
		if label == "" {
			break
		}
		labels = append(labels, label)
		offset = nextOffset
	}

	domain = strings.Join(labels, ".")

	return &Question{
		Name:  domain,
		Type:  binary.BigEndian.Uint16(buffer[offset+1 : offset+3]),
		Class: binary.BigEndian.Uint16(buffer[offset+3 : offset+5]),
	}, nil
}
