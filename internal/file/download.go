package file

import (
	"fmt"
	"os"
	"time"
)

func downloadFile(quit <-chan bool, f *File) {
	var stop bool
	for {
		select {
		case <-quit:
			stop = true
		default:
			file, err := os.Open(f.Path())
			if err != nil {
				return
			}

			stat, err := file.Stat()
			if err != nil {
				return
			}
			size := stat.Size()
			if size == 0 {
				size = 1
			}
			var percent float64 = float64(size) / float64(f.Size) * 100
			f.StatusValue = fmt.Sprintf("%0.f%%", percent)
			if stop {
				break
			}
			time.Sleep(time.Second)
		}
	}
}
