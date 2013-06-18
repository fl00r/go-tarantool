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

	res, err := space.Select(0, 0, 10, []TupleField{ String("Linda") }, []TupleField{ String("Mary") })

	if res.Count != 0 {
		t.Errorf("0 tuple should be added not %d", res.Count)
	}

	tuple := []TupleField{ String("Linda"), Int32(1), String("rider"), Int32(21) }
	res, err = space.Insert(tuple, true)

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	if res.Count != 1 {
		t.Errorf("1 tuple should be added not %d", res.Count)
	}
}

func TestSelectAddReplace(t *testing.T) {
	conn, _ := Connect("localhost:33013")
	space := conn.Space(0)

	res, err := space.Select(0, 0, 10, []TupleField{ String("Linda") }, []TupleField{ String("Mary") })

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	if res.Count != 1 {
		t.Errorf("1 tuple should be selected not %d", res.Count)
	}

	tuple := []TupleField{ String("Mary"), Int32(2), String("singer"), Int32(25) }
	space.Add(tuple, true)

	tuple = []TupleField{ String("Linda"), Int32(1), String("guitarrist"), Int32(21) }
	space.Replace(tuple, true)

	res, err = space.Select(0, 0, 10, []TupleField{ String("Linda") }, []TupleField{ String("Mary") })

	var lindaJob String
	lindaJob.Unpack(res.Tuples[0].Fields[2])

	if lindaJob != "guitarrist" {
		t.Errorf("Error: Linda's job should be replaced with guitarrist, not %s", lindaJob)
	}

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	if res.Count != 2 {
		t.Errorf("2 tuples should be selected not %d", res.Count)
	}
}

func TestCall(t *testing.T) {
	conn, _ := Connect("localhost:33013")
	space := conn.Space(0)
	res, err := space.Call("box.select_range", true, []TupleField{ String("0"), String("0"), String("10")})

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	if res.Count != 2 {
		t.Errorf("2 tuples should be returned not %d", res.Count)
	}

	res, err = space.Call("box.select_range", true, []TupleField{ String("0"), String("1"), String("10"), Int32(2)})

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	if res.Count != 1 {
		t.Errorf("1 tuples should be returned not %d", res.Count)
	}

	res, err = space.Call("box.dostring", true, []TupleField{ String("return box.space[0].index[0]:len()") })

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	var count Int32
	count.Unpack(res.Tuples[0].Fields[0])

	if count != 2 {
		t.Errorf("Count should be 2 not %d", count)
	}

	if res.Count != 1 {
		t.Errorf("1 tuples should be returned not %d", res.Count)
	}
}

func TestPing(t *testing.T) {
	conn, _ := Connect("localhost:33013")
	space := conn.Space(0)
	res, err := space.Ping()

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	if res.Count != 0 {
		t.Errorf("Ping should return nothing")
	}
}

// func TestPerformance(t *testing.T) {
// 	conn, _ := Connect("localhost:33013")
// 	space := conn.Space(0)
// 	ch := make(chan *TupleResponse)
// 	for i := 0; i < 100000; i++ {
// 		tuple := []TupleField{ String("Linda"), Int32(1), String("rider"), Int32(21) }
// 		go goo(space, tuple, ch)
// 	}
// 	for i := 0; i < 100000; i++ {
// 		_ = <- ch
// 	}
// }

// func goo(space *Space, tuple []TupleField, ch chan *TupleResponse) {
// 	res, _ := space.Insert(tuple, true)
// 	ch <- res
// }

func TestUpdate(t *testing.T) {
	conn, _ := Connect("localhost:33013")
	space := conn.Space(0)

	tuple := []TupleField{ String("Linda") }
	field1 := UpdOp{ 2, OpEq, String("dancer") }
	field2 := UpdOp{ 3, OpAdd, Int32(1) }
	res, err := space.Update(tuple, true, field1, field2)

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	var lindaJob String
	var lindaAge Int32
	lindaJob.Unpack(res.Tuples[0].Fields[2])
	lindaAge.Unpack(res.Tuples[0].Fields[3])

	if lindaJob != "dancer" {
		t.Errorf("Error: Linda's job should be replaced with dancer, not %s", lindaJob)
	}

	if lindaAge != 22 {
		t.Errorf("Error: Linda's age should be 22, noy %d", lindaAge)
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