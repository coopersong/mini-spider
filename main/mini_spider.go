package main

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "time"
)

import (
    "github.com/baidu/go-lib/log"
    "github.com/baidu/go-lib/log/log4go"
)

import (
    "github.com/coopersong/mini-spider/loader"
    "github.com/coopersong/mini-spider/scheduler"
)

const (
    Version = "v1.0"
    SpiderConfFileName = "spider.conf"
)

var (
    confPath = flag.String("c", "../conf", "root path of configuration")
    logPath  = flag.String("l", "../log", "dir path of log")
    help     = flag.Bool("h", false, "to show help")
    version  = flag.Bool("v", false, "to show version")
)

func Exit(code int) {
    log.Logger.Close()
    time.Sleep(100 * time.Millisecond)
    os.Exit(code)
}

func initLog(logSwitch string, logPath *string, stdOut bool) error {
    log4go.SetLogBufferLength(10000)
    log4go.SetLogWithBlocking(false)

    err := log.Init("mini_spider", logSwitch, *logPath, stdOut, "midnight", 5)
    if err != nil {
        return fmt.Errorf("err in log.Init(): %s", err.Error())
    }

    return nil
}

func main() {
    var err error

    flag.Parse()

    if *help {
        flag.PrintDefaults()
        return
    }

    if *version {
        fmt.Println(Version)
        return
    }

    err = initLog("INFO", logPath, true)
    if err != nil {
        fmt.Printf("initLog(): %s\n", err.Error())
        Exit(-1)
    }

    config, err := loader.ConfigLoad(filepath.Join(*confPath, SpiderConfFileName))
    if err != nil {
        log.Logger.Error("loader.ConfigLoad(): %s", err.Error())
        Exit(-1)
    }

    seeds, err := loader.SeedLoad(config.UrlListFile)
    if err != nil {
        log.Logger.Error("loader.SeedLoad(): %s", err.Error())
        Exit(-1)
    }

    miniSpider := scheduler.NewScheduler()
    miniSpider.Init(config, seeds)
    miniSpider.Start()

    Exit(0)
}