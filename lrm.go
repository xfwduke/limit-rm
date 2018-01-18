package main

import (
	"time"
	"gopkg.in/cheggaaa/pb.v1"
	"github.com/jessevdk/go-flags"
	"os"
	"fmt"
	"regexp"
	"strings"
	"strconv"
)

const (
	actionInterval = 10.0  // rm 1-chunk per 10ms
	maxSpeed       = 100 * 1024 * 1024  // useless now
)

var speed float64

var opts struct {
	Speed  string `short:"s" long:"speed" default:"10MB" description:"rm speed. support [KB, MB, GB] suffixes, or no suffix as BYTE"`
	Detail bool   `short:"v" description:"show progress detail"`
	Args struct {
		FilePaths []string `positional-arg-name:"files"`
	} `positional-args:"yes" required:"yes"`
}

func init() {
	parser := flags.NewParser(&opts, flags.Default)
	_, err := parser.ParseArgs(os.Args[1:])
	if err != nil {
		os.Exit(1)
	}

	if len(opts.Args.FilePaths) <= 0 {
		fmt.Fprint(os.Stderr, "Error: need files to be rm")
		os.Exit(1)
	}

	pattern := regexp.MustCompile(`^(?i)([0-9]*\.?[0-9]+)(|kb|mb|gb)$`)
	res := pattern.FindStringSubmatch(opts.Speed)
	if res == nil {
		fmt.Fprint(os.Stderr, "Error: bad format of speed")
		os.Exit(1)
	}

	var speedTimes float64 = 1
	switch strings.ToLower(res[2]) {
	case "kb":
		speedTimes = 1024
	case "mb":
		speedTimes = 1024 * 1024
	case "gb":
		speedTimes = 1024 * 1024 * 1024
	default:
	}
	speed, err = strconv.ParseFloat(res[1], 64)
	if err != nil {
		panic(err)
	}
	speed *= speedTimes
}

func main() {
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

	for _, filePath := range opts.Args.FilePaths {
		fmt.Println(filePath)
		file, err := os.OpenFile(filePath, os.O_RDWR, 0777)
		if err != nil {
			panic(err)
		}

		fileStat, err := file.Stat()
		if err != nil {
			panic(err)
		}

		fileOrgSize := fileStat.Size()
		truncateSize := fileOrgSize

		bar := pb.New64(truncateSize)
		bar.SetRefreshRate(actionInterval * time.Millisecond)
		bar.SetMaxWidth(80)
		bar.SetUnits(pb.U_BYTES)
		bar.ShowSpeed = true
		bar.Start()

		for ; truncateSize > chunkSize; {
			<-token
			bar.Add64(chunkSize)
			truncateSize -= chunkSize
			err = file.Truncate(truncateSize)
			if err != nil {
				panic(err)
			}
		}
		bar.Add64(truncateSize)
		err = os.Remove(filePath)
		if err != nil {
			panic(err)
		}
		bar.Finish()
	}
}
