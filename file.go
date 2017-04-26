package commons

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

var (
	ConsFile     *File
	fileCacheDir string
)

type File struct{}

func init() {
	ConsFile = &File{}
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
		err = errors.New(PathEmpty)
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
// @Returns file name:string
func (f *File) GetFileFullNameByPath(path string) string {
	var (
		index int
	)
	if strings.Contains(path, "/") {
		index = strings.LastIndexAny(path, "/")
		path = path[index+1 : len(path)]
	}
	goto RETURN
RETURN:
	return path
}

// @Title GetFileExtension
// @Description get file extension from path
// @Parameters
//            path            string        file path
// @Returns file extension:string
func (f *File) GetFileExtension(path string) string {
	var (
		index int
	)
	if strings.Contains(path, "/") {
		path = f.GetFileFullNameByPath(path)
	}
	if strings.Contains(path, ".") {
		index = strings.LastIndex(path, ".")
		path = strings.ToLower(path[index+1 : len(path)])
	}
	path = ""
	goto RETURN
RETURN:
	return path
}

// @Title GetFileNameByPath
// @Description get file name from path
// @Parameters
//            filename      string                    file name
// @Returns file name:string
func (f *File) GetFileNameByPath(path string) string {
	var (
		index int
	)
	if strings.Contains(path, "/") {
		path = f.GetFileFullNameByPath(path)
	}
	if strings.Contains(path, ".") {
		index = strings.LastIndex(path, ".")
		path = path[:index]
	}
	goto RETURN
RETURN:
	return path
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
	err = os.RemoveAll(path)
	goto RETURN
RETURN:
	return err
}
