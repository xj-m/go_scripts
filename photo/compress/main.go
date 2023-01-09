// package main
package main

import (
	"github.com/xj-m/go_scripts/photo"
)

func main() {
	srcDir := "."
	copyToDirName := "compressed_jpg"
	photo.BulkCompressPhoto(srcDir, copyToDirName)
}
