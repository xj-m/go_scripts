// package video
package video

import (
	"os/exec"
	"path/filepath"

	"github.com/xj-m/go_scripts/file"
)

func CompressMp4File(fp string, dstDir string) error {
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

