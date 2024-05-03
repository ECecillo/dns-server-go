package dns

import (
	"reflect"
	"testing"
)

func TestReadLabel(t *testing.T) {
	data := []byte{0x03, 'w', 'w', 'w', 0x06, 'g', 'o', 'o', 'g', 'l', 'e', 0x03, 'c', 'o', 'm', 0x00}
	expectedLabel := "www"
	expectedOffset := 4

	label, offset, _ := ReadLabel(data, 0)
	if label != expectedLabel || offset != expectedOffset {
		t.Errorf("readLabel failed, expected '%s' and %d, got '%s' and %d", expectedLabel, expectedOffset, label, offset)
	}
}

func TestParse(t *testing.T) {
	data := []byte{0x03, 'w', 'w', 'w', 0x06, 'g', 'o', 'o', 'g', 'l', 'e', 0x03, 'c', 'o', 'm', 0x00, 0x00, 0x01, 0x00, 0x01}
	expectedQuestion := &Question{
		Name:  "www.google.com",
		Type:  1,
		Class: 1,
	}

	question, _ := NewQuestion(data)

	if !reflect.DeepEqual(question, expectedQuestion) {
		t.Errorf("Parse failed, expected %+v, got %+v", expectedQuestion, question)
	}
}

func TestReadLabelWithIncompleteData(t *testing.T) {
	data := []byte{0x03, 'w', 'w', 'w', 0x06}
	expectedLabel := ""
	expectedOffset := -1
	label, offset, err := ReadLabel(data, 0)

	if err != nil && label != expectedLabel && offset != expectedOffset {
		t.Errorf("Reading incomplete data should throw an error and an empty label, got %+v", label)
	}
}

func TestEncodeDomainToBytes(t *testing.T) {
	domain := "www.google.com"
	expected := []byte{0x03, 'w', 'w', 'w', 0x06, 'g', 'o', 'o', 'g', 'l', 'e', 0x03, 'c', 'o', 'm', 0x00}

	result := EncodeDomainToBytes(domain)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("EncodeDomainToBytes failed, expected %v, got %v", expected, result)
	}
}
