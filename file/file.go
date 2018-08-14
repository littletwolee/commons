package file

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var (
	consFile     *File
	fileCacheDir string
)

const (
	ERROR_PATH_EMPTY = "Your path is empty?"
)

type File struct{}

func GetFile() *File {
	if consFile == nil {
		consFile = &File{}
	}
	return consFile
	// fileCacheDir = ConsConfigHelper.GetValue("filecache", "path")
	// ConsFileHelper.PathExists("", true)
}

// @Title PathExists
// @Description check directory/file is exists
// @Parameters
//            path              string          file/directory path
//            iscreate          bool            need to create directory?
// @Returns err:error
func (f *File) PathExists(path string, iscreate bool) error {
	var (
		err error
	)
	if path == "" {
		err = errors.New(ERROR_PATH_EMPTY)
		goto RETURN
	}
	_, err = os.Stat(path)
	if err == nil {
		goto RETURN
	}
	if os.IsNotExist(err) {
		if iscreate {
			err = os.Mkdir(path, os.ModePerm)
			if err != nil {
				goto RETURN
			}
			goto RETURN
		}
		goto RETURN
	}
	goto RETURN
RETURN:
	return err
}

// @Title GetFilesInfo
// @Description get all files info from path
// @Parameters
//            path              string       path
// @Returns list:[]os.FileInfo err:error
func (f *File) GetFilesInfo(path string) ([]os.FileInfo, error) {
	var (
		err     error
		dirList []os.FileInfo
	)
	if path == "" {
		err = errors.New(ERROR_PATH_EMPTY)
		goto RETURN
	}
	dirList, err = ioutil.ReadDir(path)
	if err != nil {
		goto RETURN
	}
	goto RETURN
RETURN:
	return dirList, err
}

// @Title GetFileFullNameByPath
// @Description get file full name from path
// @Parameters
//            path           string       path
// @Returns file name:string err:error
func (f *File) GetFileFullNameByPath(path string) (string, error) {
	var (
		index int
		err   error
	)
	if path == "" {
		err = errors.New(ERROR_PATH_EMPTY)
		goto RETURN
	}
	if strings.Contains(path, "/") {
		index = strings.LastIndexAny(path, "/")
		path = path[index+1 : len(path)]
	}
	goto RETURN
RETURN:
	return path, err
}

// @Title GetFileExtension
// @Description get file extension from path
// @Parameters
//            path            string        file path
// @Returns file extension:string err:error
func (f *File) GetFileExtension(path string) (string, error) {
	var (
		index int
		err   error
	)
	if path == "" {
		err = errors.New(ERROR_PATH_EMPTY)
		goto RETURN
	}
	if strings.Contains(path, "/") {
		path, err = f.GetFileFullNameByPath(path)
		if err != nil {
			goto RETURN
		}
	}
	if strings.Contains(path, ".") {
		index = strings.LastIndex(path, ".")
		path = strings.ToLower(path[index+1 : len(path)])
	}
	path = ""
	goto RETURN
RETURN:
	return path, err
}

// @Title GetFileNameByPath
// @Description get file name from path
// @Parameters
//            filename      string                    file name
// @Returns file name:string err:error
func (f *File) GetFileNameByPath(path string) (string, error) {
	var (
		index int
		err   error
	)
	if path == "" {
		err = errors.New(ERROR_PATH_EMPTY)
		goto RETURN
	}
	if strings.Contains(path, "/") {
		path, err = f.GetFileFullNameByPath(path)
		if err != nil {
			goto RETURN
		}
	}
	if strings.Contains(path, ".") {
		index = strings.LastIndex(path, ".")
		path = path[:index]
	}
	goto RETURN
RETURN:
	return path, err
}

// @Title PathDel
// @Description delete directory/file from disk
// @Parameters
//            path              string       path
// @Returns err:error
func (f *File) PathDel(path string) error {
	var (
		err error
	)
	if path == "" {
		err = errors.New(ERROR_PATH_EMPTY)
		goto RETURN
	}
	err = os.RemoveAll(path)
	goto RETURN
RETURN:
	return err
}

// @Title OpenFile
// @Description open file in io.Writer
// @Parameters
//            path            string           path
//            flags           int              flags
// @Returns writer:*os.File err:error
func (f *File) OpenFile(path string, flags int) (*os.File, error) {
	var (
		err  error
		file *os.File
	)
	if path == "" {
		err = errors.New(ERROR_PATH_EMPTY)
		goto RETURN
	}
	file, err = os.OpenFile(path, flags, 0666)
	goto RETURN
RETURN:
	return file, err
}

// @Title FormatPath
// @Description format path from xx/xx to xx/xx/
// @Parameters
//            path              string       path
// @Returns path:string err:error
func (f *File) FormatPath(path string) (string, error) {
	var (
		err         error
		currentPath string
	)
	if path == "" {
		err = errors.New(ERROR_PATH_EMPTY)
		goto RETURN
	}
	if path[len(path)-1:] != "/" {
		path = fmt.Sprintf("%s%s", path, "/")
	}
	currentPath, err = f.GetCurrentDirectory()
	if err != nil {
		goto RETURN
	}
	path = fmt.Sprintf("%s%s", currentPath, path)
	goto RETURN
RETURN:
	return path, err
}

// @Title FormatPath
// @Description format path from xx/xx to xx/xx/
// @Returns path:string err:error
func (f *File) GetCurrentDirectory() (string, error) {
	var (
		err         error
		currentPath string
		index       int
	)
	currentPath, err = exec.LookPath(os.Args[0])
	if err != nil {
		goto RETURN
	}
	index = strings.LastIndex(currentPath, "\\")
	currentPath = string(currentPath[0 : index+1])
	goto RETURN
RETURN:
	return currentPath, err
}

// @Title ReadFile
// @Description read file from path
// @Parameters
//            path              string       path
// @Returns data:[]byte err:error
func (f *File) ReadFile(path string) ([]byte, error) {
	var (
		data []byte
		err  error
	)
	data, err = ioutil.ReadFile(path)
	goto RETURN
RETURN:
	return data, err
}
