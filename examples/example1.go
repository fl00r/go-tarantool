package main

import (
	"github.com/fl00r/go-tarantool"
	"bytes"
	"fmt"
)

// Let's define our model
// tarantool.Int32, tarantool.String, tarantool.Int8 are default tarantool types
// You can define any type of data which implements Pack / Unpack method (tarantool.TupleField interface)
type Employee struct {
	Id        tarantool.Int32
	Name      tarantool.String
	Age       tarantool.Int8
	Job       tarantool.string
	BestYears MyCustomType
}

// Our Custom Type should be unpacked to array of ints
type MyCustomType []tarantool.Int32

// Here we pack our custom type into separate fields in tuple, 
// but it is possible to pack it into BER or serialized array into string etc
func (val *MyCustomType) Pack(bytes *bytes.Buffer) (err error) {
	// ...
}

// We are unpacking our 
func (val *MyCustomType) Unpack() (err error) {
	// ...
}

func main() {
	// ...
}