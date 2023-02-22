// package main
package main

import (
	"context"

	"github.com/xj-m/go_scripts/photo"
)

func main() {
	srcDir := "."
	copyToDirName := "compressed_jpg"
	photo.BulkCompressPhoto(context.Background(), srcDir, copyToDirName)
}
