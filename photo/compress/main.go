// package main
package main

import (
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/xj-m/go_scripts/file"
)

func main() {
	extNames := []string{".jpeg", ".jpg"}

	// copy all files to copyToDirName
	copyToDirName := "compressed_jpg"
	files, err := file.GetAllFilesWithExtension(".", extNames)
	if err != nil {
		logrus.Fatal(err)
		return
	}
	for _, fp := range files {
		dst := filepath.Join(copyToDirName, fp)
		file.MkdirIfNotExist(copyToDirName, fp)
		cmd := exec.Command("cp", fp, dst)
		err := cmd.Run()
		if err != nil {
			logrus.Fatal(err)
		}
	}

	// compress copyToDirName
	file.BatchWork(copyToDirName, extNames, compressJPGFile)
}

func compressJPGFile(fp string, dstDir string) error {
	file.MkdirIfNotExist(dstDir, fp)
	// compress fp file, then save to dstDir
	cmd := exec.Command(
		"jpegoptim",
		"--strip-all",
		"--all-progressive",
		"--max=80",
		fp,
	)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
