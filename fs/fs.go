package fs

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type fileSystem struct {
	Root    string
	indexes []IndexFileItem
}

var instance *fileSystem
var once sync.Once

func New(root string) {
	if root == "" {
		root = "./"
	}
	root = filepath.ToSlash(root)
	if !strings.HasSuffix(root, "/") {
		root = root + "/"
	}
	log.Infof("root path: %s", root)

	instance = &fileSystem{
		Root: root,
	}

	go func() {
		time.Sleep(1 * time.Second)
		for {
			startTime := time.Now()
			log.Info("Started making search index")
			instance.makeIndex()
			log.Infof("Completed search index in %v", time.Since(startTime))
			//time.Sleep(time.Second * 1)
			time.Sleep(time.Minute * 10)
		}
	}()
}

func GetInstance() *fileSystem {
	return instance
}

func (s *fileSystem) Mkdir(path string) (result bool, err error) {
	dir := filepath.Dir(path)

	name := filepath.Base(path)
	if err = checkFilename(name); err != nil {
		return
	}
	err = os.Mkdir(filepath.Join(s.Root, dir, name), 0755)
	if err != nil {
		return
	}

	return true, nil
}

func (s *fileSystem) Delete(path string) (err error) {
	err = os.RemoveAll(filepath.Join(s.Root, path))
	if err != nil {
		pathErr, ok := err.(*os.PathError)
		if ok {
			err = fmt.Errorf(pathErr.Op + " " + path + ": " + pathErr.Err.Error())
			return
		}
		return
	}
	return nil
}

func (s *fileSystem) UploadOrMkdir(data *UploadInfo) (reuslt bool, err error) {
	dirpath := filepath.Join(s.Root, data.Path)
	var file multipart.File
	if data.FileHeader != nil {
		file, err = data.FileHeader.Open()
		if err != nil {
			return false, err
		}
		defer func() {
			file.Close()
		}()
	}
	info, err := os.Stat(dirpath)

	if os.IsNotExist(err) {
		if err = os.MkdirAll(dirpath, os.ModePerm); err != nil {
			log.Errorf("Create directory: %s", err)
			return
		}
	}

	if file == nil { // only mkdir
		if info != nil {
			return false, errors.New("directory exist")
		}
		return true, nil
	}

	if err != nil {
		log.Errorf("Parse form file: %s", err)
		return
	}

	if err = checkFilename(data.FileHeader.Filename); err != nil {
		return
	}

	dstPath := filepath.Join(dirpath, data.FileHeader.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		log.Errorf("Create file: %s", err)
		return
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		log.Errorf("Handle upload file: %s", err)
		return
	}

	return true, nil
}

func (s *fileSystem) Info(path string) (fji *FileJSONInfo, err error) {
	relPath := filepath.Join(s.Root, path)

	fi, err := os.Stat(relPath)
	if err != nil {
		return
	}
	fji = &FileJSONInfo{
		Name:    fi.Name(),
		Size:    fi.Size(),
		Path:    relPath,
		ModTime: fi.ModTime().UTC(),
	}
	ext := filepath.Ext(path)
	switch ext {
	case ".md":
		fji.Type = "markdown"
	case ".apk":
		fji.Type = "apk"
	case "":
		fji.Type = "dir"
	default:
		fji.Type = "text"
	}

	return
}

func (s *fileSystem) JSONList(requestPath string, search string, deep bool) (lrs []HTTPFileInfo, err error) {
	localPath := filepath.Join(s.Root, requestPath)

	// path string -> info os.FileInfo
	fileInfoMap := make(map[string]os.FileInfo, 0)

	if search != "" {
		results := s.findIndex(search)
		if len(results) > 50 { // max 50
			results = results[:50]
		}
		for _, item := range results {
			if filepath.HasPrefix(item.Path, requestPath) {
				fileInfoMap[item.Path] = item.Info
			}
		}
	} else {
		var infos []os.FileInfo
		infos, err = ioutil.ReadDir(localPath)
		if err != nil {
			return
		}
		for _, info := range infos {
			fileInfoMap[filepath.Join(requestPath, info.Name())] = info
		}
	}

	// turn file list -> json
	lrs = make([]HTTPFileInfo, 0)
	for path, info := range fileInfoMap {
		lr := HTTPFileInfo{
			Name:    info.Name(),
			Path:    path,
			ModTime: info.ModTime().UTC(),
		}
		if info.IsDir() {
			name := info.Name()
			if deep {
				name = deepPath(localPath, info.Name())
			}
			lr.Name = name
			lr.Path = filepath.Join(filepath.Dir(path), name)
			lr.Type = "dir"
			lr.Size = s.historyDirSize(lr.Path)
		} else {
			lr.Type = "file"
			lr.Size = info.Size() // formatSize(info)
		}
		lrs = append(lrs, lr)
	}

	return
}

var dirSizeMap = make(map[string]int64)

func (s *fileSystem) makeIndex() error {
	var indexes = make([]IndexFileItem, 0)
	var err = filepath.Walk(s.Root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Errorf("WARN: Visit path: %s error: %v", strconv.Quote(path), err)
			return filepath.SkipDir
			// return err
		}
		if info.IsDir() {
			return nil
		}

		path, _ = filepath.Rel(s.Root, path)
		path = filepath.ToSlash(path)
		indexes = append(indexes, IndexFileItem{path, info})
		return nil
	})
	s.indexes = indexes
	dirSizeMap = make(map[string]int64)
	return err
}

func (s *fileSystem) historyDirSize(dir string) int64 {
	var size int64
	if size, ok := dirSizeMap[dir]; ok {
		return size
	}
	for _, fitem := range s.indexes {
		if filepath.HasPrefix(fitem.Path, dir) {
			size += fitem.Info.Size()
		}
	}
	dirSizeMap[dir] = size
	return size
}

func (s *fileSystem) findIndex(text string) []IndexFileItem {
	ret := make([]IndexFileItem, 0)
	for _, item := range s.indexes {
		ok := true
		// search algorithm, space for AND
		for _, keyword := range strings.Fields(text) {
			needContains := true
			if strings.HasPrefix(keyword, "-") {
				needContains = false
				keyword = keyword[1:]
			}
			if keyword == "" {
				continue
			}
			ok = (needContains == strings.Contains(strings.ToLower(item.Path), strings.ToLower(keyword)))
			if !ok {
				break
			}
		}
		if ok {
			ret = append(ret, item)
		}
	}
	return ret
}

func (s *fileSystem) IsFile(path string) bool {
	fullpath := filepath.Join(s.Root, path)
	info, err := os.Stat(fullpath)
	return err == nil && info.Mode().IsRegular()
}

func (s *fileSystem) IsDir(path string) bool {
	fullpath := filepath.Join(s.Root, path)
	info, err := os.Stat(fullpath)
	return err == nil && info.Mode().IsDir()
}

func deepPath(basedir, name string) string {
	isDir := true
	// loop max 5, incase of for loop not finished
	maxDepth := 5
	for depth := 0; depth <= maxDepth && isDir; depth += 1 {
		finfos, err := ioutil.ReadDir(filepath.Join(basedir, name))
		if err != nil || len(finfos) != 1 {
			break
		}
		if finfos[0].IsDir() {
			name = filepath.ToSlash(filepath.Join(name, finfos[0].Name()))
		} else {
			break
		}
	}
	return name
}

func checkFilename(name string) error {
	if strings.ContainsAny(name, "\\/:*<>|") {
		return errors.New("Name should not contains \\/:*<>|")
	}
	return nil
}
