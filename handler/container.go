package handler

import (
	"os/exec"
	"regexp"

	"dev.hocngay.com/hocngay/compile-test/constant"
)

// Khởi tạo container
func CreateContainer(containerName, language string, ch chan string) (string, error) {

	// C và C++ dùng chung image
	if language == "c++" {
		language = "c"
	}

	isExist := isContainerExist(containerName)
	if isExist {
		return containerName, nil
	}

	_, err := exec.Command("docker", "run", "-id", "--rm", "--name", containerName, constant.ImageCompilerPrefix+language).Output()
	ch <- "done"
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
