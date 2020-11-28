package endpoints

import (
	perceptron "../Perceptron"
	"github.com/gin-gonic/gin"


	"../entities"
	"net/http"
	"../consenso"
)




func PerceptronPredict(context *gin.Context) {
	
	var pacient entities.Pacient
	channel := make(chan bool)
	if context.Bind(&pacient) == nil {

		url := "https://raw.githubusercontent.com/arian100y/retoBCP-UserSubscriptions-API/main/diagnostic%20(2).csv"
		df, _ := perceptron.LeerCSVdeURL(url)

		feature := []float64{
			
			pacient.RadiusMea,
			pacient.TextureMean  , 
			pacient.PerimeterMean ,
			pacient.AreaMean  ,
			pacient.SmoothnessMea ,
			pacient.CompactnessMean  ,
			pacient.ConcavityMean ,
			pacient.ConcavePointsMea,
			pacient.SymmetryMean ,
			pacient.FractalDimensionMean ,
			pacient.RadiusSe  ,
			pacient.TextureSe ,
			pacient.PerimeterSe  ,
			pacient.AreaS ,
			pacient.SmoothnessSe  ,
			pacient.CompactnessSe ,
			pacient.ConcavitySe   , 
			pacient.ConcavePointsSe , 
			pacient.SymmetryS ,
			pacient.FractalDimensionSe ,  
			pacient.RadiusWorst  , 
			pacient.TextureWorst  ,
			pacient.PerimeterWors ,
			pacient.AreaWorst,
			pacient.SmoothnessWorst  ,
			pacient.CompactnessWorst  ,
			pacient.ConcavityWors ,
			pacient.ConcavePointsWorst  ,
			pacient.SymmetryWorst ,
			pacient.FractalDimensionWors ,
		}

		X_train, y_train, _, _ := perceptron.HoldoutSets(df, 0.8)

	
		//forest := perceptron.BuildForest(X_train, y_train, 10, 500, len(X_train[0]))
		percep := perceptron.Perceptron{0.0001, []float64{}, 10}

		go percep.Fit(X_train, y_train,channel)

		<- channel	
		//test_predicted := make([]string, 0)

		//X := make([]interface{}, 0)	
		X := make([]float64,0)
		
		for _, x := range feature {
			
			X = append(X, x)
		}  
		
		var result float64 = percep.Predict(X)
		result_int := int(result)

		//----------------------------------------------

		msg := consenso.Message{consenso.CNum, consenso.LocalAddress, result_int, pacient}

		for _, addr := range consenso.Addresses {
			
			consenso.Send(addr, msg)
		}

		for consenso.Prediction == -1 {

		}

		resultpred := consenso.Prediction
		consenso.Prediction = -1

		//---------------------------------------------

		message := ""

		if resultpred == 0 {
			message = "El paciente podria no padecer de cancer."
		} else {
			message = "El paciente podria padecer cancer."
		}

		context.JSON(http.StatusOK, gin.H{"result": resultpred, "message": message})
	}

}
