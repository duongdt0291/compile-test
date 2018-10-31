package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"dev.hocngay.com/hocngay/compile-test/handler"
	"dev.hocngay.com/hocngay/compile-test/model"
	"github.com/gin-gonic/gin"
	"github.com/panjf2000/ants"
)

var m sync.Mutex
var queue = model.Queue{}

func main() {
	r := gin.Default()

	// handler.Init()

	handler.InitCreCont(&queue)

	fmt.Println(queue.Go)
	// Test with worker pool
	p, _ := ants.NewPoolWithFunc(15, func(i interface{}) error {
		handleTest(i.(*gin.Context))
		return nil
	})

	r.GET("/test", func(ctx *gin.Context) {
		p.Serve(ctx)
	})

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
	//C1 (hiện tại) : khi nào người dùng req mới tạo container
	// ch := make(chan string, 20)

	// for i := 0; i < 20; i++ {
	// 	// go handler.CreateContainer(strconv.Itoa(i), "go", ch)

	// 	//test time khi nguoi dung request moi tao container
	// 	containerName := xid.New().String()
	// 	go handler.HandlerCompile(containerName, "go", `
	// 	package main\n
	// 	import "fmt"\n
	// 	func main() {\n
	// 		sum := 0\n
	// 		for i := 0; i < 10; i++ {\n
	// 			sum += i\n
	// 		}\n
	// 		fmt.Println(sum)\n
	// 	}`, ch)

	// }
	// <-ch

	// C2: tạo sẵn 5 container cho từng ngôn ngữ
	// ch := make(chan string, 20)
	// for i := 0; i < 20; i++ {
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

	//C3 Pool worker
	ch := make(chan string, 1)
	startTime := time.Now()
	go handler.HandlerCompile3("go", `
		package main\n
		import "fmt"\n
		func main() {\n
			sum := 0\n
			for i := 0; i < 10; i++ {\n
				sum += i\n
			}\n
			fmt.Println(sum)\n
		}`, &queue, &m, ch)
	<- ch	

	totalTime := time.Since(startTime)
	fmt.Println("Time 3:", totalTime)
}
