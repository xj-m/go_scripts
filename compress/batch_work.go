package compress

import (
	"path/filepath"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/xj-m/go_scripts/file"
)

const COMPRESSED_DIR_NAME = "compressed"

func BatchWork(rootDir string, extNames []string, workFunc func(fp string, dstDir string) error) {
	files, err := file.GetAllFilesWithExtension(rootDir, extNames)
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
			case file.ErrorFileAlreadyExist:
				logrus.Infof("(%v/%v) file \"%v\" already exist, skip", i+1, len(files), fp)
			case nil:
				logrus.Infof("(%v/%v) compress file \"%v\" success, size before: %v, after compress: %v", i+1, len(files), fp, file.GetFileSize(fp), file.GetFileSize(filepath.Join("compressed", fp)))
			default:
				logrus.Errorf("(%v/%v) compress file \"%v\" failed: %v", i+1, len(files), fp, err)
			}
		}(fp, i)
	}
	wg.Wait()
}
