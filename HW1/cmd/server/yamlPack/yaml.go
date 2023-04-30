package yamlPack

import (
	"HW1/cmd/server/structForSer"

	yaml "gopkg.in/yaml.v3"
)

func Serial_yaml(m structForSer.Message) ([]byte, error) {
	xmlText, err := yaml.Marshal(&m)
	if err != nil {
		return nil, err
	}
	return xmlText, nil
}

func Deserial_yaml(b []byte) (structForSer.Message, error) {
	var m2 structForSer.Message
	er := yaml.Unmarshal(b, &m2)
	if er != nil {
		return structForSer.Message{}, er
	}
	return m2, nil
}
