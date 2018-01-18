# limit-rm
speed controlled rm

```
go build ./

Usage:
  lrm [OPTIONS] files...

Application Options:
  -s, --speed= rm speed. support [KB, MB, GB] suffixes, or no suffix as BYTE (default: 10MB)
  -v           show progress detail

Help Options:
  -h, --help   Show this help message
```

```
lrm -s 70MB -v CLion-2017.3.tar.gz pycharm-professional-2017.3.1.tar.gz django-rest-framework-cn.pdf                              
CLion-2017.3.tar.gz
 315.68 MiB / 315.68 MiB [==============================] 100.00% 70.14 MiB/s 4s
pycharm-professional-2017.3.1.tar.gz
 339.49 MiB / 339.49 MiB [==============================] 100.00% 70.15 MiB/s 4s
django-rest-framework-cn.pdf
 652.26 KiB / 652.26 KiB [===============================] 100.00% 4.12 GiB/s 0s
```
