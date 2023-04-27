package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"

	"HW1/cmd/dto"
	"HW1/cmd/server/testProtocols"
)

func main() {

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", host, port))

	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			conn.Close()
			continue
		}
		go handleConnection(conn) // запускаем горутину для обработки запроса
	}
}

// обработка подключения
func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		// считываем полученные в запросе данные
		input := make([]byte, (1024 * 4))
		n, err := conn.Read(input)
		if n == 0 || err != nil {
			//fmt.Println("Read error:", err)
			break
		}

		str := string(input[0:n])
		var an []byte
		if str == "json" {

			fmt.Println("test json")
			res, er := testProtocols.Test_json()
			if er != nil {
				an, _ = json.Marshal(dto.Ans{dto.Answer{}, "", er})

			} else {

				an, _ = json.Marshal(dto.Ans{res, "json", nil})
			}
			//conn.Write(an)

		} else if str == "xml" {

			fmt.Println("test xml")
			res, er := testProtocols.Test_xml()
			if er != nil {
				an, _ = json.Marshal(dto.Ans{dto.Answer{}, "", er})

			} else {

				an, _ = json.Marshal(dto.Ans{res, "xml", nil})
			}
			//conn.Write(an)

		} else if str == "msgpack" {
			fmt.Println("test msgpack")
			res, er := testProtocols.Test_msgpack()
			if er != nil {
				an, _ = json.Marshal(dto.Ans{dto.Answer{}, "", er})

			} else {

				an, _ = json.Marshal(dto.Ans{res, "msgpack", nil})
			}

		} else {
			//fmt.Println("AAA")
			an, _ = json.Marshal(dto.Ans{dto.Answer{}, "", errors.New("fail request")})

		}
		// отправляем данные клиенту

		conn.Write(an)
	}
}
