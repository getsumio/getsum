package file

import (
	"fmt"
	"os"
	"time"

	"github.com/getsumio/getsum/internal/status"
)

//watches download process
//and update status somethig like: DOWNLOAD 12%
//quit channel immediatly terminates it
func downloadFile(quit <-chan bool, f *File) {
	f.Status.Type = status.DOWNLOAD
	for {
		select {
		case <-quit:
			return
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
			f.Status.Value = fmt.Sprintf("%0.f%%", percent)
			time.Sleep(time.Second)
		}
	}
}
