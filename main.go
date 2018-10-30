package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"dev.hocngay.com/hocngay/compile-test/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	handler.Init()

	r.GET("/test", handleTest)

	//Run service
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	err := r.Run(":" + port)
	if err != nil {
		panic(err)
	}
}

func handleTest(ctx *gin.Context) {
	ch := make(chan string, 200)
	startTime := time.Now()
	for i := 0; i < 200; i++ {
		go handler.CreateContainer(strconv.Itoa(i), "go", ch)
	}
	<-ch

	totalTime := time.Since(startTime)
	fmt.Println("Time execute:", totalTime)
}
