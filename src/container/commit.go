package container

import (
	log "github.com/sirupsen/logrus"
	"gocker/src/record"
	"os/exec"
	"path"
)

//
// CommitContainer
// @Description: 打包容器
// @param cID
// @param imageName
//
func CommitContainer(cID, imageName string) {
	mntURL := path.Join(record.RootURL, "mnt", cID)
	imageTarURL := "./" + imageName + ".tar"
	if _, err := exec.Command("tar", "-czf", imageTarURL, "-C", mntURL, ".").CombinedOutput(); err != nil {
		log.Errorf("commit container error:%v", err)
	}

}
