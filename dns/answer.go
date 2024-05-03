package dns

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

// Reference : https://www.rfc-editor.org/rfc/rfc1035#section-3.2.1
type Answer struct {
	Name     string // (Variable) The domain name encoded as a sequence of labels.
	Type     uint16 // (16 bits Big-Endian) 1 for an A record, 5 for a CNAME record etc., full list https://www.rfc-editor.org/rfc/rfc1035#section-3.2.2
	Class    uint16 // (16 bits Big-Endian) Usually set to 1 (full list https://www.rfc-editor.org/rfc/rfc1035#section-3.2.4)
	TTL      uint32 // (32 bits Big-Endian) The duration in seconds a record can be cached before requerying.
	RDLENGTH uint16 // (16 bits Big-Endian) Length of the RDATA field in bytes.
	RDATA    []byte // (Variable) Data specific to the record type.
}

// TODO: Refactoriser et rendre la méthode privé.
func (a *Answer) Read() {
	v := reflect.ValueOf(*a)
	t := v.Type()

	fmt.Println("+----------------+----------------+")
	fmt.Println("               Answer              ")
	fmt.Println("+----------------+----------------+")
	fmt.Println("|     Field      |     Value     |")
	fmt.Println("+----------------+----------------+")

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Pour les champs de type '[]byte' (comme RDATA), convertissez-les en string pour une meilleure lisibilité
		if value.Kind() == reflect.Slice {
			fmt.Printf("| %15s | %v \n", field.Name, string(value.Bytes()))
		} else {
			fmt.Printf("| %15s | %v \n", field.Name, value.Interface())
		}
		fmt.Println("+----------------+----------------+")
	}
}

func (a *Answer) Write() []byte {
	var buffer bytes.Buffer

	// TODO: Ne pas appeler la méthode publique.
	buffer.Write(EncodeDomainToBytes(a.Name))
	binary.Write(&buffer, binary.BigEndian, a.Type)
	binary.Write(&buffer, binary.BigEndian, a.Class)
	binary.Write(&buffer, binary.BigEndian, a.TTL)
	binary.Write(&buffer, binary.BigEndian, a.RDLENGTH)

	return buffer.Bytes()
}
