package prepareData

import (
	"HW1/cmd/server/pb"
	"HW1/cmd/server/structForSer"
	"math/rand"
)

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

func Prepare_data() structForSer.Message {
	ar := make_array()
	m := make_map()
	//m := make(map[string]string)
	mes := structForSer.Message{Str: "Alice",
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
