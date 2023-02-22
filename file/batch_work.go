package file

import (
	"context"
	"fmt"
	"sync"

	"github.com/xj-m/go_scripts/log"
)

func BatchWork(ctx context.Context, srcDir string, extNames []string, dstDir string, workFunc func(fp string, dstDir string) error, workerNum int) error {
	files, err := GetAllFilesWithExtension(srcDir, extNames)
	if err != nil {
		return err
	}

	log.GetLogger(nil).Info(fmt.Sprintf("[BatchWork] found %d files", len(files)))

	errCh := make(chan error, len(files))
	workCh := make(chan string, len(files))
	var wg sync.WaitGroup

	// spawn workerNum goroutines to process files
	for i := 0; i < workerNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for fp := range workCh {
				select {
				case <-ctx.Done():
					return
				default:
					errCh <- workFunc(fp, dstDir)
				}
			}
		}()
	}

	// enqueue files into workCh
	go func() {
		for _, fp := range files {
			workCh <- fp
		}
		close(workCh)
	}()

	// wait for all work to complete or an error occurs
	go func() {
		wg.Wait()
		close(errCh)
	}()

	var processed int
	for err := range errCh {
		if err != nil {
			return err
		}
		processed++
		log.GetLogger(nil).Info(fmt.Sprintf("[BatchWork] processed %d/%d files", processed, len(files)))
	}

	return nil
}
