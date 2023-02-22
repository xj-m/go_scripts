// package main
package main

import (
	"context"
	"os/exec"

	"github.com/xj-m/go_scripts/file"
)

func main() {
	srcDir := "/Volumes/Untitled/DCIM/Camera01"
	dstDir := "/Volumes/xiangjun-1T/mxj_media"

	// move all "insv" and "lrv" files from srcPath to dstPath
	file.BatchWork(
		context.Background(),
		srcDir,
		[]string{".insv", ".lrv", ".dng", ".insp"},
		dstDir,
		func(fp string, dstDir string) error {
			// move fp to dstDir
			cmd := exec.Command("mv", fp, dstDir)
			return cmd.Run()
		},
		5,
	)
}
