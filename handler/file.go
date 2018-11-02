package handler

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"git.hocngay.com/hocngay/compile-test/constant"
	"github.com/rs/xid"
)

// Trả về file name đã được khởi tạo
func CreateFileCompile(language, content string) (string, string, error) {
	var localFolderPath, localFilePath, fileName string

	folderID := xid.New().String()
	localFolderPath = constant.LocalTempDir + "/" + language + "/" + folderID

	switch language {
	case "c":
		fileName = "main.c"
	case "c++":
		fileName = "main.cpp"
	case "go":
		fileName = "main.go"
	case "java":
		fileName = "Main.java"
	case "node":
		fileName = "main.js"
	case "php":
		fileName = "main.php"
	case "python":
		fileName = "main.py"
	case "ruby":
		fileName = "main.rb"
	}

	localFilePath = localFolderPath + "/" + fileName

	if _, err := os.Stat(localFolderPath); os.IsNotExist(err) {
		os.MkdirAll(localFolderPath, os.FileMode(0777))
	}

	file, err := os.Create(localFilePath)
	if err != nil {
		return localFilePath, fileName, err
	}
	defer file.Close()

	code := []byte(content)
	err = ioutil.WriteFile(localFilePath, code, 0644)

	return localFilePath, fileName, err
}

// Copy file vào trong container
func copyFileCompile(localFilePath, containerName, containerRunDir string) error {
	_, err := exec.Command("docker", "cp", localFilePath, fmt.Sprintf("%s:%s", containerName, containerRunDir)).CombinedOutput()
	return err
}
