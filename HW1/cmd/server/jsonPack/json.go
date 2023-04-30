package jsonPack

import (
	"HW1/cmd/server/structForSer"
	"encoding/json"
)

func Serial_json(m structForSer.Message) ([]byte, error) {
	b, er := json.Marshal(m)
	if er != nil {
		return nil, er
	}
	return b, nil
}

func Deserial_json(b []byte) (structForSer.Message, error) {
	var m_deck structForSer.Message
	er := json.Unmarshal(b, &m_deck)
	if er != nil {
		return structForSer.Message{}, er
	}
	return m_deck, nil
}
