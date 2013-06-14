package tarantool

import (
	"testing"
)

func TestConnect(t *testing.T) {
	conn, err := Connect("localhost:33013")
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	conn.Space(4)
}