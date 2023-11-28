package main

import (
	"archive/zip"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	router := gin.Default()
	router.Static("/js", "static/js/")
	router.Static("/css", "static/css/")
	router.Static("/img", "static/img/")
	router.LoadHTMLGlob("static/html/*")

	manager := SessionManager{}
	manager.Start()
	downloadManager := DownloadManager{}

	router.GET("/myHost", func(context *gin.Context) {
		info, _ := manager.IsExist(GetClientIP(context))
		context.JSON(OKWithData(info))
	})

	router.GET("/", func(c *gin.Context) {
		clientIp := GetClientIP(c)
		manager.AddSession(GetClientIP(c))
		list := downloadManager.List(clientIp)
		for i := 0; i < len(list); i++ {
			list[i].SrcNickName = manager.GetNickNameByHost(list[i].SrcHost)
		}

		log.Println(list)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"list":    manager.cache,
			"msgList": list,
		})
	})

	router.GET("/deleteFile/:id", func(c *gin.Context) {
		host := GetClientIP(c)
		id := c.Param("id")
		downloadManager.RemoveInfo(host, id)
		//删除目录
		err := os.RemoveAll(GetFileDir(host, id))
		if err != nil {
			c.JSON(FailWithMsg("目录删除失败"))
			log.Println(err)
			return
		}
		c.JSON(OK())
	})

	router.POST("/message", func(c *gin.Context) {
		form, _ := c.MultipartForm()
		msg := form.Value["msg"][0]
		host := GetClientIP(c)
		targetIp := form.Value["ip"][0]
		downloadManager.SaveFile(nil, nil, host, host, targetIp, TYPE_MSG, msg)
	})

	router.POST("/upload", func(c *gin.Context) {
		// Multipart form
		form, _ := c.MultipartForm()
		files := form.File["files[]"]
		targetIp := form.Value["ip"][0]
		host := GetClientIP(c)

		res := downloadManager.SaveFile(files, c, host, host, targetIp, TYPE_FILE, "")

		c.JSON(OKWithMsg(res))
	})

	router.GET("/downloadFile/:id/:fileName", func(c *gin.Context) {
		id := c.Param("id")
		fileName := c.Param("fileName")
		host := GetClientIP(c)
		// 设置响应头信息
		c.Header("Content-Disposition", "attachment; filename="+filepath.Base(fileName))
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Access-Control-Allow-Origin", "*")
		c.File(GetFileDir(host, id) + fileName)
	})

	router.GET("/download/:id", func(c *gin.Context) {
		id := c.Param("id")
		host := GetClientIP(c)
		info := downloadManager.FindById(host, id)
		log.Println(id, host, info)
		filenames := info.FileName
		zipFilePath := GetFileDir(host, id) + info.FileDir + ".zip"

		// 检查文件是否存在
		if _, err := os.Stat(zipFilePath); os.IsNotExist(err) {
			fmt.Println("文件不存在")
			f, err := os.Create(zipFilePath)
			if err != nil {
				log.Println(err)
				c.JSON(FailWithMsg("找不到文件"))
			}

			zipFile := zip.NewWriter(f)
			for i := 0; i < len(filenames); i++ {
				f, err := zipFile.Create(filenames[i])
				if err != nil {
					log.Println(err)
				}
				data, err := os.ReadFile(GetFileDir(host, id) + filenames[i])
				if err != nil {
					log.Println(err)
				}
				_, err = f.Write(data)
				if err != nil {
					log.Println(err)
				}
				zipFile.Flush()
			}
			err = zipFile.Close()
			if err != nil {
				log.Println(err)
			}
		}
		c.Header("Content-Disposition", "attachment; filename="+filepath.Base(id+".zip"))
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Access-Control-Allow-Origin", "*")
		c.File(zipFilePath)
	})

	router.GET("/rename/:name", func(c *gin.Context) {
		name := c.Param("name")
		manager.Rename(GetClientIP(c), name)
		c.JSON(OK())
	})

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	manager.Stop()
}
