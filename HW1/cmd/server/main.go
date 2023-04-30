package main

import (
	"fmt"
	"net"
	"os"

	"HW1/cmd/dto"
	"HW1/cmd/server/testProtocols"
)

func handlReq(test_func func() (dto.Answer, error), name string) []byte {
	an, er := test_func()
	if er != nil {
		fmt.Println("error ", er)
		return []byte("error")
	} else {
		res := fmt.Sprintln(name, "-", an.Mem, "-", an.TimeOfSer, "-", an.TimeOfDes)
		return []byte(res)
	}
}

func main() {
	str := os.Args[1]

	listener, _ := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("localhost"), Port: 8080}) // открываем слушающий UDP-сокет
	for {
		handleClient(listener, str) // обрабатываем запрос клиента
	}
}

func handleClient(conn *net.UDPConn, str string) {
	buf := make([]byte, 1024) // буфер для чтения клиентских данных

	readLen, addr, err := conn.ReadFromUDP(buf) // читаем из сокета
	if err != nil {
		fmt.Println(err)
		return
	}

	//fmt.Println(string(buf[:readLen]))
	var res []byte
	if "get_result" == string(buf[:readLen]) || "get_result" == string(buf[:readLen-1]) {
		if str == "json" {
			res = handlReq(testProtocols.Test_json, "json")
		} else if str == "xml" {
			res = handlReq(testProtocols.Test_xml, "xml")
		} else if str == "msgpack" {
			res = handlReq(testProtocols.Test_msgpack, "msgpack")
		} else if str == "avro" {
			res = handlReq(testProtocols.Test_avro, "avro")
		} else if str == "yaml" {
			res = handlReq(testProtocols.Test_yaml, "yaml")
		} else if str == "protobuf" {
			res = handlReq(testProtocols.Test_protobuf, "protobuf")

		} else {
			res = []byte("No such format " + str) // пишем в сокет
		}
	} else {
		conn.WriteToUDP(append([]byte("Wrong req, you said: "), buf[:readLen]...), addr) // пишем в сокет
	}
	conn.WriteToUDP(res, addr)
}
