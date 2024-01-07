package files

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gochat/middleware"
	"gochat/models"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var DB *gorm.DB

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func FileRoutes(router *gin.RouterGroup, db *gorm.DB) {
	DB = db
	r := router.Group("/file")
	r.POST("", saveFile)
	r.POST("/link/", linkFile)
	r.GET("/img/:id", getFile)
	r.GET("/download/:id", downloadFile)
	r.GET("/info/:id", fileInfo)
	r.GET("/getFiles/:id", getFiles)
}

func saveFile(ctx *gin.Context) {
	var (
		form     Form
		response ResponseFile
		file     = &models.File{}
		db       = ctx.Value(middleware.KeyDB).(*gorm.DB)
		user     = ctx.Value(middleware.KeyAuthUser).(*models.User)
	)
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("file read error: %s", err.Error()))
		return
	}

	path, err := os.Getwd()
	if err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}
	t := time.Now()
	if form.ID == "" {
		if form.Type == "user" {
			form.ID = strconv.Itoa(user.ID)
		}
	}

	for _, formFile := range form.Files {
		fullPath := filepath.Join(path, "files", form.Type, form.ID, form.Path)
		if form.Type == "chat" && form.Path == "file" {
			chatID := form.ID
			file.ParentID, _ = strconv.Atoi(chatID)
			fullPath = filepath.Join(path, "files", form.Type, chatID, form.Path)
		} else {
			file.ParentID, _ = strconv.Atoi(form.ID)
		}
		if _, err = os.Stat(fullPath); os.IsNotExist(err) {
			err = os.MkdirAll(fullPath, os.ModeDir)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, err.Error())
				return
			}
		}

		file.Path = formFile.Filename
		//change filename if it already exists by path
		sb := strings.Builder{}
		for i := 0; i < 20; i++ {
			sb.WriteByte(charset[rand.Intn(len(charset))])
		}
		file.Extension = filepath.Ext(file.Path)
		formFile.Filename = sb.String() + file.Extension
		if _, err = os.Stat(filepath.Join(fullPath, formFile.Filename)); err == nil {
			sb = strings.Builder{}
			for i := 0; i < 20; i++ {
				sb.WriteByte(charset[rand.Intn(len(charset))])
			}
			formFile.Filename = sb.String() + filepath.Ext(formFile.Filename)
		}
		file.Guid = filepath.Join("files", form.Type, strconv.Itoa(file.ParentID), form.Path)
		file.Type = form.Type
		file.Name = formFile.Filename
		file.Fullname = filepath.Join(fullPath, formFile.Filename)
		file.Size = int(formFile.Size)
		file.MimeType = formFile.Header.Get("Content-Type")
		imageIndex := strings.LastIndex(file.MimeType, "/")
		file.Mtype = file.MimeType[:imageIndex]
		file.Created = t

		if err = ctx.SaveUploadedFile(formFile, filepath.Join(fullPath, formFile.Filename)); err != nil {
			ctx.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
		if db.Create(file).Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Can't save file"})
			return
		}
		response.FileID = file.ID
		response.FileName = file.Name
		response.Extension = file.Extension
		response.Type = file.Mtype
		response.Guid = file.Guid
		if form.Type == "chat" && form.Path == "avatar" {
			db.Model(&models.Chat{ID: file.ParentID}).Update("logo", file.ID)
		}
	}
	response.file = form

	ctx.JSON(http.StatusOK, &response)
}

func linkFile(ctx *gin.Context) {
	var (
		arr []int
		db  = ctx.Value(middleware.KeyDB).(*gorm.DB)
	)
	b, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}
	err = json.Unmarshal(b, &arr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}
	db.Model(&models.Message{ID: arr[1]}).Update("attachment", arr[0])
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}
}

func getFile(ctx *gin.Context) {
	var (
		file   = &models.File{}
		fileID = ctx.Param("id")
		db     = ctx.Value(middleware.KeyDB).(*gorm.DB)
	)
	query := db.Select("*").
		Table("files").
		Where("id = ?", fileID).Order("id desc").Limit(1)

	if err := query.Scan(&file).Error; err != nil || file.Path == "" {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	fileName := file.Guid + "/" + file.Name

	ctx.JSON(http.StatusOK, fileName)
}

func downloadFile(ctx *gin.Context) {
	var (
		param   = ctx.Param("id")
		id, err = strconv.Atoi(param)
		db      = ctx.Value(middleware.KeyDB).(*gorm.DB)
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	var file = models.File{ID: id}
	if db.First(&file).Error != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	fileBytes, err := ioutil.ReadFile(file.Fullname)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	info, err := os.Stat(file.Fullname)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s := info.Size()
	size := strconv.Itoa(int(s))
	f, err := os.Open(file.Fullname)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}
	defer f.Close()

	var r io.Reader
	r = f
	println(r)
	println(fileBytes)

	ctx.Writer.Header().Set("Content-Type", file.MimeType)
	ctx.Writer.Header().Set("Content-Disposition", "attachment; filename="+file.Name)
	ctx.Writer.Header().Set("Content-Length", size)
	_, err = io.Copy(ctx.Writer, r)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	return
}

func fileInfo(ctx *gin.Context) {
	var (
		param   = ctx.Param("id")
		id, err = strconv.Atoi(param)
		db      = ctx.Value(middleware.KeyDB).(*gorm.DB)
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	var file = models.File{ID: id}
	if db.First(&file).Error != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, &file.Path)
}

func getFiles(ctx *gin.Context) {
	var (
		files  = make(ResponseFiles, 0)
		chatID = ctx.Param("id")
		db     = ctx.Value(middleware.KeyDB).(*gorm.DB)
	)
	query := db.Select("*").
		Table("files").
		Where("parent_id = ?", chatID).
		Where("type = 'chat'")

	if err := query.Scan(&files).Error; err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, &files)
}
