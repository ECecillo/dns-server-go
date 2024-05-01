package dns

import (
	"reflect"
	"testing"
)

func TestReadLabel(t *testing.T) {
	data := []byte{0x03, 'w', 'w', 'w', 0x06, 'g', 'o', 'o', 'g', 'l', 'e', 0x03, 'c', 'o', 'm', 0x00}
	expectedLabel := "www"
	expectedOffset := 4

	label, offset, _ := readLabel(data, 0)
	if label != expectedLabel || offset != expectedOffset {
		t.Errorf("readLabel failed, expected '%s' and %d, got '%s' and %d", expectedLabel, expectedOffset, label, offset)
	}
}

func TestCreate(t *testing.T) {
	data := []byte{0x03, 'w', 'w', 'w', 0x06, 'g', 'o', 'o', 'g', 'l', 'e', 0x03, 'c', 'o', 'm', 0x00, 0x00, 0x01, 0x00, 0x01}
	expectedQuestion := &Question{
		Name:  "www.google.com",
		Type:  1,
		Class: 1,
	}

	q := new(Question)
	question, _ := q.Create(data)

	if !reflect.DeepEqual(question, expectedQuestion) {
		t.Errorf("Create failed, expected %+v, got %+v", expectedQuestion, question)
	}
}

func TestCreateWithIncompleteData(t *testing.T) {
	data := []byte{0x03, 'w', 'w', 'w', 0x06}
	q := new(Question)
	question, err := q.Create(data)

	if err != nil && question != nil {
		t.Errorf("Create with incomplete data should result in default Question, got %+v", question)
	}
}
