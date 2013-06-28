package main

import (
	"github.com/fl00r/go-tarantool"
	// "bytes"
	// "encoding/binary"
	"fmt"
)

// All error handling is skipped
func main() {
	//
	// ESTABLISHING CONNECTION
	//

	// Connectiong to host:port
	connection, _ := tarantool.Connect("localhost:33013")
	// Selecting space
	space := connection.Space(0)

	//
	// INSERTING DATA
	//

	// Prepare tuple
	// Tuple is an array of TupleField interfaces
	//
	//   type TupleField interface {
	//   	Pack(*bytes.Buffer) error
	//   }
	//
	// There are some internaly implemented types, such as
	//   * tarantool.Int8
	//   * tarantool.Int32
	//   * tarantool.String
	// This types are just helpers to pack data in terms of tarantool specification.
	// You are free to add your own types (varint, int64, BER, json, etc.).
	//
	// Let's imaginge you have got following tuple structure:
	//   { id(int32), name(string), age(int8), job(string) }
	tuple := []tarantool.TupleField{
		tarantool.Int32(1),
		tarantool.String("Peter"),
		tarantool.Int8(18),
		tarantool.String("janitor"),
	}
	// space.Insert accepts tuple and returnTuple flag/
	// If returnTuple is true, you will receive your tuple back after insertion/
	// Insert returns result and error.
	// Result is an array of arrays of arrays X)
	//   * each array item represents a Tuple (array)
	//   * each item of Tuple array represent a field (array)
	//   * each field array represents a byte
	res, _ := space.Insert(tuple, true)

	fmt.Println(res)
	//=> [
	//=>   [ [1 0 0 0] [80 101 116 101 114] [18] [106 97 110 105 116 111 114] ]
	//=> ]
	// Here is 1 row (tuple) returned with four fields in it
	// each field reasonably has got different bytes amount
	// so you need to unpack the, properly accordingly to yours data scheme

	//
	// ADDING AND REPLACING DATA
	//

	// Insert is conventionaly ADD/REPLACE operation
	// It means that if tuple is already exists it will *Replace* it.
	// If tuple doesn't exist it will *Add* it.
	// If you want to Add data explicitly:
	res, err := space.Add(tuple, true)

	// So if you will try to add existing tuple you will get an error
	fmt.Println(err)
	//=> Return code is not 0, but 14082; Error message: Duplicate key exists in unique index 0

	addTuple := []tarantool.TupleField{
		tarantool.Int32(2),
		tarantool.String("Mary"),
		tarantool.Int8(20),
		tarantool.String("singer"),
	}
	res, err = space.Add(addTuple, true)

	fmt.Println(res)
	//=> [
	//=>   [ [2 0 0 0] [77 97 114 121] [20] [115 105 110 103 101 114] ]
	//=> ]
	// Good job!

	// To replace data, use `space.Replace` function
	// It will return an error if current tuple doen't exist.
	// And it will replace existing tuple.
	replaceTuple := []tarantool.TupleField{
		tarantool.Int32(2),
		tarantool.String("Mary"),
		tarantool.Int8(20),
		tarantool.String("musician"),
	}
	space.Replace(replaceTuple, true)

	//
	// SELECTING DATA
	//

	// Prepare select tuple. Select tuple includes indexed fields.
	// space.Select accepts:
	//   * indexNo - int32, number of index
	//   * offset  - int32, amount of documents to skip
	//   * limit   - int32, count of documents to return
	//   * keys    - ... []tarantool.TupleField, select tuples
	// You could pass as many keys as you need.
	// Assuming that you have got two index fields: id and name
	key1 := []tarantool.TupleField{tarantool.String("Peter")}
	key2 := []tarantool.TupleField{tarantool.String("Mary")}
	indexNo := int32(1) // Name is our first index, id is zero
	offset := int32(0)
	limit := int32(100)

	res, _ = space.Select(indexNo, offset, limit, key1, key2)

	fmt.Println(res)
	//=> [
	//=>   [ [1 0 0 0] [80 101 116 101 114] [18] [106 97 110 105 116 111 114] ]
	//=>   [ [2 0 0 0] [77 97 114 121] [20] [109 117 115 105 99 105 97 110] ]
	//=> ]
	// You get back two tuples you've inserted earlier.

	//
	// UPDATING DATA
	//

	// Updating data in tarantool means updating a field or a number of fields.
	// There are a bunch of operations available:
	//   * tarantool.OpEq
	//   * tarantool.OpAdd
	//   * tarantool.OpAnd
	//   * tarantool.OpXor
	//   * tarantool.OpOr
	//   * tarantool.OpSplice
	//   * tarantool.OpDelete
	//   * tarantool.OpInsert
	// For more information check tarantool.org documentation.

	// So for space.Update you should specify following arguments:
	//   * key         - []tarantool.TupleField, tuple to update (primary index key)
	//   * returnTuple - bool, returning tuple after operation
	//   * fields      - ... tarantool.UpdOp
	// tarantool.UpdOp is a following structure:
	//   * FieldNo - int32, number of a field to apply operation
	//   * OpCode  - int8, operation code
	//   * Field   - tarantool.TupleField, argument to op

	// Let's fetch Mary by her primary index (id, 2)
	updKey := []tarantool.TupleField{tarantool.Int32(2)}
	// And increase  Mary's age (field 2)
	updField2 := tarantool.UpdOp {
		2,
		tarantool.OpEq,
		tarantool.Int8(21),
	}
	// And change her job (field 3)
	updField1 := tarantool.UpdOp {
		3,
		tarantool.OpEq,
		tarantool.String("guitarist"),
	}
	res, _ = space.Update(updKey, true, updField1, updField2)

	fmt.Println("error", res)
	//=> [
	//=>   [ [2 0 0 0] [77 97 114 121] [21] [103 117 105 116 97 114 105 115 116] ]
	//=> ]

	//
	// CALL (lua procedure)
	//

	// space.Call accepts following arguments:
	//   * procName    - string, procedure name
	//   * returnTuple - bool, returning tuple after operation
	//   * args        - ... tarantool.TupleField, arguments
	// For example, `box.select_range` accepts (spaceNo, indexNo, limit and index offset)
	// Args should be passed in its string representation.
	callSpaceNo := tarantool.String("0")
	callIndexNo := tarantool.String("0")
	callLimit   := tarantool.String("10")
	// As far as this particular value is stored as int32 bit representation
	startAt     := tarantool.Int32(2)
	res, _ = space.Call("box.select_range", true, callSpaceNo, callIndexNo, callLimit, startAt)

	fmt.Println(res)
	//=> [
	//=>   [ [2 0 0 0] [77 97 114 121] [21] [103 117 105 116 97 114 105 115 116] ]
	//=> ]
	// So we passed index offset to 2, so we get only one row with id 2, because 1 < 2

	//
	// DELETE
	//

	// space.Delete accepts:
	//   * key         - []tarantool.TupleField, primary key to delete
	//   * returnTuple - bool, returning tuple after operation

	// Let's delete all rows we create
	deleteKey1 := []tarantool.TupleField{tarantool.Int32(1)}
	deleteKey2 := []tarantool.TupleField{tarantool.Int32(2)}

	space.Delete(deleteKey1, false)
	space.Delete(deleteKey2, false)

	// And let's ensure there is no more tuples:
	res, _ = space.Call("box.dostring", true, tarantool.String("return box.space[0].index[0]:len()"))
	fmt.Println("count", res[0][0])
	//=> [0 0 0 0]
	// Which is a byte representation of 0 in int32

	// Happy hacking!
}
