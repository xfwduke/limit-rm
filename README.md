# limit-rm
speed controlled rm

go build ./

Usage:
  lrm [OPTIONS] files...

Application Options:
  -s, --speed= rm speed. support [KB, MB, GB] suffixes, or no suffix as BYTE (default: 10MB)
  -v           show progress detail

Help Options:
  -h, --help   Show this help message
