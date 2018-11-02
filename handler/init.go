package handler

import (
	"fmt"
	"os/exec"

	"git.hocngay.com/hocngay/compile-test/constant"

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

func InitCreCont(containers []string) {
	for _, language := range constant.AllowLanguage {
		buildImages2(language, containers)
	}
}

func buildImages2(language string, containers []string) {
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

	for i := 0; i < 5; i++ {
		id := xid.New().String()
		_, err := CreateContainer(id, language)
		if err != nil {
			i--
			continue
		}
		containers = append(containers, id)
		fmt.Println(i, containers)
	}
}
