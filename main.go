package main

import (
	"fmt"
	"os"
	"sync"

	"git.hocngay.com/hocngay/compile-test/handler"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

var m sync.Mutex
var containers []string = []string{}

func init() {
	fmt.Println(1)
	containers = handler.InitCreCont(containers)
	fmt.Println(containers)
}

func main() {
	r := gin.Default()

	go TrackingContainers(containers)

	r.GET("/containers", handleTestContainers)

	// // Test with worker pool
	// p, _ := ants.NewPoolWithFunc(15, func(i interface{}) error {
	// 	handleTest(i.(*gin.Context))
	// 	return nil
	// })

	// r.GET("/test", func(ctx *gin.Context) {
	// 	p.Serve(ctx)
	// })

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

// C2: tạo sẵn 5 container cho từng ngôn ngữ
func handleTestContainers(ctx *gin.Context) {
	ch := make(chan string, 10)
	for i := 0; i < 10; i++ {
		go handler.HandlerCompile2("go", `
		package main\n
		import "fmt"\n
		func main() {\n
			sum := 0\n
			for i := 0; i < 10; i++ {\n
				sum += i\n
			}\n
			fmt.Println(sum)\n
		}`, ch, containers, &m)
	}
	for v := range <-ch {
		fmt.Println(v)
	}
	fmt.Println(containers)
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

	//C3 Pool worker
	// ch := make(chan string, 1)
	// startTime := time.Now()
	// go handler.HandlerCompile3("go", `
	// 	package main\n
	// 	import "fmt"\n
	// 	func main() {\n
	// 		sum := 0\n
	// 		for i := 0; i < 10; i++ {\n
	// 			sum += i\n
	// 		}\n
	// 		fmt.Println(sum)\n
	// 	}`, &queue, &m, ch)
	// <-ch

	// totalTime := time.Since(startTime)
	// fmt.Println("Time 3:", totalTime)
}

func TrackingContainers(containers []string) {
	for {
		if len(containers) < 5 {
			newContainerName := xid.New().String()
			_, err := handler.CreateContainer(newContainerName, "go")
			if err == nil {
				containers = append(containers, newContainerName)
			}
		}
	}
}
