package testProtocols

import (
	"HW1/cmd/dto"
	"HW1/cmd/server/avroPack"
	"HW1/cmd/server/jsonPack"
	"HW1/cmd/server/msgpackPack"
	"HW1/cmd/server/prepareData"
	"HW1/cmd/server/protoPack"
	"HW1/cmd/server/structForSer"
	"HW1/cmd/server/xmlPack"
	"HW1/cmd/server/yamlPack"
	"math/rand"
	"time"
)

func Test_method(Des func(b []byte) (structForSer.Message, error), Ser func(m structForSer.Message) ([]byte, error)) (dto.Answer, error) {
	rand.Seed(123)

	mes := prepareData.Prepare_data()

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
	fun_des_json := jsonPack.Deserial_json
	fun_ser_json := jsonPack.Serial_json
	return Test_method(fun_des_json, fun_ser_json)

}

func Test_xml() (dto.Answer, error) {
	fun_des_xml := xmlPack.Deserial_xml
	fun_ser_xml := xmlPack.Serial_xml
	return Test_method(fun_des_xml, fun_ser_xml)
}

func Test_msgpack() (dto.Answer, error) {
	fun_des_msgpack := msgpackPack.Deserial_msgpack
	fun_ser_msgpack := msgpackPack.Serial_msgpack
	return Test_method(fun_des_msgpack, fun_ser_msgpack)
}

func Test_avro() (dto.Answer, error) {
	fun_des_avro := avroPack.Deserial_avro
	fun_ser_avro := avroPack.Serial_avro
	return Test_method(fun_des_avro, fun_ser_avro)
}

func Test_yaml() (dto.Answer, error) {
	fun_des_yaml := yamlPack.Deserial_yaml
	fun_ser_yaml := yamlPack.Serial_yaml
	return Test_method(fun_des_yaml, fun_ser_yaml)
}

func Test_protobuf() (dto.Answer, error) {
	rand.Seed(123)

	mes := prepareData.Prepare_data_for_protobuf()

	start := time.Now()
	var byt []byte
	var er error
	for i := 0; i < 1000; i++ {
		byt, er = protoPack.Serial_protobuf(mes)
		if er != nil {
			return dto.Answer{}, er
		}
	}
	duration := time.Since(start) / 1000

	sz_of_str := len(byt)

	start2 := time.Now()

	for i := 0; i < 1000; i++ {
		_, err := protoPack.Deserial_protobuf(byt)
		if err != nil {
			return dto.Answer{}, er
		}
	}
	duration2 := time.Since(start2) / 1000

	return dto.Answer{duration, duration2, sz_of_str}, nil
}
