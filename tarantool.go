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

type TupleResponse struct {
	Count  int32
	Tuples []Tuple
}

type Tuple struct {
	Fields [][]byte
}

type SelectRequestBody struct {
	spaceNo int32
	indexNo int32
	offset  int32
	limit   int32
	count   int32
}

type SelectField struct {
	Cardinality int32
	Fields		*bytes.Buffer
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

func (space *Space) Select(indexNo, offset, limit int32, keys ... *SelectField) (tuples *TupleResponse, err error) {

	body := new(bytes.Buffer)

	count := int32(len(keys))
	requestBody := &SelectRequestBody{ space.spaceNo, indexNo, offset, limit, count }
	err = binary.Write(body, binary.LittleEndian, requestBody)
	if err != nil {
		return
	}

	for _, key := range keys {
		binary.Write(body, binary.LittleEndian, key.Cardinality)
		fieldLength := key.Fields.Len()
		buf := make([]byte, 8)
		l := binary.PutUvarint(buf, uint64(fieldLength))
		_, err = body.Write(buf[0:l])
		if err != nil {
			return
		}
		_, err = body.Write(key.Fields.Bytes())
		if err != nil {
			return
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
	err = binary.Read(response.Body, binary.LittleEndian, &tuplesSize)
	if err != nil {
		return
	}
	tuples = &TupleResponse{ tuplesCount, make([]Tuple, tuplesCount) }

	for i := int32(0); i < tuplesCount; i++ {
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

func (tuple *Tuple) PackTuple(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.LittleEndian, len(tuple.Fields))

	for _, key := range tuple.Fields {
		fieldLength := key.Fields.Len()
		buf := make([]byte, 8)
		l := binary.PutUvarint(buf, uint64(fieldLength))
		_, err = body.Write(buf[0:l])
		if err != nil {
			return
		}
		_, err = body.Write(key.Fields.Bytes())
		if err != nil {
			return
		}
	}

}

// func (space *Space) Insert()