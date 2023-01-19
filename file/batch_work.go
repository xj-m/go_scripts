package file

import (
	"path/filepath"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

const COMPRESSED_DIR_NAME = "compressed"

func BatchWork(rootDir string, extNames []string, workFunc func(fp string, dstDir string) error) {
	files, err := GetAllFilesWithExtension(rootDir, extNames)
	if err != nil {
		logrus.Fatal(err)
		return
	}
	logrus.Infof("files total %v: \n\t%v", len(files), strings.Join(files, "\n\t"))
	wg := sync.WaitGroup{}
	for i, fp := range files {
		wg.Add(1)
		go func(fp string, i int) {
			defer wg.Done()
			logrus.Infof("[compressing](%v/%v): %v", i+1, len(files), fp)
			err := workFunc(fp, COMPRESSED_DIR_NAME)
			switch err {
			case ErrFileAlreadyExist:
				logrus.Infof("(%v/%v) file \"%v\" already exist, skip", i+1, len(files), fp)
			case nil:
				logrus.Infof("(%v/%v) compress file \"%v\" success, size before: %v, after compress: %v", i+1, len(files), fp, GetFileSize(fp), GetFileSize(filepath.Join("compressed", fp)))
			default:
				logrus.Errorf("(%v/%v) compress file \"%v\" failed: %v", i+1, len(files), fp, err)
			}
		}(fp, i)
	}
	wg.Wait()
}
