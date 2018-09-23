package g

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"
)

var StartTime int64

func Init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	RunPid()
	MyLog()
	InitRootDir()
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	ParseConfig(*cfg)
	StartTime = time.Now().Unix()
}
