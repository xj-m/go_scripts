// package main
package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/xj-m/go_scripts/file"
)

func main() {
	todoDir := "todo_gif"
	dstDir := "generated_gif"

	file.MkdirIfNotExist(dstDir, dstDir)

	files, err := file.GetAllFilesWithExtension(todoDir, []string{".jpg", ".jpeg"})
	if err != nil {
		panic(err)
	}
	err = makeGif(files, dstDir)
	if err != nil {
		panic(err)
	}
	logrus.Info("gif generated")
}

func makeGif(filePaths []string, dstDir string) error {
	logrus.Info(filePaths)
	//  create the animation.gif using filePaths, with 0.1 second delay between frames
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
	cmd.Args = append(cmd.Args, filepath.Join(dstDir, "animation.gif"))
	logrus.Info(cmd.String())
	var stdErr strings.Builder
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error: %v, stderr: %s", err, stdErr.String())
	}
	return nil
}
