package handler

import (
	"log"
	"os/exec"
	"regexp"
)

// Khởi tạo container
func CreateContainer(containerName, language string) (string, error) {

	// // C và C++ dùng chung image
	// if language == "c++" {
	// 	language = "c"
	// }

	isExist := isContainerExist(containerName)
	if isExist {
		return containerName, nil
	}
	//constant.ImageCompilerPrefix+language
	_, err := exec.Command("docker", "run", "-id", "--name", containerName, "compiler-go").Output()
	if err != nil {
		log.Println("err: ", err)
	}
	// ch <- "done"
	return containerName, err
}

// RemoveContainer xóa container
func RemoveContainer(containerName string) (string, error) {
	_, err := exec.Command("docker", "kill", containerName).CombinedOutput()
	return containerName, err
}

// Kiểm tra container đã tồn tại chưa
func isContainerExist(containerName string) (isExist bool) {
	out, _ := exec.Command("docker", "inspect", "--format=\"{{.Name}}\"", containerName).CombinedOutput()

	regexContainerExist, _ := regexp.Compile("No such object: " + containerName)
	isExist = !regexContainerExist.MatchString(string(out))

	return isExist
}
