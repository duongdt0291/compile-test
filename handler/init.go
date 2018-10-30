package handler

import (
	"os/exec"

	"dev.hocngay.com/hocngay/compile-test/constant"
	"github.com/sirupsen/logrus"
)

func Init() {
	for _, language := range constant.AllowLanguage {
		buildImages(language)
	}
}

func buildImages(language string) {
	// C và C++ dùng chung image
	if language == "c++" {
		language = "c"
	}

	out, err := exec.Command("docker", "build", "-t", constant.ImageCompilerPrefix+language, constant.LocalBuildDir+"/"+language+"/.").Output()

	if err != nil {
		logrus.Errorf("%s", err)
	}

	logrus.Infof("%s", out)
}
