package dns

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
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
	Name  string // (Variable) The domain name, encoded as a sequence of labels.
	Type  uint16 // (16 bits Big-Endian) The record type. (1 for an A record, 5 for a CNAME record etc..., see RFC 1035 for a full list of types.)
	Class uint16 // (16 bits Big-Endian) The class, in practice always set to 1.
}

// TODO: Refactoriser et rendre la méthode privé.
func (q *Question) Read() {
	v := reflect.ValueOf(*q)
	t := v.Type()

	fmt.Println("+----------------+----------------+")
	fmt.Println("             Question              ")
	fmt.Println("+----------------+----------------+")
	fmt.Println("|     Field      |     Value     |")
	fmt.Println("+----------------+----------------+")

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if value.Kind() == reflect.Slice {
			fmt.Printf("| %15s | %v \n", field.Name, string(value.Bytes()))
		} else {
			fmt.Printf("| %15s | %v \n", field.Name, value.Interface())
		}
		fmt.Println("+----------------+----------------+")
	}
}

func ReadLabel(data []byte, offset int) (string, int, *ParseError) {
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

func NewQuestion(buffer []byte) (*Question, *ParseError) {
	var domain string

	domain, offset := DecodeDomainToString(buffer)

	return &Question{
		Name:  domain,
		Type:  binary.BigEndian.Uint16(buffer[offset+1 : offset+3]),
		Class: binary.BigEndian.Uint16(buffer[offset+3 : offset+5]),
	}, nil
}

func DecodeDomainToString(buffer []byte) (string, int) {
	var labels []string
	var domain string

	offset := 0

	for {
		label, nextOffset, _ := ReadLabel(buffer, offset)
		if label == "" {
			break
		}
		labels = append(labels, label)
		offset = nextOffset
	}

	domain = strings.Join(labels, ".")

	return domain, offset
}

// TODO: Refactoriser et rendre la méthode privé.
func EncodeDomainToBytes(domain string) []byte {
	var buffer bytes.Buffer

	labels := strings.Split(domain, ".")
	for _, label := range labels {
		buffer.WriteByte(byte(len(label)))
		buffer.WriteString(label)
	}
	buffer.WriteByte(0x00)

	return buffer.Bytes()
}

func (q *Question) Write() []byte {
	var buffer bytes.Buffer

	buffer.Write(EncodeDomainToBytes(q.Name))
	binary.Write(&buffer, binary.BigEndian, q.Type)
	binary.Write(&buffer, binary.BigEndian, q.Class)

	fmt.Println("Question buffer: ", len(buffer.Bytes()))

	return buffer.Bytes()
}
