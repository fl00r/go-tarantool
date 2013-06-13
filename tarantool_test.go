package tarantool

import (
	"testing"
	"encoding/binary"
	"bytes"
	"fmt"
)

func TestConnect(t *testing.T) {
	conn, err := Connect("localhost:33013")
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	space := conn.Space(4)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	field := new(bytes.Buffer)
	err = binary.Write(field, binary.LittleEndian, int32(1))
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	f := SelectField{ 1, field }
	response, err := space.Select(0, 0, 1, &f)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	fmt.Println(response)
}