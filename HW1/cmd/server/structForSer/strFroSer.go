package structForSer

import "HW1/cmd/server/mapForXml"

type Message struct {
	//XMLName    xml.Name `xml:"mess_for_xml"`
	Str        string        `xml:"str,attr" avro:"str" yaml:"str"`
	Numb       int64         `xml:"numb" avro:"numb" yaml:"numb"`
	M          mapForXml.Map `xml:"m" avro:"m" yaml:"m"`
	Arr        []string      `xml:"arr" avro:"arr" yaml:"arr"`
	Float_numb float32       `xml:"float_numb" avro:"float_numb" yaml:"float_numb"`
}
