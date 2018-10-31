package handler

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"dev.hocngay.com/hocngay/compile-test/constant"
	"dev.hocngay.com/hocngay/compile-test/model"
	"github.com/kr/pty"
	"github.com/rs/xid"
)

func StartManual(param []string) error {
	var err error
	cmd := exec.Command(param[0], param[1:]...)
	_, err = pty.Start(cmd)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// func (wp *WsPty) Stop() {
// 	wp.Pty.Close()
// 	wp.Cmd.Wait()
// }

type WsPty struct {
	Cmd         *exec.Cmd
	Pty         *os.File
	ID          string
	Language    string
	Content     string
	Type        string
	IsReConnect bool
	Payload     []byte
}

func HandlerCompile(containerName, language, content string, ch chan string) ([]byte, error) {
	startTime := time.Now()
	isAllow := checkLanguage(language)
	if !isAllow {
		return nil, errors.New("Không hỗ trợ ngôn ngữ này")
	}

	var localFilePath, fileName string
	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)

	// Tạo file trên local từ code
	go func() {
		var err error
		localFilePath, fileName, err = CreateFileCompile(language, content)

		if err != nil {
			fmt.Println("Loi tao file")
			ch1 <- err.Error()
		}
		ch1 <- "done"
	}()

	// Tạo container
	go func() {

		_, err := CreateContainer(containerName, language)
		if err != nil {
			fmt.Println("Loi tao container:", err)
			ch2 <- err.Error()
		}
		ch2 <- "done"
	}()

	result1 := <-ch1
	result2 := <-ch2

	if result1 == "done" && result2 == "done" {
		err := copyFileCompile(localFilePath, containerName, constant.ContainerRunDir)
		if err != nil {
			fmt.Println("Loi copy file")
			log.Println(err)
		}

		// Thực thi file
		switch language {

		case "c":
			// Bỏ extesion .c khỏi tên file
			binaryFileName := fileName[:len(fileName)-len(".c")]

			// Biên dịch file c sang binary
			result, err := exec.Command("docker", "exec", containerName, "gcc", fileName, "-o", binaryFileName).CombinedOutput()

			if err != nil {
				return result, errors.New("Không thể biên dịch code")
			}

			StartManual([]string{"docker", "exec", "-it", containerName, "./" + binaryFileName})

		case "c++":
			// Bỏ extesion .cpp khỏi tên file
			binaryFileName := fileName[:len(fileName)-len(".cpp")]

			// Biên dịch file cpp sang binary
			result, err := exec.Command("docker", "exec", containerName, "g++", fileName, "-o", binaryFileName).CombinedOutput()

			if err != nil {
				return result, errors.New("Không thể biên dịch code")
			}

			StartManual([]string{"docker", "exec", "-it", containerName, "./" + binaryFileName})

		case "go":

			StartManual([]string{"docker", "exec", "-it", containerName, "go", "run", fileName})

		case "java":
			// Bỏ extesion .java khỏi tên file
			binaryFileName := fileName[:len(fileName)-len(".java")]

			// Biên dịch file java sang class
			result, err := exec.Command("docker", "exec", containerName, "javac", "-encoding", "UTF-8", fileName).CombinedOutput()

			if err != nil {
				return result, errors.New("Không thể biên dịch code")
			}

			StartManual([]string{"docker", "exec", "-it", containerName, "java", binaryFileName})

		case "node":
			StartManual([]string{"docker", "exec", "-it", containerName, "node", fileName})

		case "php":
			StartManual([]string{"docker", "exec", "-it", containerName, "php", fileName})

		case "python":
			StartManual([]string{"docker", "exec", "-it", containerName, "python", fileName})

		case "ruby":
			StartManual([]string{"docker", "exec", "-it", containerName, "ruby", fileName})
		}
	}
	totalTime := time.Since(startTime)
	fmt.Println("Time 1:", totalTime)

	ch <- "done"
	return nil, nil
}

func checkLanguage(language string) bool {
	isAllow := false

	for _, al := range constant.AllowLanguage {
		if al == language {
			isAllow = true
			break
		}
	}

	return isAllow
}

func HandlerCompile2(language, content string, ch chan string, queue *model.Queue, m *sync.Mutex) ([]byte, error) {
	startTime := time.Now()

	// fmt.Println(queue[language])
	isAllow := checkLanguage(language)
	if !isAllow {
		return nil, errors.New("Không hỗ trợ ngôn ngữ này")
	}

	var localFilePath, fileName string
	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)

	// Tạo file trên local từ code
	go func() {
		var err error
		localFilePath, fileName, err = CreateFileCompile(language, content)

		if err != nil {
			fmt.Println("Loi tao file")
			ch1 <- err.Error()
		}
		ch1 <- "done"
	}()

	lenCon := len(queue.Go)
	containerName := xid.New().String()
	count := 0

	m.Lock()
	for i := 0; i < lenCon; i++ {
		if queue.Go[i].IsRunning == false {
			containerName = queue.Go[i].Id
			queue.Go[i].IsRunning = true
			break
		}
		count++
	}
	m.Unlock()

	if count == lenCon {
		// Tạo container
		go func() {
			_, err := CreateContainer(containerName, language)
			if err != nil {
				fmt.Println("Loi tao container:", err)
				ch2 <- err.Error()
			} else {
				newContainer := &model.ContainerInfo{
					Id:        containerName,
					IsRunning: true,
				}
				queue.Go = append(queue.Go, newContainer)
			}
			ch2 <- "done"
		}()
	} else {
		ch2 <- "done"
		newContainerId := xid.New().String()
		//TODO: có cần chạy while để chắc chắc container mới được tạo ko?
		go func() {
			CreateContainer(newContainerId, language)
			newContainer := &model.ContainerInfo{
				Id:        newContainerId,
				IsRunning: false,
			}
			queue.Go = append(queue.Go, newContainer)
		}()
	}

	result1 := <-ch1
	result2 := <-ch2

	if result1 == "done" && result2 == "done" {
		err := copyFileCompile(localFilePath, containerName, constant.ContainerRunDir)
		if err != nil {
			fmt.Println("Loi copy file")
			log.Println(err)
		}

		// Thực thi file
		switch language {
		case "go":
			StartManual([]string{"docker", "exec", "-it", containerName, "go", "run", fileName})

		}
	}
	// RemoveContainer(containerName)
	totalTime := time.Since(startTime)
	fmt.Println("Time 2:", totalTime)
	ch <- "done"
	return nil, nil
}

func HandlerCompile3(language, content string, queue *model.Queue, m *sync.Mutex, ch chan string) ([]byte, error) {

	// fmt.Println(queue[language])
	isAllow := checkLanguage(language)
	if !isAllow {
		return nil, errors.New("Không hỗ trợ ngôn ngữ này")
	}

	var localFilePath, fileName string
	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)

	// Tạo file trên local từ code
	go func() {
		var err error
		localFilePath, fileName, err = CreateFileCompile(language, content)

		if err != nil {
			fmt.Println("Loi tao file")
			ch1 <- err.Error()
		}
		ch1 <- "done"
	}()

	lenCon := len(queue.Go)
	containerName := xid.New().String()
	count := 0

	m.Lock()
	for i := 0; i < lenCon; i++ {
		if queue.Go[i].IsRunning == false {
			containerName = queue.Go[i].Id
			queue.Go[i].IsRunning = true
			break
		}
		count++
	}
	m.Unlock()

	if count == lenCon {
		// Tạo container
		go func() {
			_, err := CreateContainer(containerName, language)
			if err != nil {
				fmt.Println("Loi tao container:", err)
				ch2 <- err.Error()
			} else {
				newContainer := &model.ContainerInfo{
					Id:        containerName,
					IsRunning: true,
				}
				queue.Go = append(queue.Go, newContainer)
			}
			ch2 <- "done"
		}()
	} else {
		ch2 <- "done"
		newContainerId := xid.New().String()
		//TODO: có cần chạy while để chắc chắc container mới được tạo ko?
		go func() {
			CreateContainer(newContainerId, language)
			newContainer := &model.ContainerInfo{
				Id:        newContainerId,
				IsRunning: false,
			}
			queue.Go = append(queue.Go, newContainer)
		}()
	}

	result1 := <-ch1
	result2 := <-ch2

	if result1 == "done" && result2 == "done" {
		err := copyFileCompile(localFilePath, containerName, constant.ContainerRunDir)
		if err != nil {
			fmt.Println("Loi copy file")
			log.Println(err)
		}

		// Thực thi file
		switch language {
		case "go":
			StartManual([]string{"docker", "exec", "-it", containerName, "go", "run", fileName})

		}
	}
	ch <- "done"
	return nil, nil
}
