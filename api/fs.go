package api

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"fileserver/fs"
)

func (h *Handler) getList(c *gin.Context) {
	fileSystem := fs.GetInstance()
	path := strings.TrimPrefix(c.Param("path"), "/")
	isFile := fileSystem.IsFile(path)
	switch op := c.Query("op"); op {
	case "download":

		if isFile {
			fullpath := filepath.Join(fileSystem.Root, path)
			c.FileAttachment(fullpath, filepath.Base(path))
			return
		}
	case "":
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
		lrs, err := fileSystem.JSONList(path, search)
		if err != nil {
			abortWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, lrs)

	default:
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func (h *Handler) upload(c *gin.Context) {
	fileSystem := fs.GetInstance()
	path := c.Param("path")

	file, _ := c.FormFile("file")

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

	c.String(http.StatusOK, "Uploaded...")
}

func (h *Handler) delete(c *gin.Context) {
	fileSystem := fs.GetInstance()
	path := c.Param("path")
	isFile := fileSystem.IsFile(path)
	if !isFile {
		abortWithError(c, http.StatusBadRequest, "can only delete file")
		return
	}

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

	c.String(http.StatusOK, "Deleted...")
}
