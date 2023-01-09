// package main
package main

import (
	"github.com/xj-m/go_scripts/file"
	"github.com/xj-m/go_scripts/video"
)

func main() {
	file.BatchWork(".", []string{".mp4"}, video.CompressMp4File)
}
