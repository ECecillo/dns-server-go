package dns

import (
	"encoding/binary"
	"fmt"
	"reflect"
)

// https://github.com/EmilHernvall/dnsguide/blob/b52da3b32b27c81e5c6729ac14fe01fef8b1b593/chapter1.md

// 12 bytes long
// Intger in Big Endian
type Header struct {
	ID      uint16 // (16 bits Big-Endian) Package identifier : A random ID to query packets. Response packet must reply with the smae ID
	QR      byte   // (1 bit)   Query/Response Indicator : 1 for reply, 0 for a question packet.
	OPCODE  byte   // (4 bits)  Operation Code : Specifies the kind of query in a message.
	AA      byte   // (1 bit)   Authoritative Answer: 1 if the responding server "owns" the domain queried, i.e., it's authoritative.
	TC      byte   // (1 bit)   Truncation : 1 if the message is larger than 512 bytes. Always 0 in UDP responses.
	RD      byte   // (1 bit)   Recursion Desired : Sender sets this to 1 if the server should recursively resolve this query, 0 otherwise.
	RA      byte   // (1 bit)   Recursion Available : Server sets this to 1 to indicate that recursion is available.
	Z       byte   // (3 bits)  Reserved : Used by DNSSEC queries. At inception, it was reserved for future use.
	RCODE   byte   // (4 bits)  Response code indicating the status of the response.
	QDCOUNT uint16 // (16 bits Big-Endian) Question Count : Number of questions in the Question section.
	ANCOUNT uint16 // (16 bits Big-Endian) Answer Record Count : Number of records in the Answer section.
	NSCOUNT uint16 // (16 bits Big-Endian) Authority Record Count : Number of records in the Authority section.
	ARCOUNT uint16 // (16 bits Big-Endian) Additional Record Count : Number of records in the Additional section.
}

// TODO: Refactoriser et rendre la méthode privé.
func (h *Header) Read() {
	v := reflect.ValueOf(*h)
	t := v.Type()

	fmt.Println("+----------------+----------------+")
	fmt.Println("               Header              ")
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

// +-------+-------+----------------------------------+---------+---------+---------+---------+
// |  _ _  |  _ _  | 0 0 0 0 0 0 0 1  0 0 1 0 0 0 0 0 |  _   _  |   _ _   |  _   _  |  _	 _  |
// | ID[0] + ID[1] | - -+-+-+- - - -  - -+-+- -+-+-+- | QDCOUNT | ANCOUNT | NSCOUNT | ARCOUNT |
//                   Q    O    A T R  R   Z      R
//                   R    P    A C D  A          C
//                        C                      O
//                        O                      D
//                        D                      E
//                        E
//+-------+-------+----------------------------------+---------+---------+---------+---------+

// Write Header segment of a DNS request into an array of 12 bytes.
func (h *Header) Write() [12]byte {
	var buffer [12]byte

	binary.BigEndian.PutUint16(buffer[:2], h.ID)
	buffer[2] = h.QR<<7 | h.OPCODE<<3 | h.AA<<2 | h.TC<<1 | h.RD
	buffer[3] = h.RA<<7 | h.Z<<4 | h.RCODE
	binary.BigEndian.PutUint16(buffer[4:6], h.QDCOUNT)
	binary.BigEndian.PutUint16(buffer[6:8], h.ANCOUNT)
	binary.BigEndian.PutUint16(buffer[8:10], h.NSCOUNT)
	binary.BigEndian.PutUint16(buffer[10:12], h.ARCOUNT)

	return buffer
}

func NewHeader(buffer [12]byte) *Header {
	return &Header{
		ID:      binary.BigEndian.Uint16(buffer[:2]),
		QR:      buffer[2] >> 7,        // Shift 7 bits to the right to get the first bit.
		OPCODE:  buffer[2] >> 3 & 0x0F, // Shift 3 bits to the right to get the 4 bits of the OPCODE using AND 0000 1111.
		AA:      buffer[2] >> 2 & 0x01, // Shift 2 bits to the right to get the 1 bit of the AA using AND 000 0001.
		TC:      buffer[2] >> 1 & 0x01, // Shift 1 bit to the right to get the 1 bit of the TC using AND 000 0001.
		RD:      buffer[2] & 0x01,      // Apply AND with 0000 0001.
		RA:      buffer[3] >> 7,        // ....
		Z:       buffer[3] >> 4 & 0x07,
		RCODE:   buffer[3] & 0x0F,
		QDCOUNT: binary.BigEndian.Uint16(buffer[4:6]),
		ANCOUNT: binary.BigEndian.Uint16(buffer[6:8]),
		NSCOUNT: binary.BigEndian.Uint16(buffer[8:10]),
		ARCOUNT: binary.BigEndian.Uint16(buffer[10:12]),
	}
}
