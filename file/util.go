// package file
package file

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/djherbis/times"
)

// ErrorFileAlreadyExist is a error when file already exist
var ErrorFileAlreadyExist = fmt.Errorf("file already exist")

// GetFileSize return file size in human readable format, like 1.2M, 1.2G
func GetFileSize(fp string) string {
	fi, err := os.Stat(fp)
	if err != nil {
		return "unknown"
	}
	// in human readable format, like 1.2M
	return fmt.Sprintf("%.1f%s", float64(fi.Size())/float64(1024*1024), "M")
}

// MkdirParentIfNotExist create parent dir if not exist
func MkdirParentIfNotExist(dstDir string) error {
	// create parent dir for dstDir if not exist
	parentDir := filepath.Dir(dstDir)
	MkdirIfNotExist(parentDir)
	return nil
}

// MkdirIfNotExist create dir if not exist
func MkdirIfNotExist(dirName string) error {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err = os.MkdirAll(dirName, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetAllFilesWithExtension get all files that has target extension only in current dir
func GetAllFilesWithExtension(rootDir string, extNames []string) ([]string, error) {
	ret := []string{}
	err := filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
		// if in sub dir, skip
		if d.IsDir() && path != rootDir {
			return filepath.SkipDir
		}
		if err != nil {
			return err
		}
		// keep path if to lower is in extNames
		if !d.IsDir() && contains(extNames, strings.ToLower(filepath.Ext(path))) {
			ret = append(ret, path)
		}
		return nil
	})
	return ret, err
}

func contains[T comparable](elements []T, v T) bool {
	for _, s := range elements {
		if v == s {
			return true
		}
	}
	return false
}

func DeleteDir(dir string) error {
	cmd := exec.Command("rm", "-rf", dir)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func MatchFileTime(srcPath, dstPath string) error {
	// for dstPath, match the create time and modify  time and modify time to match with srcPath
	// get times
	t, err := times.Stat(srcPath)
	if err != nil {
		panic(err)
	}
	// print t
	fmt.Printf("AccessTime: %s, ModifyTime: %s, ChangeTime: %s\n", t.AccessTime(), t.ModTime(), t.ChangeTime())
	err = os.Chtimes(dstPath, t.AccessTime(), t.ModTime())
	if err != nil {
		return err
	}
	return nil
}

func GetAllSubDirs(dir string) ([]string, error) {
	// get all sub dirs under dir
	ret := []string{}
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && path != dir {
			ret = append(ret, path)
		}
		return nil
	})
	return ret, err
}
