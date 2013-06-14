package tarantool

import (
	"fmt"
	"github.com/fl00r/go-iproto"
	"bytes"
	"encoding/binary"
	// "errors"
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

type Int32 int32

type Int8 int8

type String string

type SelectKey interface {
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
	err = binary.Write(buffer, binary.LittleEndian, val)
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

func (space *Space) Select(indexNo, offset, limit int32, typeToReturn TypeToReturn, keys ... []SelectKey) (tuples *TupleResponse, err error) {

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

	response, err := space.conn.Request(17, body)
	if err != nil {
		return
	}

	var returnCode int32
	err = binary.Read(response.Body, binary.LittleEndian, &returnCode)
	if err != nil {
		return
	}

	if returnCode != 0 {
		err = fmt.Errorf("Return code is not 0, but %d; Error message: %s", returnCode, response.Body.String())
		return
	}

	var (
		tuplesCount int32
		tuplesSize  int32
		cardinality int32
	)
	err = binary.Read(response.Body, binary.LittleEndian, &tuplesCount)
	if err != nil {
		return
	}
	tuples = &TupleResponse{ tuplesCount, make([]Tuple, tuplesCount) }

	for i := int32(0); i < tuplesCount; i++ {
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
			var size uint64
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
// func (space *Space) Insert()