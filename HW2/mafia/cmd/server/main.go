package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mafia/internal/game"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

var maxPl int

var roles []int32

func remove(s []*Subscription, i int) []*Subscription {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

type Subscription struct {
	ch    chan *game.GameRoom
	mutex sync.Mutex
}

func (s *Subscription) send(mes *game.GameRoom) {
	s.mutex.Lock()
	s.ch <- mes
	s.mutex.Unlock()
}

type Topic struct {
	subs  []*Subscription
	mutex sync.Mutex
}

type GamesManager struct {
	cuWait int32
	cuRoom int32
	rooms  map[int32]*RoomManager
	mutex  sync.Mutex
	ch     *amqp.Channel
}

type Player struct {
	sub    *Subscription
	chatCh chan []byte
	name   string
	alive  bool
	role   int32
}

type RoomManager struct {
	memb   []Player
	top    *Topic
	mutex  sync.Mutex
	cuRoom int32
	ch     *amqp.Channel
	qu     *amqp.Queue
	//isDay        bool
	cuDay              int32
	cntVotesToFi       int32
	isKil              bool
	isCheck            bool
	cuKilPl            int32
	checks             []int32
	checksRoles        []int32
	votes              [10]int32
	cntVotes           int32
	numbLivPl          int32
	igGameFinish       bool
	numbPlToFinishChat int32
}

func (r *RoomManager) FinishGame() {
	fi := &game.FinishGame{
		GameRes: 0,
	}
	resp := &game.GameRoom{
		GameState: &game.GameRoom_Fi{
			Fi: fi,
		},
	}
	r.top.Publish(resp)

	r.igGameFinish = true
}

func (r *RoomManager) CheckGameFinish() bool {
	for _, pl := range r.memb {
		if pl.role == 1 && !pl.alive {
			fi := &game.FinishGame{
				GameRes: 2,
			}
			resp := &game.GameRoom{
				GameState: &game.GameRoom_Fi{
					Fi: fi,
				},
			}
			r.top.Publish(resp)

			r.igGameFinish = true
			return true
		}
	}

	if r.numbLivPl == 2 {
		fi := &game.FinishGame{
			GameRes: 1,
		}
		resp := &game.GameRoom{
			GameState: &game.GameRoom_Fi{
				Fi: fi,
			},
		}
		r.top.Publish(resp)

		r.igGameFinish = true
		return true
	}

	return false
}

func (s *Subscription) Messages() chan *game.GameRoom {
	return s.ch
}

func (t *Topic) Subscribe() (*Subscription, error) {
	t.mutex.Lock()
	t.subs = append(t.subs, &Subscription{ch: make(chan *game.GameRoom)})
	res := t.subs[len(t.subs)-1]
	t.mutex.Unlock()
	return res, nil
}

func (t *Topic) Publish(message *game.GameRoom) error {
	t.mutex.Lock()

	for i, _ := range t.subs {
		//fmt.Println("send mes")
		t.subs[i].ch <- message
	}
	t.mutex.Unlock()
	return nil
}

func (r *RoomManager) PublishToAlive(message *game.GameRoom) error {
	r.mutex.Lock()

	for _, pl := range r.memb {
		//fmt.Println("send mes")
		if pl.alive {
			pl.sub.ch <- message
		}
		//t.subs[i].ch <- message
	}
	r.mutex.Unlock()
	return nil
}

func (g *GamesManager) Add(name_ string) (*RoomManager, *Subscription, int32, int32) {
	g.mutex.Lock()
	g.cuWait += 1
	fmt.Println("current number of wait", g.cuWait, maxPl)

	if int(g.cuWait) == maxPl {
		//fmt.Println("game start")
		g.cuWait = 0
		go g.rooms[g.cuRoom].start_game()
	}

	if g.cuWait == 1 {
		g.cuRoom += 1
		g.rooms[g.cuRoom] = &RoomManager{cuRoom: g.cuRoom, top: &Topic{}}
		g.rooms[g.cuRoom].ch = g.ch
	}
	sub_, _ := g.rooms[g.cuRoom].top.Subscribe()
	cu := len(g.rooms[g.cuRoom].memb)
	cuR := g.cuRoom
	g.rooms[g.cuRoom].memb = append(g.rooms[g.cuRoom].memb, Player{sub: sub_, name: name_, chatCh: make(chan []byte)})
	//res := &g.rooms[g.cuRoom]
	g.mutex.Unlock()
	return g.rooms[g.cuRoom], sub_, int32(cu), int32(cuR)
}

func (g *GamesManager) Unsub(pl int32, cuG int32) {
	g.mutex.Lock()
	g.cuWait -= 1
	if int(g.cuWait) == maxPl {
		go g.rooms[cuG].FinishGame()
	} else {
		for i, _ := range g.rooms[g.cuRoom].memb {
			if int32(i) == pl {
				fmt.Println("del", i)
				for j, sub := range g.rooms[g.cuRoom].top.subs {
					if sub == g.rooms[cuG].memb[i].sub {
						g.rooms[g.cuRoom].top.subs = remove(g.rooms[g.cuRoom].top.subs, j)
						break
					}
				}

				g.rooms[cuG].memb[i] = g.rooms[cuG].memb[len(g.rooms[cuG].memb)-1]
				g.rooms[cuG].memb = g.rooms[cuG].memb[:len(g.rooms[cuG].memb)-1]

				break
			}
		}
	}
	//fmt.Println("current number of wait", g.cuWait)

	g.mutex.Unlock()
}

func (r *RoomManager) start_game() {

	//g := GameManager{}
	fmt.Println("game func start")
	q, err := r.ch.QueueDeclare(
		string(r.cuRoom), // name
		false,            // durable
		true,             // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		log.Fatalf("failed to declare a queue. Error: %s", err)
	}
	r.qu = &q

	go r.sendMessageToClient()
	//rand.Shuffle(roles)
	rand.Shuffle(len(roles),
		func(i, j int) { roles[i], roles[j] = roles[j], roles[i] })

	var names []string
	for _, pl := range r.memb {
		names = append(names, pl.name)
	}
	for i, _ := range r.memb {
		r.votes[i] = 0
		r.memb[i].role = roles[i]
		r.memb[i].alive = true
		r.numbLivPl += 1
		g_st := game.GameStart{
			RoomNumb: r.cuRoom,
			Role:     r.memb[i].role,
			GameMemb: names,
			PlNumb:   int32(i),
		}

		resp := &game.GameRoom{
			GameState: &game.GameRoom_GameSt{
				GameSt: &g_st,
			},
		}

		r.memb[i].sub.send(resp)
	}
	r.dayStart()
	//r.top.Publish(&resp)
}

func (r *RoomManager) dayStart() {
	fmt.Println("Day start")
	if r.igGameFinish {
		return
	}
	r.cntVotesToFi = 0
	r.isCheck = false
	r.isKil = false
	r.cntVotes = 0
	var aliveMemb []int32
	for i, val := range r.memb {
		if val.alive {
			aliveMemb = append(aliveMemb, int32(i+1))
		}
	}
	g_da := &game.DayState{
		CurDay:        r.cuDay,
		AliveGameMemb: aliveMemb,
	}

	resp := &game.GameRoom{
		GameState: &game.GameRoom_Day{
			Day: g_da,
		},
	}
	r.cuDay += 1
	r.top.Publish(resp)
	if r.cuDay != 1 {
		r.askComis()
	} else {
		r.startChat()
		//r.ascUsToFinish()
	}
}

func (r *RoomManager) startChat() {
	r.numbPlToFinishChat = 0
	g_st := &game.StartChat{}

	resp := &game.GameRoom{
		GameState: &game.GameRoom_Chat{
			Chat: g_st,
		},
	}
	r.PublishToAlive(resp)
}

func (r *RoomManager) reqToFinishChat() bool {
	r.mutex.Lock()
	if r.numbPlToFinishChat >= r.numbLivPl {
		r.mutex.Unlock()
		return false
	}
	r.numbPlToFinishChat += 1
	fmt.Println(r.numbPlToFinishChat, r.numbLivPl)
	if r.numbPlToFinishChat == r.numbLivPl {
		go r.FinishChat()
	}
	r.mutex.Unlock()
	return true

}

func (r *RoomManager) FinishChat() {
	fmt.Println("finishChat")
	if r.cuDay == 1 {
		r.askUsToFinish()
	} else {
		r.AskStartVote()
	}
}

func (r *RoomManager) askComis() {
	if r.igGameFinish {
		return
	}
	g_da := &game.ComisReq{}

	resp := &game.GameRoom{
		GameState: &game.GameRoom_ComReq{
			ComReq: g_da,
		},
	}

	for i, _ := range r.memb {
		if r.memb[i].role == 2 {
			if r.memb[i].alive {
				r.memb[i].sub.send(resp)
			} else {
				r.startChat()
			}
		}
	}

}

func (r *RoomManager) publishCheck(fl bool) {
	if r.igGameFinish {
		return
	}
	if fl {
		checks := &game.ComisChecks{
			Players:      r.checks,
			PlayersRoles: r.checksRoles,
		}

		resp := &game.GameRoom{
			GameState: &game.GameRoom_Checks{
				Checks: checks,
			},
		}
		r.top.Publish(resp)
	}
	r.startChat()
}

func (r *RoomManager) AskStartVote() {
	fmt.Println("ask to start vote")
	if r.igGameFinish {
		return
	}
	var aliveMemb []int32
	for i, val := range r.memb {
		if val.alive {
			aliveMemb = append(aliveMemb, int32(i+1))
		}
	}

	votes := &game.StartVote{
		AliveGameMemb: aliveMemb,
	}

	resp := &game.GameRoom{
		GameState: &game.GameRoom_Votes{
			Votes: votes,
		},
	}

	r.PublishToAlive(resp)
}

func (r *RoomManager) FinishDay() {
	if r.igGameFinish {
		return
	}
	var aliveMemb []int32
	for i, val := range r.memb {
		if val.alive {
			aliveMemb = append(aliveMemb, int32(i+1))
		}
	}

	ni := &game.NightState{
		AliveGameMemb: aliveMemb,
	}

	resp := &game.GameRoom{
		GameState: &game.GameRoom_Night{
			Night: ni,
		},
	}

	for _, pl := range r.memb {
		if pl.role == 2 && !pl.alive {
			r.isCheck = true
		}
	}
	r.top.Publish(resp)
}

func (r *RoomManager) reqToFinishDay() {
	if r.igGameFinish {
		return
	}
	r.mutex.Lock()
	r.cntVotesToFi += 1
	r.mutex.Unlock()

	if r.cntVotesToFi == int32(len(r.memb)) {
		r.FinishDay()
	}
}

func (r *RoomManager) killPl(pl int32, cuPl int32) bool {
	if r.igGameFinish {
		return true
	}
	if pl >= 0 && int(pl) < len(r.memb) && r.memb[pl].alive && r.memb[cuPl].role == 1 {
		fl := false
		r.mutex.Lock()
		r.memb[pl].alive = false
		r.numbLivPl -= 1
		r.cuKilPl = pl
		r.isKil = true
		if r.isCheck {
			fl = true
		}
		r.mutex.Unlock()
		if fl {
			go r.fiNight()
		}
		return true
	} else {
		return false
	}
}

func (r *RoomManager) checkPl(pl int32, cuPl int32) (bool, bool) {
	if r.igGameFinish {
		return true, true
	}
	if pl >= 0 && int(pl) < len(r.memb) && r.memb[cuPl].role == 2 {
		fl := false
		r.mutex.Lock()
		r.isCheck = true
		r.checks = append(r.checks, pl)
		r.checksRoles = append(r.checksRoles, r.memb[pl].role)
		if r.isKil {
			fl = true
			//r.fiNight()
		}
		r.mutex.Unlock()
		if fl {
			go r.fiNight()
		}
		return true, r.memb[pl].role == 1
	} else {
		return false, false
	}
}

func (r *RoomManager) fiNight() {
	if r.igGameFinish {
		return
	}
	r.mutex.Lock()
	fiNi := &game.FiNightState{
		KilPl: r.cuKilPl,
	}

	resp := &game.GameRoom{
		GameState: &game.GameRoom_FiNight{
			FiNight: fiNi,
		},
	}
	r.mutex.Unlock()
	r.top.Publish(resp)

	if !r.CheckGameFinish() {
		r.dayStart()
	}
}

func (r *RoomManager) askUsToFinish() {
	fmt.Println("Ask User to finish")
	resp := &game.GameRoom{
		GameState: &game.GameRoom_FiR{
			FiR: &game.FinishReq{},
		},
	}
	r.top.Publish(resp)
}

type messageFromPl struct {
	Mes    string
	PlNumb int32
}

func enc(st messageFromPl) []byte {
	b, _ := json.Marshal(st)
	return b
}

func dec(by []byte) messageFromPl {
	var st messageFromPl
	json.Unmarshal(by, &st)
	return st
}

func (r *RoomManager) sendMessage(mes string, plNumb int32) bool {
	r.mutex.Lock()
	fmt.Println("start publish")
	mesWithPl := messageFromPl{
		Mes:    mes,
		PlNumb: plNumb,
	}

	err := r.ch.PublishWithContext(context.Background(),
		"",        // exchange
		r.qu.Name, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        enc(mesWithPl),
		})
	if err != nil {
		fmt.Println("Get error")
		fmt.Println(err)
	}
	//r.qu
	r.mutex.Unlock()
	return true

}

func (r *RoomManager) sendMessageToClient() {
	messages, err := r.ch.Consume(
		r.qu.Name, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Fatalf("failed to register a consumer. Error: %s", err)
	}

	for message := range messages {
		for _, mem := range r.memb {
			mem.chatCh <- message.Body
		}
		//log.Printf("received a message: %s", message.Body)
	}

}

func (r *RoomManager) vote(pl int32) bool {
	if r.igGameFinish {
		return true
	}
	if pl >= 0 && int(pl) < len(r.memb) && r.memb[pl].alive {
		r.mutex.Lock()
		r.cntVotes += 1
		r.votes[pl] += 1
		fmt.Println(r.cntVotes, r.numbLivPl)
		if int(r.cntVotes) == int(r.numbLivPl) {
			fmt.Println("vote end")
			ma := 0
			pl_ := 0
			fl := 0
			for i, _ := range r.memb {
				if int(r.votes[i]) > ma {
					ma = int(r.votes[i])
					pl_ = i
					fl = 1
				} else if int(r.votes[i]) == ma {
					fl = 0
				}
			}
			res := &game.ResultOfVotes{
				IsCh: false,
				Pl:   0,
			}
			if fl == 1 {
				r.memb[pl_].alive = false
				r.numbLivPl -= 1
				res = &game.ResultOfVotes{
					IsCh: true,
					Pl:   int32(pl_),
				}
			}
			resp := &game.GameRoom{
				GameState: &game.GameRoom_Res{
					Res: res,
				},
			}
			r.top.Publish(resp)
			if !r.CheckGameFinish() {
				r.askUsToFinish()
			}
		}
		r.mutex.Unlock()
		return true
	} else {
		return false
	}
}

func main() {
	fmt.Println(os.Args)

	serverAddr := os.Getenv("SERVER_HOST")
	port := os.Getenv("SERVER_PORT")
	listenAddr := fmt.Sprintf("amqp://guest:guest@%s:%s/", serverAddr, port)

	fmt.Println(listenAddr)

	conn, err := amqp.Dial(listenAddr) // Создаем подключение к RabbitMQ
	for err != nil {
		conn, err = amqp.Dial(listenAddr)
		//log.Fatalf("unable to open connect to RabbitMQ server. Error: %s", err)
	}

	defer func() {
		_ = conn.Close() // Закрываем подключение в случае удачной попытки
	}()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open channel. Error: %s", err)
	}

	defer func() {
		_ = ch.Close() // Закрываем канал в случае удачной попытки открытия
	}()

	maxPl, _ = strconv.Atoi(os.Args[1])
	rand.Seed(time.Now().UnixNano())

	roles = []int32{1, 2}
	for i := 0; i < maxPl-2; i++ {
		roles = append(roles, 0)
	}
	//curRomNumb = 1
	// create listiner

	lis, err := net.Listen("tcp", ":50005")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create grpc server
	s := grpc.NewServer()
	game.RegisterConnectToGameServer(s, &server{manag: GamesManager{rooms: make(map[int32]*RoomManager), ch: ch}})

	log.Println("start server")
	// and start...
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// type RoomHelper struct {
// 	streams []game.ConnectToGame_JoinGameRoomServer
// 	mutex   sync.Mutex
// }

type server struct {
	game.UnimplementedConnectToGameServer
	manag GamesManager
}

func (s *server) JoinGameRoom(req *game.JoinGameRoomRequest, stream game.ConnectToGame_JoinGameRoomServer) error {

	//fmt.Println("new user come")

	man, sub, nu, cuR := s.manag.Add(req.Name)
	fmt.Println("Add", nu, "player")
	//sub, _ := top.Subscribe()
	mes := sub.Messages()
	//fmt.Println("wait for room")

	select {
	case room := <-mes:
		//fmt.Println("res mes")
		err := stream.Send(room)
		if err != nil {
			return err
		}
	case <-stream.Context().Done():
		s.manag.Unsub(nu, cuR)
		return nil
	}

	for {
		select {
		case room := <-mes:
			//fmt.Println("res mes")
			err := stream.Send(room)
			if err != nil {
				return err
			}
		case <-stream.Context().Done():
			man.FinishGame()
			return nil
		}
	}
	//return nil
}

func (s *server) GetChatStream(req *game.ChatStreamRequest, stream game.ConnectToGame_GetChatStreamServer) error {

	for {
		select {
		case str := <-s.manag.rooms[req.RoomNumb].memb[req.PlNumb].chatCh:
			res := dec(str)
			resp := game.ChatStream{
				Mes:    res.Mes,
				PlNumb: res.PlNumb,
			}
			stream.Send(&resp)
		case <-stream.Context().Done():
			return nil
			//str := <-s.manag[req.RoomNumb].memb[req.PlNumb].
		}
	}
	return nil
}

func (s *server) FinishDay(c context.Context, request *game.FinishDayRequest) (*game.FinishDayResp, error) {
	s.manag.rooms[request.RoomNumb].reqToFinishDay()
	resp := game.FinishDayResp{
		Suc: true,
	}
	return &resp, nil
}

func (s *server) KillPlayer(c context.Context, request *game.KillPlayerRequest) (*game.KillPlayerResp, error) {
	suc := s.manag.rooms[request.RoomNumb].killPl(request.Goal, request.PlNumb)
	resp := game.KillPlayerResp{
		Suc: suc,
	}
	return &resp, nil
}

func (s *server) CheckPlayer(c context.Context, request *game.CheckPlayerRequest) (*game.CheckPlayerResp, error) {
	suc, res := s.manag.rooms[request.RoomNumb].checkPl(request.Goal, request.PlNumb)
	resp := game.CheckPlayerResp{
		Suc:   suc,
		IsMaf: res,
	}
	return &resp, nil
}

func (s *server) PublishChecks(c context.Context, request *game.Checks) (*game.PubResp, error) {
	s.manag.rooms[request.RoomNumb].publishCheck(request.PubCheck)
	resp := game.PubResp{
		Suc: request.PubCheck,
	}
	return &resp, nil
}

func (s *server) Vote(c context.Context, request *game.VoteRequest) (*game.VoteResp, error) {
	suc := s.manag.rooms[request.RoomNumb].vote(request.Goal)
	resp := game.VoteResp{
		Suc: suc,
	}
	return &resp, nil
}

func (s *server) FinishChat(c context.Context, request *game.FinishChatRequest) (*game.FinishChatResp, error) {
	suc := s.manag.rooms[request.RoomNumb].reqToFinishChat()
	resp := game.FinishChatResp{
		Suc: suc,
	}
	return &resp, nil
}

func (s *server) SendMessage(c context.Context, request *game.ChatMes) (*game.ChatMesResp, error) {
	fmt.Println("want to send Mes")
	suc := s.manag.rooms[request.RoomNumb].sendMessage("Player "+strconv.Itoa(int(request.PlNumb))+" said: "+request.Mes, request.PlNumb)
	resp := game.ChatMesResp{
		Suc: suc,
	}
	return &resp, nil
}
