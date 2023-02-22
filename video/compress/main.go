// package main
package main

import (
	"context"

	"github.com/xj-m/go_scripts/file"
	"github.com/xj-m/go_scripts/video"
)

func main() {
	file.BatchWork(context.Background(), ".", []string{".mp4"}, "compressed", video.CompressMp4File, 5)
}
