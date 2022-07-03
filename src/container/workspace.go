package container

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gocker/src/utils"
	"os"
	"os/exec"
	"path"
	"strings"
)

//
// NewWorkSpace
// @Description: 创建新的文件workspace
// @param rootURL
// @param mntURL
//
func NewWorkSpace(rootURL, mntURL, ImageTarPath, volume, cID string) {
	imageName := verifyTarPath(ImageTarPath)
	if imageName == "" {
		return
	}
	//创建init只读层
	createReadOnlyLayer(rootURL, ImageTarPath, imageName)
	//创建读写层
	createWriteLayer(rootURL, cID)
	//创建挂载点
	createMountPoint(rootURL, mntURL, imageName, cID)
	//挂载数据卷
	if volume != "" {
		volumeURLs, err := volumeUrlExtract(volume)
		if err != nil {
			log.Warn(err)
			return
		}
		mountVolume(mntURL, volumeURLs)
		log.Infof("workspace mount volume success")
	}
}

//
// createReadOnlyLayer
// @Description: 解压缩镜像tar，创建只读层
// @param rootURL
// @param ImageTarPath
// @param imageName
//
func createReadOnlyLayer(rootURL, ImageTarPath, imageName string) {
	imageDir := path.Join(rootURL, "diff", imageName)
	if has, err := utils.DirOrFileExist(imageDir); err == nil && !has {
		//不存在文件
		if err := os.MkdirAll(imageDir, 0777); err != nil {
			log.Errorf("create readonly layer dir error:%v", err)
		}
	}
	if _, err := exec.Command("tar", "-xvf", ImageTarPath, "-C", imageDir).CombinedOutput(); err != nil {
		log.Errorf("exec tar xvf command error:%v", err)
	}
}

//
// createWriteLayer
// @Description: 创建读写层
// @param rootURL
// @param cID
//
func createWriteLayer(rootURL, cID string) {
	writeURL := path.Join(rootURL, "diff", cID+"_writeLayer")
	if has, err := utils.DirOrFileExist(writeURL); err == nil && has {
		log.Infof("write layer dir already exist,recreate new write layer dir")
		deleteWriteLayer(rootURL, cID)
	}
	if err := os.Mkdir(writeURL, 0777); err != nil {
		log.Errorf("create write layer error:%v", err)
	}
}

//
// createMountPoint
// @Description:创建挂载点
// @param rootURL
// @param mntURL
// @param imageName
// @param cID
//
func createMountPoint(rootURL, mntURL, imageName, cID string) {
	if has, err := utils.DirOrFileExist(mntURL); err == nil && has {
		log.Info("mnt dir already exist,recreate mnt dir")
		deleteMountPoint(mntURL)
	}
	if err := os.MkdirAll(mntURL, 0777); err != nil {
		log.Errorf("create mnt dir error:%v", err)
	}
	//将读写层和只读层重新挂载
	writeURL := path.Join(rootURL, "diff", cID+"_writeLayer")
	imageDir := path.Join(rootURL, "diff", imageName)
	dirs := "dirs=" + writeURL + ":" + imageDir
	//todo mount命令详解
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "mnt_"+cID[:4], mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Errorf("mount aufs error:%v", err)
	}
}

//
// mountVolume
// @Description: 挂载数据卷
// @param mntURL
// @param volumeURL
//
func mountVolume(mntURL string, volumeURL []string) {
	//创建宿主机文件目录
	parentURL, containerURL := volumeURL[0], path.Join(mntURL, volumeURL[1])
	//创建文件夹
	if has, err := utils.DirOrFileExist(parentURL); err == nil && !has {
		if err := os.Mkdir(parentURL, 0777); err != nil {
			log.Errorf("create parent volume dir error:%v", err)
			return
		}
	}
	//在容器中创建挂载点目录
	if has, err := utils.DirOrFileExist(containerURL); err == nil && !has {
		if err := os.RemoveAll(containerURL); err != nil {
			log.Errorf("delete container mount volume dir error:%v", err)
			return
		}
	}
	if err := os.MkdirAll(containerURL, 0777); err != nil {
		log.Errorf("create container mount volume dir error:%v", err)
		return
	}
	//开始挂载
	dirs := "dirs=" + parentURL
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "gocker", containerURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("mount volume failed,the error:%v", err)
	}
}

//
// verifyTarPath
// @Description: 验证镜像tar路径
// @param ImageTarPath
// @return string
//
func verifyTarPath(ImageTarPath string) string {
	if has, err := utils.DirOrFileExist(ImageTarPath); err != nil {
		log.Errorf("verifyTarPath error:%v", err)
		return ""
	} else if err == nil && !has {
		log.Errorf("can not found this image through this path:%v", err)
		return ""
	}
	paths := strings.Split(ImageTarPath, "/")
	tarFilename := paths[len(paths)-1]
	if !strings.HasSuffix(tarFilename, "tar") {
		log.Errorf("ImageTarPath has no tar file")
	}
	return strings.Split(tarFilename, ".")[0]
}

//
// volumeUrlExtract
// @Description: 解析volume字符串
// @param volume
// @return []string
// @return error
//
func volumeUrlExtract(volume string) ([]string, error) {
	volumeArray := strings.Split(volume, ":")
	//暂只支持一对数据卷挂载
	if len(volumeArray) != 2 || volumeArray[0] == "" || volumeArray[1] != "" {
		return nil, fmt.Errorf("mount volume args error")
	}
	return volumeArray, nil
}

//
// DeleteWorkspace
// @Description: 容器删除时调用此函数删除容器工作空间
// @param rootURL
// @param mntURL
// @param volume
// @param cID
//
func DeleteWorkspace(rootURL, mntURL, volume, cID string) {
	if volume != "" {
		volumeURLs, err := volumeUrlExtract(volume)
		if err != nil {
			log.Warn(err)
			deleteMountPoint(mntURL)
		} else {
			deleteMountPointWithVolume(mntURL, volumeURLs)
		}
	} else {
		deleteMountPoint(mntURL)
	}
	deleteWriteLayer(rootURL, cID)
}

//
// deleteMountPoint
// @Description: 取消挂载点
// @param mntURL
//
func deleteMountPoint(mntURL string) {
	cmd := exec.Command("unmount", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("unmount mnt dir error:%v", err)
	}
	if err := os.RemoveAll(mntURL); err != nil {
		log.Errorf("delete mnt dir error:%v")
	}
}

//
// deleteMountPointWithVolume
// @Description: 取消挂载点，取消volume挂载
// @param mntURL
// @param volumeURL
//
func deleteMountPointWithVolume(mntURL string, volumeURL []string) {
	containerURL := path.Join(mntURL, volumeURL[1])
	cmd := exec.Command("unmount", containerURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("unmount volume dir error:%v", err)
	}
	deleteMountPoint(mntURL)
}

//
// deleteWriteLayer
// @Description: 删除读写层
// @param rootURL
// @param cID
//
func deleteWriteLayer(rootURL, cID string) {
	writeURL := path.Join(rootURL, "diff", cID+"_writeLayer")
	if err := os.RemoveAll(writeURL); err != nil {
		log.Errorf("delete write layer error:%v", err)
	}
}
