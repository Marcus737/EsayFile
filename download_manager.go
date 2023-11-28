package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"os"
)

type DownloadInfo struct {
	TypeName    string   `json:"typeName"`
	SrcHost     string   `json:"srcHost"`
	TargetHost  string   `json:"targetHost"`
	FileName    []string `json:"fileName"`
	Message     string   `json:"message"`
	Progress    uint8    `json:"progress"`
	FileDir     string   `json:"fileDir"`
	SrcNickName string   `json:"SrcNickName"`
}

type DownloadManager struct {
	hostMap map[string][]DownloadInfo
}

const TYPE_FILE = "file"
const TYPE_MSG = "msg"

const fileDir = "trans"

func GetFileDir(host string, id string) string {
	return os.TempDir() + string(os.PathSeparator) + fileDir + string(os.PathSeparator) + host + string(os.PathSeparator) + id + string(os.PathSeparator)
}

func (manager *DownloadManager) SaveFile(files []*multipart.FileHeader, c *gin.Context, host string, src string, dest string, typeName string, message string) string {
	id := GenerateRandomID(16)
	var fileNames []string

	if typeName == TYPE_FILE {
		log.Println(len(files))
		log.Println(files[0].Filename)
		log.Println(src, dest, typeName)
		for i := 0; i < len(files); i++ {
			file := files[i]
			path := GetFileDir(dest, id) + file.Filename
			fileNames = append(fileNames, file.Filename)

			err := c.SaveUploadedFile(file, path)
			log.Println(path)
			if err != nil {
				log.Println("写入文件失败", err)
				return "写入文件失败"
			}
		}
	}
	if manager.hostMap == nil {
		manager.hostMap = make(map[string][]DownloadInfo)
	}
	_, exist := manager.hostMap[dest]
	if !exist {
		manager.hostMap[dest] = make([]DownloadInfo, 0)
	}
	manager.hostMap[dest] = append(manager.hostMap[dest], DownloadInfo{
		TypeName:    typeName,
		SrcHost:     src,
		TargetHost:  dest,
		FileName:    fileNames,
		Message:     message,
		Progress:    0,
		FileDir:     id,
		SrcNickName: src,
	})
	log.Println(manager.hostMap[dest])
	return id
}

func (manager *DownloadManager) List(host string) []DownloadInfo {
	infos, _ := manager.hostMap[host]
	return infos
}

func (manager *DownloadManager) FindById(host string, id string) *DownloadInfo {
	infos := manager.List(host)
	for i := 0; i < len(infos); i++ {
		if infos[i].FileDir == id {
			return &infos[i]
		}
	}
	return nil
}

func (manager *DownloadManager) RemoveInfo(host string, id string) {
	info, exist := manager.hostMap[host]
	if !exist {
		return
	}
	for i := 0; i < len(info); i++ {
		if info[i].FileDir == id {
			info = append(info[:i], info[i+1:]...)
		}
	}
	manager.hostMap[host] = info

}
