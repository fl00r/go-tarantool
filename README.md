# Go, Tarantool, Go

[Tarantool](http://tarantool.org/) client on Go.

## Usage

You have got space (Employee) with two indexes: `id` (field 0) and `name` (field 1).

Currently you should explicitly define how to pack your query tuples and how to unpack returned tuple.

So you need to define `Employee` struct, `SelectIdTuple` struct and `SelectNameTuple`.

```go
package main

import (
	"github.com/fl00r/go-tarantool"
	"encoding/binary"
	"bytes"
)

type Employee struct {
	id   int32   // index 0
	name string  // index 1
	job  string
	age  int8
}

type SelectIdTuple struct {
	id int32
}

type SelectNameTuple struct {
	name string
}

func (tuple *SelectTuple) pack(buffer *bytes.Buffer) (n int, err error) {
	<!-- tarantool.PackInt32(tuple.id) -->
	n, err = buffer.Write(buffer.NewBuffer(tuple.Name).Bytes())
	return
}

func (tuple *SelectNameTuple) pack(buffer *bytes.Buffer) (n int, err error) {
	n, err = binary.Write(buffer, binary.LittleEndian, tuple.id)
	return
}

func (tuple Employee) unpack(buffer) (tuple *Employee) {
	
}



func main() {
	connection = tarantool.Connect("loaclhost:33013")
	space = connection.Space(1)

	//INSERT
	tuple = MyTuple{ 1 }
	// box.select(0, 0, 1)
	// indexNo, offset, limit, tuple
	space.Insert(0, 0, 1, &SelectIdTuple{ 1, 2, 3, 4 }, &SelectIdTuple{ 2 })
	// space.Select()
	// space.Delete()
	// space.Update()
}
```

TO BE WRITTEN