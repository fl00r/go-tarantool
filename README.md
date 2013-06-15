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
	name      tarantool.String  // index 0
	id        tarantool.Int32   // index 1
	job       tarantool.String
	age       tarantool.Int8
	bestYears []tarantool.Int8
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