package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"log"
	"mafia/internal/game"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func recieveChatMes(roomNumb int32, plNumb int32, client game.ConnectToGameClient) {
	colorReset := "\033[0m"
	colors := []string{"\033[32m", "\033[33m", "\033[34m", "\033[35m", "\033[36m", "\033[37m"}
	req := game.ChatStreamRequest{
		RoomNumb: int32(roomNumb),
		PlNumb:   int32(plNumb),
	}
	streamChat, _ := client.GetChatStream(context.Background(), &req)
	for {
		recv, _ := streamChat.Recv()
		//fmt.Println(recv.Mes)
		fmt.Println(string(colors[int(recv.PlNumb)%len(colors)]), recv.Mes, string(colorReset))
	}
}

func main() {
	colorReset := "\033[0m"

	colorRed := "\033[31m"
	colorGreen := "\033[32m"
	colorYellow := "\033[33m"
	// colorBlue := "\033[34m"
	colorPurple := "\033[35m"
	colorCyan := "\033[36m"
	// colorWhite := "\033[37m"

	names := []string{"Mary", "Bob", "Tim", "Lola", "Kate", "Sara",
		"John", "Bill", "Lily", "Dima", "Rick", "Morty", "Paul", "Derek",
		"Ira", "Leo", "Vini", "Ney", "Cris", "Olya", "Darya", "Mark", "Mayk", "Peter"}
	auto := os.Args[1] == "bot"

	rand.Seed(time.Now().UnixNano())
	conn, err := grpc.Dial(":50005", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}

	// create stream
	roles := make(map[int32]string)
	roles[0] = "civilian"
	roles[1] = "mafia"
	roles[2] = "сommissar"
	client := game.NewConnectToGameClient(conn)

	var name_ string

	name := names[rand.Int()%len(names)]
	fmt.Println("Print your name or press enter to use default")
	if !auto {
		fmt.Scanf("%s\n", &name_)
	}
	if name_ != "" {
		name = name_
	} else {
		fmt.Println("Your name is", name)
	}

	request := &game.JoinGameRoomRequest{
		Name: name,
	}
	stream, err := client.JoinGameRoom(context.Background(), request)

	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}

	roomNumb := 0
	role := 0
	plNumb := 0
	cntPl := 0
	alive := true
	fmt.Println("Please wait for another players")
	for {
		resp, err := stream.Recv()

		if err == io.EOF {
			fmt.Println("recieve i0.EOF")
			//done <- true //means stream is finished
			return
		}
		if resp.GetGameSt() != nil {
			roomNumb = int(resp.GetGameSt().RoomNumb) //RoomNumb
			role = int(resp.GetGameSt().Role)
			plNumb = int(resp.GetGameSt().GetPlNumb())
			cntPl = int(len(resp.GetGameSt().GetGameMemb()))
			fmt.Println(string(colorPurple), "Game Start, your room numb is", roomNumb, string(colorReset))
			for i, memb := range resp.GetGameSt().GameMemb {
				fmt.Println("The name of", i+1, "player is", memb)
			}
			fmt.Println("Your number is", resp.GetGameSt().PlNumb+1)
			fmt.Println(string(colorYellow), "Your role is", roles[resp.GetGameSt().Role], string(colorReset))
			go recieveChatMes(int32(roomNumb), int32(plNumb), client)

		}
		if resp.GetDay() != nil {
			fmt.Println(string(colorYellow), "Day starts", string(colorReset))
			if resp.GetDay().CurDay == 0 {
				fmt.Println("This is the first day, so you can't kill anyone")
			}
		}

		if resp.GetFiR() != nil {
			fmt.Println("If you want to finish this day press enter")
			if !auto {
				var enter_ string
				fmt.Scanf("%s\n", &enter_)
			}
			fmt.Println("Please, wait for other players to finish the day")
			req := &game.FinishDayRequest{
				RoomNumb: int32(roomNumb),
			}
			client.FinishDay(context.Background(), req)
		}

		if resp.GetComReq() != nil {
			fmt.Println("You are comissar, do you want to publish your checks?")

			var ans string
			if !auto {
				fmt.Scanf("%s\n", &ans)
			} else {
				ans = "yes"
			}
			fl := false
			if ans == "Yes" || ans == "yes" || ans == "Yes\n" || ans == "yes\n" {
				fl = true
			}
			req := game.Checks{
				RoomNumb: int32(roomNumb),
				PubCheck: fl,
			}
			client.PublishChecks(context.Background(), &req)
		}

		if resp.GetChecks() != nil {
			fmt.Println(string(colorPurple), "Comissar published his checks", string(colorReset))
			for i, v := range resp.GetChecks().Players {
				fmt.Println("The role of", v+1, "player is", roles[resp.GetChecks().PlayersRoles[i]])
			}
		}

		if resp.GetVotes() != nil {
			fmt.Println("Please, select the player you want to vote for")
			fmt.Println("Choose one of these players", resp.GetVotes().AliveGameMemb)

			var numb int32

			if !auto {
				fmt.Scanf("%d\n", &numb)
			} else {
				le := len(resp.GetVotes().AliveGameMemb)
				numb = resp.GetVotes().AliveGameMemb[rand.Int()%le]
			}
			req := game.VoteRequest{
				Goal:     numb - 1,
				RoomNumb: int32(roomNumb),
			}
			suc, _ := client.Vote(context.Background(), &req)
			for !suc.Suc {
				fmt.Println(string(colorRed), "You cannot vote for this player", string(colorReset))
				fmt.Println("Choose one of these players", resp.GetVotes().AliveGameMemb)
				var numb int32
				fmt.Scanf("%d\n", &numb)
				req := game.VoteRequest{
					Goal:     numb - 1,
					RoomNumb: int32(roomNumb),
				}
				suc, _ = client.Vote(context.Background(), &req)
			}
			fmt.Println("Your vote is taken into account")
		}

		if resp.GetRes() != nil {
			if resp.GetRes().IsCh {
				if resp.GetRes().Pl == int32(plNumb) {
					alive = false
					fmt.Println(string(colorRed), "You are killed, watch the game", string(colorReset))
				} else {
					fmt.Println("The city decided to kill player number", resp.GetRes().Pl+1)
				}
			} else {
				fmt.Println("The city has not decided who to kill")
			}
		}

		if resp.GetNight() != nil {
			fmt.Println(string(colorCyan), "Night starts", string(colorReset))
			if role == 1 {
				fmt.Println("Your role is mafia, so choose the number of the player you want to kill")
				fmt.Println("Live player numbers are", resp.GetNight().AliveGameMemb)
				var numb int32
				if !auto {
					fmt.Scanf("%d\n", &numb)
				} else {
					le := len(resp.GetNight().AliveGameMemb)
					numb = resp.GetNight().AliveGameMemb[rand.Int()%le]
				}
				req := game.KillPlayerRequest{
					Goal:     numb - 1,
					PlNumb:   int32(plNumb),
					RoomNumb: int32(roomNumb),
				}
				suc, _ := client.KillPlayer(context.Background(), &req)
				for !suc.Suc {
					fmt.Println(string(colorRed), "You can't kill this player, please choose someone from the live players", string(colorReset))
					fmt.Println("Live player numbers are", resp.GetNight().AliveGameMemb)

					fmt.Scanf("%d\n", &numb)
					req.Goal = numb - 1
					suc, _ = client.KillPlayer(context.Background(), &req)
				}
				fmt.Println("You kill player number", numb)
			} else if role == 2 && alive {
				fmt.Println("Your role is commissar, so choose the number of the player you want to check from 1 to", cntPl)
				var numb int32
				if !auto {
					fmt.Scanf("%d\n", &numb)
				} else {
					numb = int32(rand.Int()%4) + 1
				}
				req := game.CheckPlayerRequest{
					Goal:     numb - 1,
					PlNumb:   int32(plNumb),
					RoomNumb: int32(roomNumb),
				}
				suc, _ := client.CheckPlayer(context.Background(), &req)
				for !suc.Suc {
					fmt.Println(string(colorRed), "You can't check this player, please choose the number from 1 to", cntPl, string(colorReset))

					fmt.Scanf("%d\n", &numb)
					req.Goal = numb - 1
					suc, _ = client.CheckPlayer(context.Background(), &req)
				}
				if suc.IsMaf {
					fmt.Println("This player is mafia")
				} else {
					fmt.Println("This player is not mafia")
				}

			} else {
				fmt.Println("Please wait")
			}
		}

		if resp.GetFiNight() != nil {
			if resp.GetFiNight().KilPl == int32(plNumb) {
				alive = false
				fmt.Println(string(colorRed), "You are killed, watch the game", string(colorReset))
			} else {
				fmt.Println("The night is over, player number was killed that night is", resp.GetFiNight().KilPl+1)
			}
		}

		if resp.GetChat() != nil {
			fmt.Println("Now you can discuss the game with other players")
			fmt.Println("Print", string(colorCyan), "stop", string(colorReset), "when you want to end the chat")
			if !auto {
				for {
					//var name_ string
					in := bufio.NewReader(os.Stdin)

					name_, _ := in.ReadString('\n')
					mes := game.ChatMes{
						RoomNumb: int32(roomNumb),
						PlNumb:   int32(plNumb),
						Mes:      name_[:len(name_)-1],
					}
					if name_ == "stop" || name_ == "stop\n" {
						break
					}
					client.SendMessage(context.Background(), &mes)
				}
			}
			fmt.Println("Please wait other players")
			req := game.FinishChatRequest{
				RoomNumb: int32(roomNumb),
			}
			client.FinishChat(context.Background(), &req)

		}

		if resp.GetFi() != nil {
			fmt.Println(string(colorYellow), "Game finished", string(colorReset))
			if resp.GetFi().GameRes == 1 {
				fmt.Println("Mafia win")
				if role == 1 {
					fmt.Println(string(colorGreen), "You win", string(colorReset))
				} else {
					fmt.Println(string(colorRed), "You lose", string(colorReset))
				}
			} else if resp.GetFi().GameRes == 2 {
				fmt.Println("City win")
				if role == 1 {
					fmt.Println(string(colorRed), "You lose", string(colorReset))
				} else {
					fmt.Println(string(colorGreen), "You win", string(colorReset))
				}
			} else {
				fmt.Println("Someone left")
			}

			break
		}
	}

	//fmt.Println(response.Message)
}
