// package file
package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

// MkdirIfNotExist create dir if not exist
func MkdirIfNotExist(dstDir string, fp string) error {
	if _, err := os.Stat(dstDir); os.IsNotExist(err) {
		err := os.Mkdir(dstDir, 0o755)
		if err != nil {
			return err
		}
	}
	// if file already exist, skip
	if _, err := os.Stat(filepath.Join(dstDir, fp)); !os.IsNotExist(err) {
		return ErrorFileAlreadyExist
	}
	return nil
}

// GetAllFilesWithExtension get all files that has target extension only in current dir
func GetAllFilesWithExtension(extNames []string) ([]string, error) {
	ret := []string{}
	err := filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		// if in sub dir, skip
		if d.IsDir() && path != "." {
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
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func contains[T comparable](elements []T, v T) bool {
	for _, s := range elements {
		if v == s {
			return true
		}
	}
	return false
}
