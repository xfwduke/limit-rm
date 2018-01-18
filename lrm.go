package main

import (
	"github.com/jessevdk/go-flags"
	"os"
	"fmt"
	"regexp"
	"strings"
	"strconv"
	"./limitRM"
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
		if opts.Detail {
			fmt.Println(filePath)
		}
		err := limitRM.RM(filePath, speed, opts.Detail)
		if err != nil {
			panic(err)
		}
	}
}
