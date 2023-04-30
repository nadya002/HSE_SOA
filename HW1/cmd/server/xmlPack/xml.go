package xmlPack

import (
	"HW1/cmd/server/structForSer"
	"encoding/xml"
)

func Serial_xml(m structForSer.Message) ([]byte, error) {
	xmlText, err := xml.MarshalIndent(&m, " ", " ")
	if err != nil {
		return nil, err
	}
	return xmlText, nil
}

func Deserial_xml(b []byte) (structForSer.Message, error) {
	var m2 structForSer.Message
	er := xml.Unmarshal(b, &m2)
	if er != nil {
		return structForSer.Message{}, er
	}
	return m2, nil
}
