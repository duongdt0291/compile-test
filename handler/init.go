package handler

import (
	"fmt"
	"os/exec"

	"dev.hocngay.com/hocngay/compile-test/constant"
	"dev.hocngay.com/hocngay/compile-test/model"

	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
)

func Init() {
	for _, language := range constant.AllowLanguage {
		buildImages(language)
	}
}

func buildImages(language string) {
	fmt.Println("Dang build image")
	// C và C++ dùng chung image
	if language == "c++" {
		language = "c"
	}

	out, err := exec.Command("docker", "build", "-t", constant.ImageCompilerPrefix+language, constant.LocalBuildDir+"/"+language+"/.").Output()

	if err != nil {
		fmt.Println("image:", err)
		logrus.Errorf("%s", err)
	}

	logrus.Infof("%s", out)
}

func InitCreCont(queue *model.Queue) {
	for _, language := range constant.AllowLanguage {
		buildImages2(language, queue)
	}
}

func buildImages2(language string, queue *model.Queue) {
	// C và C++ dùng chung image
	if language == "c++" {
		language = "c"
	}

	// capLanguage := strings.Title(language)

	out, err := exec.Command("docker", "build", "-t", constant.ImageCompilerPrefix+language, constant.LocalBuildDir+"/"+language+"/.").Output()

	if err != nil {
		fmt.Println("image:", err)
		logrus.Errorf("%s", err)
	}

	logrus.Infof("%s", out)

	for i := 0; i < 5; i++ {
		id := xid.New().String()
		_, err := CreateContainer(id, language)
		if err != nil {
			i--
			continue
		}
		newContainer := &model.ContainerInfo{
			Id:        id,
			IsRunning: false,
		}
		queue.Go = append(queue.Go, newContainer)
	}

}
