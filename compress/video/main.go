// package main
package main

import (
	"os/exec"
	"path/filepath"

	"github.com/xj-m/go_scripts/compress"
	"github.com/xj-m/go_scripts/file"
)

func main() {
	compress.BatchWork([]string{".mp4"}, compressMp4File)
}

func compressMp4File(fp string, dstDir string) error {
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
