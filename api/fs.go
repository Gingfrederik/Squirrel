package api

import (
	"fileserver/types"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"fileserver/fs"
)

func (h *Handler) getList(c *gin.Context) {
	fileSystem := fs.GetInstance()
	path := strings.TrimPrefix(c.Param("path"), "/")

	isFile := fileSystem.IsFile(path)

	op := "info"
	if isFile {
		op = c.DefaultQuery("op", "download")
	}

	switch op {
	case "info":
		if isFile {
			fji, err := fileSystem.Info(path)
			if err != nil {
				abortWithError(c, http.StatusBadRequest, err.Error())
				return
			}
			c.JSON(http.StatusOK, fji)
			return
		}
		search := c.Query("search")
		deepStr := c.DefaultQuery("deep", "true")
		deep, _ := strconv.ParseBool(deepStr)

		lrs, err := fileSystem.JSONList(path, search, deep)
		if err != nil {
			abortWithError(c, http.StatusBadRequest, err.Error())
			return
		}

		res := types.Response{
			Status:  0,
			Message: "get file list",
			Data:    lrs,
		}

		c.JSON(http.StatusOK, res)
	case "download":
		if isFile {
			fullpath := filepath.Join(fileSystem.Root, path)
			c.FileAttachment(fullpath, filepath.Base(path))
			return
		}
	default:
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func (h *Handler) upload(c *gin.Context) {
	fileSystem := fs.GetInstance()
	path := c.Param("path")
	file, _ := c.FormFile("file")

	fjis := make([]*fs.FileJSONInfo, 0)
	uploadInfo := &fs.UploadInfo{
		Path: path,
	}
	if file != nil {
		uploadInfo.FileHeader = file
	}

	result, err := fileSystem.UploadOrMkdir(uploadInfo)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	if !result {
		abortWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	uploadFilePath := filepath.Join(path, file.Filename)
	fji, err := fileSystem.Info(uploadFilePath)
	if err != nil {
		abortWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	fjis = append(fjis, fji)

	res := types.Response{
		Status:  1,
		Message: "upload success",
		Data:    fjis,
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) delete(c *gin.Context) {
	fileSystem := fs.GetInstance()
	path := c.Param("path")

	_, err := fileSystem.Info(path)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	err = fileSystem.Delete(path)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	res := types.Response{
		Status:  1,
		Message: "delete success",
	}

	c.JSON(http.StatusOK, res)
}
