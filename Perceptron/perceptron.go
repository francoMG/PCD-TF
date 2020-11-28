
package perceptron

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"strconv"
  "time"
	"net/http"
)

type Perceptron struct {
	Eta     float64
	Pesos []float64
	Numiter int
}

func shuffleDf(matrix [][]string)([][]string){
			rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(matrix), func(i, j int) {
			 matrix[:][i],matrix[:][j] =  matrix[:][j],matrix[:][i]
			  })
		return matrix
}
func LeerCSVdeURL(url string) ([][]string, error) {
		
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	dataframe := [][]string{}
	
	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	reader.Comma = ','
	reader.LazyQuotes = true

	
		for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		dataframe = append(dataframe, record)
	}

	return dataframe, nil
}

func HoldoutSets(df [][]string, alpha float64) ([][]float64, []float64, [][]float64, []float64)  {
 //hold out sets 80-20

 	X := [][]float64{}
	Y := []float64{}
	cont := 0
	for _, row := range df {

		if(cont != 0 ){
			
			//convert str slice to float slice
		temp := []float64{}
		for _, i := range row[1:] {
			parsedValue, err := strconv.ParseFloat(i, 64)
			if err != nil {
				panic(err)
			}
			temp = append(temp, parsedValue)
		}

		//agregando el row a X
		X = append(X, temp)

		//agregando su respectiva clase en Y
		if row[0] == "0"{
			Y = append(Y, 0.0)
		} else {
			Y = append(Y, 1.0)
		}

		
		}
		cont += 1
	}

	split:= int(float64(len(df))*alpha)
	X_train := X[:][:split]
	Y_train := Y[:][:split]
	X_test := X[:][split:]
	Y_test := Y[:][split:]
	
	return X_train, Y_train,  X_test, Y_test
}

func activateFunction(combiLin float64) float64 {
	if combiLin > 0 {
		return 1.0
	} else {
		return 0.0
	}
}



func (p *Perceptron)Predict(x []float64) float64 {
	var combinacionLinear float64
	
	for i := 0; i < len(x); i++ {
		combinacionLinear += x[i] + p.Pesos[i+1]
	}
	combinacionLinear += p.Pesos[0]
	return activateFunction(combinacionLinear)
}
func (p *Perceptron)PredictOnly(x []float64) float64 {
	var combinacionLinear float64
	fmt.Println("ONLY",len(p.Pesos))
	fmt.Println("I CANT HEAR U ",len(x))
	for i := 0; i < len(x); i++ {
		combinacionLinear += x[i] + p.Pesos[i+1]
	}
	combinacionLinear += p.Pesos[0]
	return activateFunction(combinacionLinear)
}

func (p *Perceptron) train(X [][]float64, Y []float64, c <-chan bool, next chan<- bool,index int) {

 for range c {
		error := 0
		for i := 0; i < len(X); i++ {
			y_pred := p.Predict(X[i])
			update := p.Eta * (Y[i] - y_pred)
			p.Pesos[0] += update
			for j := 0; j < len(X[i]); j++ {
				p.Pesos[j+1] += update * X[i][j]

				}
				if update != 0 {
					error += 1
				}
			}

			fmt.Println("Porcentaje de error: ",float64(error) / float64(len(Y)))
			fmt.Println("Numero de iteracion paralela!: ", index)
	 		next<- true
 }

 close(next)

}

func (p *Perceptron) Fit(X [][]float64, Y []float64, c chan bool) {
	//Inicializando Pesos forma aleatoria
	p.Pesos = []float64{}
	for i := 0; i <= len(X[0]); i++ {
		if i == 0 {
			p.Pesos = append(p.Pesos, 1.0)
		} else {
			p.Pesos = append(p.Pesos, rand.NormFloat64())
		}
	}

	//Ajustar Pesos en manera concurrente por Pipeline

		ch := make([]chan bool, p.Numiter + 1 )
  	ch[0] = make(chan bool)

 	 for iter := 0; iter < p.Numiter; iter++ {
				ch[iter+1] = make(chan bool)
        go p.train(X,Y,ch[iter], ch[iter+1],iter)
	
		}
  
 		go func(){
        
        ch[0]<- true
        
        close(ch[0])
    }()

		for range ch[p.Numiter] {
    }

	c <- true
}


// func main() {
// 	url := "https://raw.githubusercontent.com/Marcelpv96/Heart-Disease-UCI/master/data/heart.csv"
// 	//df, err := leerCSVdeURL(url)
// 	if err != nil {
// 		panic(err)
// 	}

// //omitiendo la primer fila que son los nombres de las columnas
// 	df = df[:][1:]
// 	//haciendo un shuffle para que el dataset este mezclado
// 	df = shuffleDf(df)



// 	//separando la data en X y Y
// 	X := [][]float64{}
// 	Y := []float64{}
// 	for _, row := range df {
		
// 		//convert str slice to float slice
// 		temp := []float64{}
// 		for _, i := range row[:len(row)-1] {
// 			parsedValue, err := strconv.ParseFloat(i, 64)
// 			if err != nil {
// 				panic(err)
// 			}
// 			temp = append(temp, parsedValue)
// 		}

// 		//agregando el row a X
// 		X = append(X, temp)

// 		//agregando su respectiva clase en Y
// 		if row[len(row)-1] == "0"{
// 			Y = append(Y, 0.0)
// 		} else {
// 			Y = append(Y, 1.0)
// 		}

// 	}



// 	channel := make(chan bool)

// 	 //hold out sets 80-20
// 	split:= int(float64(len(X))*0.8)
// 	X_train := X[:][:split]
// 	Y_train	:= Y[:][:split]
// 	X_test := X[:][split:]
// 	Y_test := Y[:][split:]

// 	//training
// 	perceptron := Perceptron{0.0001, []float64{}, 100}
// 	go perceptron.fit(X_train, Y_train,channel)
// 	<- channel
//  //testing
// 	perceptron.test(X_test,Y_test)
  
// }

func (p*Perceptron) scores(cMatrix [][] int ){
	TP := float64(cMatrix[1][1])
	FN := float64(cMatrix[1][0])
	TN := float64(cMatrix[0][0])
	FP :=float64(cMatrix[0][1])
	
	accuracy := float64((TP+TN)/(TP+TN+FP+FN))
 fmt.Println("---------------------------------")
 fmt.Println("Accuracy: ",accuracy)
 recall := TP/(TP+FN)
 fmt.Println("Recall: ", recall)
 precision := TP/(TP+FP) 
 fmt.Println("Precision: ", precision)
 fmt.Println("F1 Score: ", (2*precision*recall)/(precision+recall))
}
func (p*Perceptron) test(X_test[] [] float64,Y_test[] float64){
		confusionMatrix := [][]int{{0,0},{0,0}}
	
	for i := 0; i < len(X_test); i++ {
		realClass := int(Y_test[i])
		prediction := int(p.Predict(X_test[i]))
		//fmt.Println(realClass, prediction)
		confusionMatrix[realClass][prediction] = confusionMatrix[realClass][prediction] +1
	}
	p.scores(confusionMatrix)
}