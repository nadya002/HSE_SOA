package msgpackPack

import (
	"HW1/cmd/server/structForSer"

	gomsgpack "github.com/nnabeyang/go-msgpack"
)

func Serial_msgpack(m structForSer.Message) ([]byte, error) {
	b, er := gomsgpack.Marshal(m, false)

	if er != nil {
		return nil, er
	}
	return b, nil
}

func Deserial_msgpack(b []byte) (structForSer.Message, error) {
	var m_deck structForSer.Message
	er := gomsgpack.Unmarshal(b, &m_deck)
	if er != nil {
		return structForSer.Message{}, er
	}
	return m_deck, nil
}
