package handler

import (
	"os"
	"os/exec"

	"dev.hocngay.com/hocngay/compile-test/constant"
	"github.com/kr/pty"
)

func (wp *WsPty) StartManual(param []string) error {
	var err error
	wp.Cmd = exec.Command(param[0], param[1:]...)
	wp.Pty, err = pty.Start(wp.Cmd)
	return err
}

func (wp *WsPty) Stop() {
	wp.Pty.Close()
	wp.Cmd.Wait()
}

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

// func (wp *WsPty) HandlerCompile(containerName, language, content string) ([]byte, error) {
// 	isAllow := checkLanguage(language)
// 	if !isAllow {
// 		return nil, errors.New("Không hỗ trợ ngôn ngữ này")
// 	}

// 	var localFilePath, fileName string
// 	ch1 := make(chan string, 1)
// 	ch2 := make(chan string, 1)

// 	// Tạo file trên local từ code
// 	go func() {
// 		var err error
// 		localFilePath, fileName, err = CreateFileCompile(language, content)

// 		if err != nil {
// 			ch1 <- err.Error()
// 		}
// 		ch1 <- "done"
// 	}()

// 	// Tạo container
// 	go func() {
// 		_, err := CreateContainer(containerName, language)
// 		if err != nil {
// 			ch2 <- err.Error()
// 		}
// 		ch2 <- "done"
// 	}()

// 	result1 := <-ch1
// 	result2 := <-ch2

// 	if result1 == "done" && result2 == "done" {
// 		err := copyFileCompile(localFilePath, containerName, constant.ContainerRunDir)
// 		if err != nil {
// 			log.Println(err)
// 		}

// 		// Thực thi file
// 		switch language {

// 		case "c":
// 			// Bỏ extesion .c khỏi tên file
// 			binaryFileName := fileName[:len(fileName)-len(".c")]

// 			// Biên dịch file c sang binary
// 			result, err := exec.Command("docker", "exec", containerName, "gcc", fileName, "-o", binaryFileName).CombinedOutput()

// 			if err != nil {
// 				return result, errors.New("Không thể biên dịch code")
// 			}

// 			wp.StartManual([]string{"docker", "exec", "-it", containerName, "./" + binaryFileName})

// 		case "c++":
// 			// Bỏ extesion .cpp khỏi tên file
// 			binaryFileName := fileName[:len(fileName)-len(".cpp")]

// 			// Biên dịch file cpp sang binary
// 			result, err := exec.Command("docker", "exec", containerName, "g++", fileName, "-o", binaryFileName).CombinedOutput()

// 			if err != nil {
// 				return result, errors.New("Không thể biên dịch code")
// 			}

// 			wp.StartManual([]string{"docker", "exec", "-it", containerName, "./" + binaryFileName})

// 		case "go":
// 			wp.StartManual([]string{"docker", "exec", "-it", containerName, "go", "run", fileName})

// 		case "java":
// 			// Bỏ extesion .java khỏi tên file
// 			binaryFileName := fileName[:len(fileName)-len(".java")]

// 			// Biên dịch file java sang class
// 			result, err := exec.Command("docker", "exec", containerName, "javac", "-encoding", "UTF-8", fileName).CombinedOutput()

// 			if err != nil {
// 				return result, errors.New("Không thể biên dịch code")
// 			}

// 			wp.StartManual([]string{"docker", "exec", "-it", containerName, "java", binaryFileName})

// 		case "node":
// 			wp.StartManual([]string{"docker", "exec", "-it", containerName, "node", fileName})

// 		case "php":
// 			wp.StartManual([]string{"docker", "exec", "-it", containerName, "php", fileName})

// 		case "python":
// 			wp.StartManual([]string{"docker", "exec", "-it", containerName, "python", fileName})

// 		case "ruby":
// 			wp.StartManual([]string{"docker", "exec", "-it", containerName, "ruby", fileName})
// 		}
// 	}

// 	return nil, nil
// }

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
