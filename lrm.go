package main

import (
	"github.com/jessevdk/go-flags"
	"os"
	"fmt"
	"regexp"
	"strings"
	"strconv"
	"./limitRM"
	"gopkg.in/cheggaaa/pb.v1"
	"time"
	"context"
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
	for _, filePath := range opts.Args.FilePaths {
		progress := make(chan int64/*, 100*/)

		file, err := os.OpenFile(filePath, os.O_RDWR, 0777)
		if err != nil {
			panic(err)
		}
		fileStat, err := file.Stat()
		if err != nil {
			panic(err)
		}

		ctx, cancel := context.WithCancel(context.Background())

		bar := pb.New64(fileStat.Size())
		bar.SetRefreshRate(10 * time.Millisecond)
		bar.SetMaxWidth(80)
		bar.SetUnits(pb.U_BYTES)
		bar.ShowSpeed = true

		file.Close()

		if opts.Detail {
			go func(string, *pb.ProgressBar, chan int64) {
				fmt.Println(filePath)
				bar.Start()
				defer func() {
					bar.Finish()
					cancel()
				}()
				for p := range progress {
					bar.Add64(p)
				}
			}(filePath, bar, progress)
		}

		err = limitRM.RM(ctx, filePath, speed, progress)
		if err != nil {
			panic(err)
		}
	}
}
