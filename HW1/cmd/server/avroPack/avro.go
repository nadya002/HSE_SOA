package avroPack

import (
	"HW1/cmd/server/structForSer"
	"fmt"

	"github.com/hamba/avro"
)

func Serial_avro(m structForSer.Message) ([]byte, error) {
	schema, er := avro.Parse(`{
		"type": "record",
		"name": "simple",
		"namespace": "org.hamba.avro",
		"fields" : [
			{"name": "str", "type": "string"},
			{"name": "numb", "type": "long"},
			{"name": "m", "type" : {"type": "map", "values": "string"}},
			{"name": "arr", "type" : {"type" : "array", "items" : "string"}},
			{"name": "float_numb", "type": "float"}

		]
	}`)

	if er != nil {
		fmt.Println(er)
	}

	xmlText, err := avro.Marshal(schema, &m)
	if err != nil {
		return nil, err
	}
	return xmlText, nil

}

func Deserial_avro(b []byte) (structForSer.Message, error) {
	schema, _ := avro.Parse(`{
		"type": "record",
		"name": "simple",
		"namespace": "org.hamba.avro",
		"fields" : [
			{"name": "str", "type": "string"},
			{"name": "numb", "type": "long"},
			{"name": "m", "type" : {"type": "map", "values": "string"}},
			{"name": "arr", "type" : {"type" : "array", "items" : "string"}},
			{"name": "float_numb", "type": "float"}

		]
	}`)

	var m2 structForSer.Message
	er := avro.Unmarshal(schema, b, &m2)
	if er != nil {
		return structForSer.Message{}, er
	}
	return m2, nil
}
