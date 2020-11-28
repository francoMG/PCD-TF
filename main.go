package main

import (

	"./consenso"

	"./endpoints"
	"github.com/gin-gonic/gin"
)

func policyAPIcors() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		context.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		context.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		context.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		context.Next()
	}
}


func main() {

	go consenso.RunServer()

	router := gin.Default()

	router.Use(policyAPIcors())


	router.POST("/perceptron/predict",endpoints.PerceptronPredict)

	router.Run()

}
