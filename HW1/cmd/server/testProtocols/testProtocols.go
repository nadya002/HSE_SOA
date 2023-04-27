package testProtocols

import (
	"HW1/cmd/dto"
	"encoding/json"
	"encoding/xml"
	"io"
	"math/rand"
	"time"

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

type Message struct {
	//XMLName    xml.Name `xml:"mess_for_xml"`
	Str        string   `xml:"str,attr"`
	Numb       int64    `xml:"numb"`
	M          Map      `xml:"m"`
	Arr        []string `xml:"arr"`
	Float_numb float64  `xml:"float_numb"`
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
	mes := Message{Str: "Alice",
		Numb:       1264739,
		M:          m,
		Arr:        ar,
		Float_numb: 23.45,
	}
	return mes
}

func Test_method(Des func(b []byte) (Message, error), Ser func(m Message) ([]byte, error)) (dto.Answer, error) {
	rand.Seed(time.Now().UnixNano())

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
