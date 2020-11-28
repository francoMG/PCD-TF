package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"

	perceptron "./Perceptron"
	"./entities"

)



const (
	Cnum = iota
	maligno  = 1
	noMaligno  = 0
)

type Message struct {
	Code    int
	Addr    string
	Op      int
	Pacient entities.Pacient
}

var Pacient entities.Pacient = entities.Pacient{
	
	0.0,
	0.0,
	0.0,
	0.0,
	0.0,
	0.0,
	0.0,
	0.0,
	0.0,
	0.0,
	0.0,
	0.0,
	0.0,
	0.0,
	0.0,
	0.0,
	0.0,
	0.0,
	0.0,
	0.0,
	0.0,
	0.0,
	0.0,
	0.0,
	0.0,
	0.0,
	0.0,
	0.0,
	0.0,
	0.0,
	0.0,

}

var Msg Message
var Prediction int = -1
var predict bool = false


//para correr otro nodo con otro ip, corre este programa luego cambia al address y corre en otro terminal.
var LocalAddress string = "localhost:8300"

//var mainNodeAddress = "localhost:8100"

var chInfo chan map[string]int
var Addresses = []string{
	"localhost:8100",

}

func server() {
	if ln, err := net.Listen("tcp", LocalAddress); err != nil {
		log.Panicln("Error: can't start listener on", LocalAddress)
	} else {
		defer ln.Close()
		fmt.Println("Listening on: ", LocalAddress)
		for {
			if conn, err := ln.Accept(); err != nil {
				log.Println("Can't accept connection.", conn.RemoteAddr())
			} else {
				go handle(conn)
			}
		}
	}
}
func handle(conn net.Conn) {
	defer conn.Close()
	dec := json.NewDecoder(conn)

	if err := dec.Decode(&Msg); err != nil {
		log.Println("Can't decode from", conn.RemoteAddr())
	} else {
		fmt.Println(Msg)
		predict = true
		
	}
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



func main() {

	chInfo = make(chan map[string]int)

	go func() { chInfo <- map[string]int{} }()

	go server()
	channel := make(chan bool)
	for {

		if predict {

			fmt.Println()
			url := "https://raw.githubusercontent.com/arian100y/retoBCP-UserSubscriptions-API/main/diagnostic%20(2).csv"
			df, _ := perceptron.LeerCSVdeURL(url)

			feature := []float64{
			
				Pacient.RadiusMea,
				Pacient.TextureMean  , 
				Pacient.PerimeterMean ,
				Pacient.AreaMean  ,
				Pacient.SmoothnessMea ,
				Pacient.CompactnessMean  ,
				Pacient.ConcavityMean ,
				Pacient.ConcavePointsMea,
				Pacient.SymmetryMean ,
				Pacient.FractalDimensionMean ,
				Pacient.RadiusSe  ,
				Pacient.TextureSe ,
				Pacient.PerimeterSe  ,
				Pacient.AreaS ,
				Pacient.SmoothnessSe  ,
				Pacient.CompactnessSe ,
				Pacient.ConcavitySe   , 
				Pacient.ConcavePointsSe , 
				Pacient.SymmetryS ,
				Pacient.FractalDimensionSe ,  
				Pacient.RadiusWorst  , 
				Pacient.TextureWorst  ,
				Pacient.PerimeterWors ,
				Pacient.AreaWorst,
				Pacient.SmoothnessWorst  ,
				Pacient.CompactnessWorst  ,
				Pacient.ConcavityWors ,
				Pacient.ConcavePointsWorst  ,
				Pacient.SymmetryWorst ,
				Pacient.FractalDimensionWors ,
			}
			
			X_train, y_train, _, _ := perceptron.HoldoutSets(df, 0.8)

		
			percep := perceptron.Perceptron{00001 + rand.Float64() * (1.0 - 0.00001), []float64{}, 100}
			go percep.Fit(X_train, y_train,channel)
			<- channel	
			//test_predicted := make([]string, 0)

			X := make([]float64,0 )
			for _, x := range feature {
				X = append(X, x)
			}

			
			var result float64 = percep.Predict(X)
		result_int := int(result)

			fmt.Println("Resultado de prediccion: ", result_int)

			msg := Message{Cnum, LocalAddress, result_int, Pacient}

			
			for _, addr := range Addresses {
				Send(addr, msg)
			}

			predict = false

		}
	}

}
