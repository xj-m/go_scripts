// package file
package file

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/djherbis/times"
)

// ErrorFileAlreadyExist is a error when file already exist
var (
	ErrorFileAlreadyExist = fmt.Errorf("file already exist")
	ErrorFileNotExist     = fmt.Errorf("file not exist")
)

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

func GetYearMonDay(t time.Time) (string, string, string) {
	// get year, mon, day from time
	// year: 2021, mon: 2021-01, day: 2021-01-01
	year := fmt.Sprintf("%d", t.Year())
	mon := fmt.Sprintf("%d-%02d", t.Year(), t.Month())
	day := fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
	return year, mon, day
}

func NewPathWithYearMon(path string, t time.Time) string {
	// get year, mon, day from time
	year, mon, _ := GetYearMonDay(t)
	return fmt.Sprintf("%s/%s/%s", year, mon, path)
}

func MoveFile(srcPath, dstPath string) error {
	// move file from srcPath to dstPath
	// if dstPath exist, return error
	if _, err := os.Stat(dstPath); !os.IsNotExist(err) {
		return ErrorFileAlreadyExist
	}
	// create parent dir for dstPath if not exist
	parentDir := filepath.Dir(dstPath)
	MkdirIfNotExist(parentDir)
	// move file
	err := os.Rename(srcPath, dstPath)
	if err != nil {
		return err
	}
	return nil
}

func ExtractTimeFromFileName(filePath string) (time.Time, error) {
	// get filename from path (without ext name) and extract time from filename
	filename := filepath.Base(filePath)
	// get filename without ext
	filenameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
	// get time from filename, using cur timezone
	t, err := time.ParseInLocation("2006-01-02", filenameWithoutExt, time.Local)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func IsFileExist(path string) bool {
	// check if file exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func CopyFileWithIO(srcPath, dstPath string) error {
	// cp file from srcPath to dstPath
	// if dstPath exist, return error
	if _, err := os.Stat(dstPath); !os.IsNotExist(err) {
		return ErrorFileAlreadyExist
	}
	// create parent dir for dstPath if not exist
	parentDir := filepath.Dir(dstPath)
	MkdirIfNotExist(parentDir)
	// os cp file, not link
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	return nil
}

func TruncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func ReadLines(filePath string) ([]string, error) {
	// read all lines from file
	// if file not exist, return error
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []string{}, ErrorFileNotExist
	}
	// read file
	file, err := os.Open(filePath)
	if err != nil {
		return []string{}, err
	}
	defer file.Close()
	// read lines
	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return []string{}, err
	}
	return lines, nil
}

func WriteLines(fileName string, lines []string) error {
	// write lines to file
	// if file exist, return error
	if _, err := os.Stat(fileName); !os.IsNotExist(err) {
		return ErrorFileAlreadyExist
	}
	// create parent dir for dstPath if not exist
	parentDir := filepath.Dir(fileName)
	MkdirIfNotExist(parentDir)
	// write lines to file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, line := range lines {
		fmt.Fprintln(file, line)
	}
	return nil
}

func OverWriteFile(fileName string, content string) error {
	// write content to file
	// if file not exist, return error
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return ErrorFileNotExist
	}
	// save file to tmp file with timestamp
	tmpFileName := "tmp/" + fileName + "." + strconv.FormatInt(time.Now().Unix(), 10)
	MkdirIfNotExist("tmp")
	err := ioutil.WriteFile(tmpFileName, []byte(content), 0o644)
	if err != nil {
		return err
	}
	// write content to file
	err = ioutil.WriteFile(fileName, []byte(content), 0o644)
	if err != nil {
		return err
	}
	return nil
}

func RunCmd(cmd string, args ...string) error {
	// run cmd with args
	// if cmd return error, return error
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	if err != nil {
		return err
	}
	return nil
}

func CreateSymlink(srcPath, dstPath string) error {
	// create symlink from srcPath to dstPath
	// if dstPath exist, return error
	if _, err := os.Stat(dstPath + ".tmp"); !os.IsNotExist(err) {
		return ErrorFileAlreadyExist
	}
	// create parent dir for dstPath if not exist
	parentDir := filepath.Dir(dstPath)
	MkdirIfNotExist(parentDir)
	// create symlink
	err := os.Symlink(srcPath, dstPath+".tmp")
	if err != nil {
		return err
	}
	// rename symlink
	err = os.Rename(dstPath+".tmp", dstPath)
	if err != nil {
		return err
	}
	return nil
}
