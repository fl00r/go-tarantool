// Tarantool protocol https://github.com/mailru/tarantool/blob/master/doc/box-protocol.txt

package tarantool

import (
	"fmt"
	"github.com/fl00r/go-iproto"
	"bytes"
	"encoding/binary"
)

const (
	// Ops
	SelectOp = 17
	InsertOp = 13
	UpdateOp = 19
	DeleteOp = 21
	CallOp   = 22
	PingOp   = 65280

	// Flags
	BoxFlags       = int32(0x00)
	BoxReturnTuple = int32(0x01)
	BoxAdd         = int32(0x02)
	BoxReplace     = int32(0x04)

	// Update Ops
	OpEq     = int8(0)
	OpAdd    = int8(1)
	OpAnd    = int8(2)
	OpXor    = int8(3)
	OpOr     = int8(4)
	OpSplice = int8(5)
	OpDelete = int8(6)
	OpInsert = int8(7)
)

type Space struct {
	spaceNo int32
	conn    *iproto.IProto
}

type Connection struct {
	conn *iproto.IProto
}

type SelectRequestBody struct {
	spaceNo int32
	indexNo int32
	offset  int32
	limit   int32
	count   int32
}

type TupleResponse struct {
	Count  int32
	Tuples []Tuple
}

type Tuple struct {
	Fields [][]byte
}

type UpdOp struct {
	FieldNo int32
	OpCode  int8
	Field   TupleField
}

type Int32 int32

type Int8 int8

type String string

type TupleField interface {
	Pack(*bytes.Buffer) error
}

type TypeToReturn interface {
	Unpack([][]byte) error
}

func (val Int32) Pack(buffer *bytes.Buffer) (err error) {
	buf := make([]byte, 1)
	binary.PutUvarint(buf, uint64(4))
	_, err = buffer.Write(buf)
	if err != nil {
		return
	}
	err = binary.Write(buffer, binary.LittleEndian, val)
	return
}

func (val Int8) Pack(buffer *bytes.Buffer) (err error) {
	buf := make([]byte, 1)
	binary.PutUvarint(buf, uint64(1))
	_, err = buffer.Write(buf)
	if err != nil {
		return
	}
	err = binary.Write(buffer, binary.LittleEndian, val)
	return
}

func (val String) Pack(buffer *bytes.Buffer) (err error) {
	size := len(val)
	buf := make([]byte, 8)
	l := binary.PutUvarint(buf, uint64(size))
	_, err = buffer.Write(buf[0:l])
	if err != nil {
		return
	}
	_, err = buffer.Write([]byte(val))
	return
}

func (val *Int32) Unpack(packet []byte) (err error) {
	err = binary.Read(bytes.NewBuffer(packet), binary.LittleEndian, val)
	return 
}

func (val *Int8) Unpack(packet []byte) (err error) {
	err = binary.Read(bytes.NewBuffer(packet), binary.LittleEndian, val)
	return 
}

func (val *String) Unpack(packet []byte) (err error) {
	*val = String(bytes.NewBuffer(packet).String())
	return 
}


func Connect(addr string) (conn *Connection, err error) {
	ipr, err := iproto.Connect(addr)
	conn = &Connection{ ipr }
	return
}

func (conn *Connection) Space(spaceNo int32) (space *Space) {
	space = &Space{ spaceNo, conn.conn }
	return
}

func (space *Space) Select(indexNo, offset, limit int32, keys ... []TupleField) (tuples *TupleResponse, err error) {

	body := new(bytes.Buffer)

	count := int32(len(keys))
	requestBody := &SelectRequestBody{ space.spaceNo, indexNo, offset, limit, count }
	err = binary.Write(body, binary.LittleEndian, requestBody)
	if err != nil {
		return
	}

	for _, key := range keys {
		err = binary.Write(body, binary.LittleEndian, int32(len(key)))
		if err != nil {
			return
		}
		for _, field := range key {
			field.Pack(body)
		}
	}

	tuples, err = space.request(SelectOp, body)
	return
}

func (space *Space) Insert(tuple []TupleField, returnTuple bool) (tuples *TupleResponse, err error) {
	flags := BoxFlags
	tuples, err = space.insert(flags, returnTuple, tuple)
	return
}

func (space *Space) Add(tuple []TupleField, returnTuple bool) (tuples *TupleResponse, err error) {
	flags := BoxAdd
	tuples, err = space.insert(flags, returnTuple, tuple)
	return

}

func (space *Space) Replace(tuple []TupleField, returnTuple bool) (tuples *TupleResponse, err error) {
	flags := BoxReplace
	tuples, err = space.insert(flags, returnTuple, tuple)
	return

}

func (space *Space) insert(flags int32, returnTuple bool, tuple []TupleField) (tuples *TupleResponse, err error) {
	body := new(bytes.Buffer)

	if returnTuple == true {
		flags |= BoxReturnTuple
	}

	requestBody := []int32{ space.spaceNo, flags }
	err = binary.Write(body, binary.LittleEndian, requestBody)
	if err != nil {
		return
	}

	err = binary.Write(body, binary.LittleEndian, int32(len(tuple)))
	if err != nil {
		return
	}
	for _, field := range tuple {
		field.Pack(body)
	}
	tuples, err = space.request(InsertOp, body)
	return
}

func (space *Space) Update(tuple []TupleField, returnTuple bool, ops ... UpdOp) (tuples *TupleResponse, err error) {
	flags := BoxFlags
	tuples, err = space.update(flags, returnTuple, tuple, ops)
	return
}

func (space *Space) Upsert() {

}

func (space *Space) update(flags int32, returnTuple bool, tuple []TupleField, ops []UpdOp) (tuples *TupleResponse, err error) {
	body := new(bytes.Buffer)

	if returnTuple == true {
		flags |= BoxReturnTuple
	}

	requestBody := []int32{ space.spaceNo, flags }
	err = binary.Write(body, binary.LittleEndian, requestBody)
	if err != nil {
		return
	}

	err = binary.Write(body, binary.LittleEndian, int32(len(tuple)))
	if err != nil {
		return
	}

	for _, field := range tuple {
		field.Pack(body)
	}

	opsCount := int32(len(ops))
	err = binary.Write(body, binary.LittleEndian, opsCount)
	if err != nil {
		return
	}

	for _, op := range ops {
		err = binary.Write(body, binary.LittleEndian, op.FieldNo)
		if err != nil {
			return
		}
		err = binary.Write(body, binary.LittleEndian, op.OpCode)
		if err != nil {
			return
		}
		err = op.Field.Pack(body)
		if err != nil {
			return
		}
	}

	tuples, err = space.request(UpdateOp, body)
	return
}

// Refactor: same as Insert but Op number
func (space *Space) Delete(tuple []TupleField, returnTuple bool) (tuples *TupleResponse, err error) {
	body := new(bytes.Buffer)
	flags := BoxFlags

	if returnTuple == true {
		flags |= BoxReturnTuple
	}

	requestBody := []int32{ space.spaceNo, flags }
	err = binary.Write(body, binary.LittleEndian, requestBody)
	if err != nil {
		return
	}

	err = binary.Write(body, binary.LittleEndian, int32(len(tuple)))
	if err != nil {
		return
	}

	for _, field := range tuple {
		field.Pack(body)
	}

	tuples, err = space.request(DeleteOp, body)
	return
}

func (space *Space) Call(procName string, returnTuple bool, args []TupleField) (tuples *TupleResponse, err error) {
	body := new(bytes.Buffer)
	flags := BoxFlags

	if returnTuple == true {
		flags |= BoxReturnTuple
	}

	err = binary.Write(body, binary.LittleEndian, flags)
	if err != nil {
		return
	}

	err = String(procName).Pack(body)
	if err != nil {
		return
	}

	err = binary.Write(body, binary.LittleEndian, int32(len(args)))
	if err != nil {
		return
	}

	for _, field := range args {
		field.Pack(body)
	}

	tuples, err = space.request(CallOp, body)
	return
}

func (space *Space) Ping() {

}

func (space *Space) request(requestId int32, body *bytes.Buffer) (tuples *TupleResponse, err error) {
	var (
		returnCode  int32
		tuplesCount int32
		tuplesSize  int32
		cardinality int32
		size        uint64
		response    *iproto.Response
	)

	response, err = space.conn.Request(requestId, body)
	if err != nil {
		return
	}

	err = binary.Read(response.Body, binary.LittleEndian, &returnCode)
	if err != nil {
		return
	}

	if returnCode != 0 {
		err = fmt.Errorf("Return code is not 0, but %d; Error message: %s", returnCode, response.Body.String())
		return
	}
	err = binary.Read(response.Body, binary.LittleEndian, &tuplesCount)
	if err != nil {
		return
	}
	tuples = &TupleResponse{ tuplesCount, make([]Tuple, tuplesCount) }

	for i := int32(0); i < tuplesCount && response.Body.Len() > 0; i++ {
		err = binary.Read(response.Body, binary.LittleEndian, &tuplesSize)
		if err != nil {
			return
		}
		err = binary.Read(response.Body, binary.LittleEndian, &cardinality)
		if err != nil {
			return
		}
		tuples.Tuples[i] = Tuple{ make([][]byte, cardinality) }
		for j := int32(0); j < cardinality; j++ {
			size, err = binary.ReadUvarint(response.Body)
			if err != nil {
				return
			}
			tuples.Tuples[i].Fields[j] = make([]byte, size)
			_, err = response.Body.Read(tuples.Tuples[i].Fields[j])
			if err != nil {
				return
			}
		}
	}
	return
}