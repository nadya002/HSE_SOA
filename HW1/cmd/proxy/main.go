package main

import (
	"fmt"
	"net"
	"strings"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func main() {
	listener, _ := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("localhost"), Port: 2000}) // открываем слушающий UDP-сокет
	for {
		go handleClient(listener) // обрабатываем запрос клиента
	}
}

func handleClient(conn *net.UDPConn) {
	buf := make([]byte, 128) // буфер для чтения клиентских данных

	readLen, addr, err := conn.ReadFromUDP(buf) // читаем из сокета
	if err != nil {
		fmt.Println(err)
		return
	}

	words := strings.Fields(string(buf[:readLen]))
	if len(words) != 2 {
		conn.WriteToUDP([]byte("Wrong request\n"), addr)
		return
	}
	var res []byte

	if "get_result" == words[0] {
		fmt.Println("aaa")

		if words[1] == "json" {
			fmt.Println("json")
			ServerAddr, err := net.ResolveUDPAddr("udp", "server_json:8080")
			CheckError(err)
			res = get_ans(ServerAddr)

		} else if words[1] == "xml" {
			ServerAddr, err := net.ResolveUDPAddr("udp", "server_xml:8080")
			CheckError(err)
			res = get_ans(ServerAddr)

		} else if words[1] == "msgpack" {
			ServerAddr, err := net.ResolveUDPAddr("udp", "server_msgpack:8080")
			CheckError(err)
			res = get_ans(ServerAddr)

		} else if words[1] == "avro" {
			ServerAddr, err := net.ResolveUDPAddr("udp", "server_avro:8080")
			CheckError(err)
			res = get_ans(ServerAddr)

		} else if words[1] == "yaml" {
			ServerAddr, err := net.ResolveUDPAddr("udp", "server_yaml:8080")
			CheckError(err)
			res = get_ans(ServerAddr)

		} else if words[1] == "protobuf" {
			ServerAddr, err := net.ResolveUDPAddr("udp", "server_protobuf:8080")
			CheckError(err)
			res = get_ans(ServerAddr)

		}
		//fmt.Println(string(res))
		conn.WriteToUDP(res, addr)
	}

}

func get_ans(ServerAddr *net.UDPAddr) []byte {
	//ServerAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:10001")
	//CheckError(err)

	//LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	//CheckError(err)

	Conn, err := net.DialUDP("udp", nil, ServerAddr)
	CheckError(err)

	defer Conn.Close()

	buf := []byte("get_result")
	_, err = Conn.Write(buf)
	if err != nil {
		fmt.Println(err)
	}

	an_buf := make([]byte, 1024)
	le, _, _ := Conn.ReadFromUDP(an_buf)

	return an_buf[:le]

}
