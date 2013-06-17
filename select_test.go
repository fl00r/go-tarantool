package tarantool

import (
	"testing"
	// "fmt"
)

type Employee struct {
	name      String  // index 1
	id        Int32   // index 0
	job       String
	age       Int8
	bestYears []Int32
}

func (employee *Employee) Unpack(cortege [][]byte) (err error) {
	err = employee.name.Unpack(cortege[0])
	err = employee.id.Unpack(cortege[1])
	err = employee.job.Unpack(cortege[2])
	err = employee.age.Unpack(cortege[3])
	length := len(cortege) - 4
	employee.bestYears = make([]Int32, length)
	for i := 0; i < length; i++ {
		employee.bestYears[i] = Int32(0)
		err = employee.bestYears[i].Unpack(cortege[i+4])
	}
	return
}

func TestInsert(t *testing.T) {
	conn, _ := Connect("localhost:33013")
	space := conn.Space(0)

	tuple := []TupleField{ String("Linda"), Int32(3), String("rider"), Int32(21) }
	res, err := space.Insert(tuple, true)

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	if res.Count != 1 {
		t.Errorf("1 tuple should be added not %d", res.Count)
	}
}

func TestSelect(t *testing.T) {
	conn, _ := Connect("localhost:33013")
	space := conn.Space(0)

	res, err := space.Select(0, 0, 10, []TupleField{ String("Linda") }, []TupleField{ String("Mary") })

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	if res.Count != 1 {
		t.Errorf("1 tuple should be selected not %d", res.Count)
	}

	tuple := []TupleField{ String("Mary"), Int32(3), String("singer"), Int32(25) }
	space.Insert(tuple, true)

	res, err = space.Select(0, 0, 10, []TupleField{ String("Linda") }, []TupleField{ String("Mary") })

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	if res.Count != 2 {
		t.Errorf("2 tuples should be selected not %d", res.Count)
	}
}

func TestDelete(t *testing.T) {
	conn, _ := Connect("localhost:33013")
	space := conn.Space(0)

	key := []TupleField{ String("Linda") }
	res, err := space.Delete(key, true)

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	if res.Count != 1 {
		t.Errorf("1 tuple should be deleted not %d", res.Count)
	}

	key = []TupleField{ String("Mary") }
	space.Delete(key, true)
}