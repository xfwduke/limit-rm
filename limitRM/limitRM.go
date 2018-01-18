package limitRM

import (
	"time"
	"gopkg.in/cheggaaa/pb.v1"
	"os"
)

const (
	actionInterval = 10.0              // rm 1-chunk per 10ms
	maxSpeed       = 100 * 1024 * 1024 // useless now
)


func RM(filePath string, speed float64, detail bool) error {
	chunkSize := int64(speed * actionInterval / 1000)
	if chunkSize < 1 {
		chunkSize = 1
	}

	token := make(chan int)
	go func(chan int) {
		ticker := time.NewTicker(actionInterval * time.Millisecond)
		for range ticker.C {
			select {
			case token <- 0:
			default:
			}
		}
	}(token)

	file, err := os.OpenFile(filePath, os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}

	fileStat, err := file.Stat()
	if err != nil {
		return err
	}

	fileOrgSize := fileStat.Size()
	truncateSize := fileOrgSize

	bar := pb.New64(truncateSize)
	bar.SetRefreshRate(actionInterval * time.Millisecond)
	bar.SetMaxWidth(80)
	bar.SetUnits(pb.U_BYTES)
	bar.ShowSpeed = true
	if detail {
		bar.Start()
	}

	for ; truncateSize > chunkSize; {
		<-token

		if detail {
			bar.Add64(chunkSize)
		}
		truncateSize -= chunkSize
		err = file.Truncate(truncateSize)
		if err != nil {
			return err
		}
	}
	if detail {
		bar.Add64(truncateSize)
	}
	err = os.Remove(filePath)
	if err != nil {
		return err
	}
	if detail {
		bar.Finish()
	}
	return nil
}
