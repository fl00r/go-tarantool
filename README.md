# Go, Tarantool, Go

[Tarantool](http://tarantool.org/) client on Go.

## Usage

You have got space (Employee) with two indexes: `id` (field 0) and `name` (field 1).

Currently you should explicitly define how to pack your query tuples and how to unpack returned tuple. Later some Reflections will be added for simplifying API (developer will have got low level and high level APIs).

So you need to define following structs: `Employee`, `SelectId` and `SelectName`. Employee should support `Unpack` method (It will receives raw bytes). `SelectId` and `SelectName` should specify `Pack` method.

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

type SelectId struct {
	cardinality int32
	id          int32
}

type SelectName struct {
	cardinality int32
	name        string
}

func (tuple *SelectId) Pack(buffer *bytes.Buffer) (err error) {
	err = binary.Write(buffer, binary.LittleEndian, tuple.cardinality)
	if err != nil {
		return
	}
	err = binary.Write(buffer, binary.LittleEndian, tuple.id)
	return
}

func (tuple *SelectName) Pack(buffer *bytes.Buffer) (err error) {
	err = binary.Write(buffer, binary.LittleEndian, tuple.cardinality)
	if err != nil {
		return
	}
	err = buffer.Write(bytes.NewBuffer(tuple.name).Bytes())
	return
}

func (tuple Employee) Unpack(bytes [][]byte) err error {
	err = binary.Read(bytes.NewBuffer(bytes[0]), binary.LittleEndian, &tuple.id)
	if err != nil {
		return
	}
	tuple.name = bytes.NewBuffer(bytes[1]).String()
	tuple.job = bytes.NewBuffer(bytes[2]).String()
	err = binary.Read(bytes.NewBuffer(bytes[0]), binary.LittleEndian, &tuple.age)
	return
}

func main() {
	connection = tarantool.Connect("loaclhost:33013")
	space = connection.Space(1)

	// INSERT
	//
	// SELECT
	// indexNo, offset, limit, tuple
	space.Insert(0, 0, 1, &SelectId{ 1, 1 }, &SelectId{ 1, 2 })
	// space.Select()
	// space.Delete()
	// space.Update()
}
```

TO BE WRITTEN