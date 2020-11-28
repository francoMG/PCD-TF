package consenso

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"../entities"
)

var Prediction int = -1

const (
	CNum = iota
	maligno  = 1
	noMaligno  = 0
)

type Message struct {
	Code    int
	Address    string
	Op      int
	Pacient entities.Pacient
}

var LocalAddress string = "localhost:8100"

var Addresses = []string{
	"localhost:8200",
	"localhost:8300",
}

var chInfo chan map[string]int

func RunServer() {

	chInfo = make(chan map[string]int)

	go func() { chInfo <- map[string]int{} }()

	go server()

}

func server() {
	if ln, err := net.Listen("tcp", LocalAddress); err != nil {
		log.Panicln("Can't start listener on", LocalAddress)
	} else {
		defer ln.Close()
		fmt.Println("Listening on", LocalAddress)
		for {
			if conn, err := ln.Accept(); err != nil {
				log.Println("Can't accept", conn.RemoteAddr())
			} else {
				go handle(conn)
			}
		}
	}
}
func handle(conn net.Conn) {
	defer conn.Close()
	dec := json.NewDecoder(conn)
	var msg Message
	if err := dec.Decode(&msg); err != nil {
		log.Println("Can't decode from", conn.RemoteAddr())
	} else {
		fmt.Println(msg)
		switch msg.Code {
		case CNum:
			concensus(conn, msg)
		}
	}
}

func concensus(conn net.Conn, msg Message) {
	info := <-chInfo
	info[msg.Address] = msg.Op
	fmt.Println(info)
	if len(info) == len(Addresses) {
		ca, cb := 0, 0
		for _, op := range info {
			if op == maligno {
				ca++
			} else {
				cb++
			}
		}
		if ca > cb {
			Prediction = 1
		} else {
			Prediction = 0
		}
		info = map[string]int{}
	}
	go func() { chInfo <- info }()
}

func Send(remoteAddr string, msg Message) {
	
	if conn, err := net.Dial("tcp", remoteAddr); err != nil {
		log.Println("Can't dial", remoteAddr, err)
	} else {
		defer conn.Close()
		fmt.Println("Sending to", remoteAddr)
		enc := json.NewEncoder(conn)
		enc.Encode(msg)
	}
}
