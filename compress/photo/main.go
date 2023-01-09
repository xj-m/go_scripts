// package main
package main

import (
	"os/exec"

	"github.com/xj-m/go_scripts/compress"
	"github.com/xj-m/go_scripts/file"
)

func main() {
	compress.BatchWork([]string{".mp4"}, compressJPGFile)
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
