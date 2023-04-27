package main

import (
	"HW1/cmd/dto"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

func main() {

	host := os.Getenv("SERVER_HOST")
	port := os.Getenv("SERVER_PORT")
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", host, port))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	file, err := os.Open("inp.txt")
	if err != nil {
		fmt.Println("error ", err)
	}
	fileScanner := bufio.NewScanner(file)

	for fileScanner.Scan() {
		//var source string
		//fmt.Print("Введите слово: ")

		//handle errors while opening
		// if err != nil {
		// 	log.Fatalf("Error when opening file: %s", err)
		// }

		// read line by line

		source := fileScanner.Text()

		// отправляем сообщение серверу
		if n, err := conn.Write([]byte(source)); n == 0 || err != nil {
			fmt.Println(err)
			return
		}
		// получем ответ
		//fmt.Print("Перевод:")
		buff := make([]byte, 1024)
		n, _ := conn.Read(buff)
		//fmt.Println(buff[0:n])

		var an dto.Ans
		json.Unmarshal(buff[0:n], &an)
		if an.Err != nil {
			fmt.Println("error ", an.Err)
		} else {
			fmt.Println(an.Name, " - ", an.Ans.Mem, " - ", an.Ans.TimeOfSer, " - ", an.Ans.TimeOfDes)
		}

		// fmt.Print(string(buff[0:n]))
		// fmt.Println()
	}
	file.Close()
}
