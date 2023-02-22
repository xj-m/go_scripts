// package photo

package photo

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/xj-m/go_scripts/file"
	"github.com/xj-m/go_scripts/log"
)

func BulkCompressPhoto(ctx context.Context, srcDir, copyToDirName string) {
	extNames := []string{".jpeg", ".jpg"}

	file.MkdirIfNotExist(copyToDirName)

	// copy all files to copyToDirName
	files, err := file.GetAllFilesWithExtension(srcDir, extNames)
	if err != nil {
		panic(err)
	}
	for _, fp := range files {
		// get fp filename
		_, filename := filepath.Split(fp)
		dst := filepath.Join(copyToDirName, filename)

		file.MkdirIfNotExist(copyToDirName)
		cmd := exec.Command("cp", fp, dst)

		// log cmd string

		err := cmd.Run()
		if err != nil {
			log.GetLogger(nil).Error(fmt.Sprintf("cmd failed: %s", strings.Join(cmd.Args, " ")))
			panic(err)
		}
	}

	// compress copyToDirName
	file.BatchWork(ctx, copyToDirName, extNames, "compressed", compressJPGFile)
}

func compressJPGFile(fp string, dstDir string) error {
	file.MkdirIfNotExist(dstDir)
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

func MakeGif(filePaths []string, dstDir string) (filePath string, err error) {
	//  create the animation.gif using filePaths, with 30ms delay between frames
	cmd := exec.Command(
		"convert",
		"-delay",
		"30",
		"-loop",
		"0",
	)
	for _, fp := range filePaths {
		cmd.Args = append(cmd.Args, fp)
	}
	// use first element to name the gif
	outputGifFileName := filepath.Join(dstDir, fmt.Sprintf("%s.gif", filepath.Base(filePaths[0])))
	cmd.Args = append(cmd.Args, outputGifFileName)

	log.GetLogger(nil).Info(cmd.String())
	var stdErr strings.Builder
	cmd.Stderr = &stdErr

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error: %v, stderr: %s", err, stdErr.String())
	}

	err = compressGif(outputGifFileName, dstDir)

	return outputGifFileName, err
}

func compressGif(fp string, dstDir string) error {
	file.MkdirIfNotExist(dstDir)
	// resize the gif file to 50%
	cmd := exec.Command(
		"convert",
		fp,
		"-resize",
		"50%",
		fp,
	)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
