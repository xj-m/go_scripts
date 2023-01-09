// package main
package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/xj-m/go_scripts/file"
	compress "github.com/xj-m/go_scripts/photo/compress"
)

func main() {
	todoDir := "todo_gif"
	dstDir := "generated_gif"

	compress.BulkCompress(todoDir, "tmp_todo_gif")

	file.MkdirIfNotExist(dstDir, dstDir)

	files, err := file.GetAllFilesWithExtension(todoDir, []string{".jpg", ".jpeg"})
	if err != nil {
		logrus.Fatal(err)
	}
	err = makeGif(files, dstDir)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("gif generated")
}

func makeGif(filePaths []string, dstDir string) error {
	logrus.Info(filePaths)
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

	logrus.Info(cmd.String())
	var stdErr strings.Builder
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error: %v, stderr: %s", err, stdErr.String())
	}

	compressGif(outputGifFileName, dstDir)

	return nil
}

func compressGif(fp string, dstDir string) error {
	file.MkdirIfNotExist(dstDir, fp)
	// compress fp file, then save to dstDir
	cmd := exec.Command(
		"ffmpeg",
		"-i",
		fp,
		"-vcodec",
		"libx264",
		"-crf",
		"36",
		"-preset",
		"medium",
		filepath.Join(dstDir, fp),
	)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
