package testProtocols

import (
	"HW1/cmd/dto"
	"HW1/cmd/server/pb"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/hamba/avro"
	yaml "gopkg.in/yaml.v3"

	"github.com/golang/protobuf/proto"

	gomsgpack "github.com/nnabeyang/go-msgpack"
)

// StringMap is a map[string]string.
type Map map[string]string

type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

// MarshalXML marshals the map to XML, with each key in the map being a
// tag and it's corresponding value being it's contents.
func (m Map) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(m) == 0 {
		return nil
	}

	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	for k, v := range m {
		e.Encode(xmlMapEntry{XMLName: xml.Name{Local: k}, Value: v})
	}

	return e.EncodeToken(start.End())
}

// UnmarshalXML unmarshals the XML into a map of string to strings,
// creating a key in the map for each tag and setting it's value to the
// tags contents.
//
// The fact this function is on the pointer of Map is important, so that
// if m is nil it can be initialized, which is often the case if m is
// nested in another xml structurel. This is also why the first thing done
// on the first line is initialize it.
func (m *Map) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = Map{}
	for {
		var e xmlMapEntry

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		(*m)[e.XMLName.Local] = e.Value
	}
	return nil
}

// type Message struct {
// 	Str        string
// 	Numb       int64
// 	M          map[string]string
// 	Arr        []string
// 	Float_numb float64
// }

type SimpleRecord struct {
	A int64  `avro:"a"`
	B string `avro:"b"`
}

//var schema avro.Schema

type Message struct {
	//XMLName    xml.Name `xml:"mess_for_xml"`
	Str        string   `xml:"str,attr" avro:"str" yaml:"str"`
	Numb       int64    `xml:"numb" avro:"numb" yaml:"numb"`
	M          Map      `xml:"m" avro:"m" yaml:"m"`
	Arr        []string `xml:"arr" avro:"arr" yaml:"arr"`
	Float_numb float32  `xml:"float_numb" avro:"float_numb" yaml:"float_numb"`
}

// type BreakfastMenu struct {
// 	XMLName xml.Name `xml:"test_xmk"`
// 	Food    []struct {
// 		Name        string `xml:"name"`
// 		Price       string `xml:"price"`
// 		Description string `xml:"description"`
// 		Calories    string `xml:"calories"`
// 	} `xml:"food"`
// }

func Serial_msgpack(m Message) ([]byte, error) {
	b, er := gomsgpack.Marshal(m, false)

	if er != nil {
		return nil, er
	}
	return b, nil
}

func Deserial_msgpack(b []byte) (Message, error) {
	var m_deck Message
	er := gomsgpack.Unmarshal(b, &m_deck)
	if er != nil {
		return Message{}, er
	}
	return m_deck, nil
}

func Serial_json(m Message) ([]byte, error) {
	b, er := json.Marshal(m)
	if er != nil {
		return nil, er
	}
	return b, nil
}

func Deserial_json(b []byte) (Message, error) {
	var m_deck Message
	er := json.Unmarshal(b, &m_deck)
	if er != nil {
		return Message{}, er
	}
	return m_deck, nil
}

func Serial_xml(m Message) ([]byte, error) {
	xmlText, err := xml.MarshalIndent(&m, " ", " ")
	if err != nil {
		return nil, err
	}
	return xmlText, nil
}

func Deserial_xml(b []byte) (Message, error) {
	var m2 Message
	er := xml.Unmarshal(b, &m2)
	if er != nil {
		return Message{}, er
	}
	return m2, nil
}

func Serial_yaml(m Message) ([]byte, error) {
	xmlText, err := yaml.Marshal(&m)
	if err != nil {
		return nil, err
	}
	return xmlText, nil
}

func Deserial_yaml(b []byte) (Message, error) {
	var m2 Message
	er := yaml.Unmarshal(b, &m2)
	if er != nil {
		return Message{}, er
	}
	return m2, nil
}

// func To_proto_mes(m Message) pb.ProtoMessage {

// }

func Serial_protobuf(m pb.ProtoMessage) ([]byte, error) {
	//fmt.Println(m)
	p, err := proto.Marshal(&m)
	if err != nil {
		fmt.Println("marshaling error: ", err)
		return nil, nil
	}
	//fmt.Println(p)
	return p, nil
}

func Deserial_protobuf(b []byte) (pb.ProtoMessage, error) {
	msg := pb.ProtoMessage{}
	proto.Unmarshal(b, &msg)
	return msg, nil
}

func Serial_avro(m Message) ([]byte, error) {
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

func Deserial_avro(b []byte) (Message, error) {
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

	var m2 Message
	er := avro.Unmarshal(schema, b, &m2)
	if er != nil {
		return Message{}, er
	}
	return m2, nil
}

func make_array() []string {
	var ar []string
	for i := 0; i < 2000; i++ {
		var st string
		for j := 0; j < 10; j++ {

			c := 'a' + rune(rand.Intn('z'-'a'+1))
			st += string(c)
		}
		ar = append(ar, st)
	}
	return ar
}

func make_map() map[string]string {
	m := make(map[string]string)
	for i := 0; i < 2000; i++ {
		var key string
		for j := 0; j < 10; j++ {
			c := 'a' + rune(rand.Intn('z'-'a'+1))
			key += string(c)
		}

		var val string
		for j := 0; j < 10; j++ {
			c := 'a' + rune(rand.Intn('z'-'a'+1))
			val += string(c)
		}
		m[key] = val
	}
	return m
}

func Prepare_data() Message {
	ar := make_array()
	m := make_map()
	//m := make(map[string]string)
	mes := Message{Str: "Alice",
		Numb:       1264739,
		M:          m,
		Arr:        ar,
		Float_numb: 23.45,
	}
	return mes
}

func Prepare_data_for_protobuf() pb.ProtoMessage {
	ar := make_array()
	m := make_map()
	//m := make(map[string]string)
	mes := pb.ProtoMessage{
		Str:       "Alice",
		Numb:      1264739,
		Map:       m,
		Arr:       ar,
		FloatNumb: 23.45,
	}
	return mes
}

func Test_method(Des func(b []byte) (Message, error), Ser func(m Message) ([]byte, error)) (dto.Answer, error) {
	rand.Seed(123)

	mes := Prepare_data()

	start := time.Now()
	var byt []byte
	var er error
	for i := 0; i < 1000; i++ {
		byt, er = Ser(mes)
		if er != nil {
			return dto.Answer{}, er
		}
	}
	duration := time.Since(start) / 1000

	sz_of_str := len(byt)

	start2 := time.Now()

	for i := 0; i < 1000; i++ {
		_, err := Des(byt)
		if err != nil {
			return dto.Answer{}, er
		}
	}
	duration2 := time.Since(start2) / 1000

	return dto.Answer{duration, duration2, sz_of_str}, nil

}

func Test_json() (dto.Answer, error) {
	fun_des_json := Deserial_json
	fun_ser_json := Serial_json
	return Test_method(fun_des_json, fun_ser_json)

}

func Test_xml() (dto.Answer, error) {
	fun_des_xml := Deserial_xml
	fun_ser_xml := Serial_xml
	return Test_method(fun_des_xml, fun_ser_xml)
}

func Test_msgpack() (dto.Answer, error) {
	fun_des_msgpack := Deserial_msgpack
	fun_ser_msgpack := Serial_msgpack
	return Test_method(fun_des_msgpack, fun_ser_msgpack)
}

func Test_avro() (dto.Answer, error) {
	fun_des_avro := Deserial_avro
	fun_ser_avro := Serial_avro
	return Test_method(fun_des_avro, fun_ser_avro)
}

func Test_yaml() (dto.Answer, error) {
	fun_des_yaml := Deserial_yaml
	fun_ser_yaml := Serial_yaml
	return Test_method(fun_des_yaml, fun_ser_yaml)
}

func Test_protobuf() (dto.Answer, error) {
	rand.Seed(123)

	mes := Prepare_data_for_protobuf()

	start := time.Now()
	var byt []byte
	var er error
	for i := 0; i < 1000; i++ {
		byt, er = Serial_protobuf(mes)
		if er != nil {
			return dto.Answer{}, er
		}
	}
	duration := time.Since(start) / 1000

	sz_of_str := len(byt)

	start2 := time.Now()

	for i := 0; i < 1000; i++ {
		_, err := Deserial_protobuf(byt)
		if err != nil {
			return dto.Answer{}, er
		}
	}
	duration2 := time.Since(start2) / 1000

	return dto.Answer{duration, duration2, sz_of_str}, nil
}
