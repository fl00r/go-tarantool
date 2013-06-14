package tarantool

import (
	"testing"
	"fmt"
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

func TestSelect(t *testing.T) {
	var res *TupleResponse

	conn, err := Connect("localhost:33013")
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	space := conn.Space(1)

	tuple := []TupleField{ String("Linda"), Int32(3), String("rider"), Int32(21) }
	res, err = space.Insert(tuple, true)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
		fmt.Println(res)

	// for i := int32(0); i < res.Count; i++ {
	// 	emp := &Employee{}
	// 	err = emp.Unpack(res.Tuples[i].Fields)
	// 	if err != nil {
	// 		t.Errorf("Error: %s", err.Error())
	// 	}
	// 	fmt.Println(emp)
	// }



	key1 := []TupleField{ Int32(1) }
	key2 := []TupleField{ Int32(2) }
	var limit int32 = 10
	res, err = space.Select(1, 0, limit, &Employee{}, key1, key2)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	for i := int32(0); i < res.Count; i++ {
		emp := &Employee{}
		err = emp.Unpack(res.Tuples[i].Fields)
		if err != nil {
			t.Errorf("Error: %s", err.Error())
		}
		fmt.Println(emp)
	}
}