package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"dev.hocngay.com/hocngay/compile-test/handler"
	"dev.hocngay.com/hocngay/compile-test/model"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

var m sync.Mutex
var queue = model.Queue{}

func main() {
	r := gin.Default()

	// handler.Init()

	handler.InitCreCont(&queue)

	fmt.Println(queue.Go)

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
	ch := make(chan string, 100)
	startTime := time.Now()
	for i := 0; i < 100; i++ {
		// go handler.CreateContainer(strconv.Itoa(i), "go", ch)

		//test time khi nguoi dung request moi tao container
		containerName := xid.New().String()
		go handler.HandlerCompile(containerName, "go", `
		package main\n
		import "fmt"\n
		func main() {\n
			sum := 0\n
			for i := 0; i < 10; i++ {\n
				sum += i\n
			}\n
			fmt.Println(sum)\n
		}`, ch)

	}
	<-ch

	totalTime := time.Since(startTime)
	fmt.Println("Time execute:", totalTime)

	// C222222
	// ch := make(chan string, 50)
	// startTime := time.Now()
	// for i := 0; i < 50; i++ {
	// 	// go handler.CreateContainer(strconv.Itoa(i), "go", ch)

	// 	//test time khi nguoi dung request moi tao container
	// 	go handler.HandlerCompile2("go", `
	// 	package main\n
	// 	import "fmt"\n
	// 	func main() {\n
	// 		sum := 0\n
	// 		for i := 0; i < 10; i++ {\n
	// 			sum += i\n
	// 		}\n
	// 		fmt.Println(sum)\n
	// 	}`, ch, &queue, &m)

	// }
	// <-ch

	// totalTime := time.Since(startTime)
	// fmt.Println("Time execute:", totalTime)

}
