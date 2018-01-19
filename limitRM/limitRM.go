package limitRM

import (
	"time"
	"os"
	"context"
)

const (
	actionInterval = 10.0              // rm 1-chunk per 10ms
	maxSpeed       = 100 * 1024 * 1024 // useless now
)

func RM(ctx context.Context, filePath string, speed float64, progress chan int64) error {
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

	for ; fileStat.Size() > chunkSize; {
		<-token

		fileStat, err = file.Stat()
		if err != nil {
			return err
		}

		if fileStat.Size()-chunkSize <= 0 {
			break
		}

		err = file.Truncate(fileStat.Size() - chunkSize)
		if err != nil {
			return err
		}

		progress <- chunkSize
	}

	fileStat, err = file.Stat()
	if err != nil {
		return err
	}

	progress <- fileStat.Size()
	err = os.Remove(filePath)
	if err != nil {
		return err
	}

	close(progress)
	select {
	case <-ctx.Done():
	}
	return nil
}
