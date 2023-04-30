package main

import (
	"fmt"
	"net"
	"os"

	"HW1/cmd/server/testProtocols"
)

func main() {
	str := os.Args[1]
	//var str string

	//flag.IntVar(&a, "a", 0, "The integer param")
	//flag.StringVar(&str, "format", "json", "The name of a format")

	//flag.Parse()

	fmt.Println(str)
	//fmt.Println("dump:", dump)

	listener, _ := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("localhost"), Port: 8080}) // открываем слушающий UDP-сокет
	for {
		handleClient(listener, str) // обрабатываем запрос клиента
	}
}

func handleClient(conn *net.UDPConn, str string) {
	buf := make([]byte, 128) // буфер для чтения клиентских данных

	readLen, addr, err := conn.ReadFromUDP(buf) // читаем из сокета
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(buf[:readLen]))
	if "get_result" == string(buf[:readLen]) || "get_result" == string(buf[:readLen-1]) {
		if str == "json" {
			an, er := testProtocols.Test_json()
			if er != nil {
				fmt.Println("error ", er)
			} else {
				res := fmt.Sprintln("json", "-", an.Mem, "-", an.TimeOfSer, "-", an.TimeOfDes)
				conn.WriteToUDP([]byte(res), addr)
			}
		} else if str == "xml" {
			an, er := testProtocols.Test_xml()
			if er != nil {
				fmt.Println("error ", er)
			} else {
				res := fmt.Sprintln("xml", "-", an.Mem, "-", an.TimeOfSer, "-", an.TimeOfDes)
				conn.WriteToUDP([]byte(res), addr)
			}

		} else if str == "msgpack" {
			an, er := testProtocols.Test_msgpack()
			if er != nil {
				fmt.Println("error ", er)
			} else {
				res := fmt.Sprintln("xml", "-", an.Mem, "-", an.TimeOfSer, "-", an.TimeOfDes)
				conn.WriteToUDP([]byte(res), addr)
			}

		} else if str == "avro" {
			an, er := testProtocols.Test_avro()
			if er != nil {
				fmt.Println("error ", er)
			} else {
				res := fmt.Sprintln("xml", "-", an.Mem, "-", an.TimeOfSer, "-", an.TimeOfDes)
				conn.WriteToUDP([]byte(res), addr)
			}

		} else if str == "yaml" {
			an, er := testProtocols.Test_yaml()
			if er != nil {
				fmt.Println("error ", er)
			} else {
				res := fmt.Sprintln("xml", "-", an.Mem, "-", an.TimeOfSer, "-", an.TimeOfDes)
				conn.WriteToUDP([]byte(res), addr)
			}

		} else if str == "protobuf" {
			an, er := testProtocols.Test_protobuf()
			if er != nil {
				fmt.Println("error ", er)
			} else {
				res := fmt.Sprintln("xml", "-", an.Mem, "-", an.TimeOfSer, "-", an.TimeOfDes)
				conn.WriteToUDP([]byte(res), addr)
			}

		} else {
			conn.WriteToUDP([]byte("No such format "+str), addr) // пишем в сокет
		}
	} else {
		conn.WriteToUDP(append([]byte("Wrong querry, you said: "), buf[:readLen]...), addr) // пишем в сокет
	}
}
