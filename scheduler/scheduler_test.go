package scheduler

import (
    "os/exec"
    "testing"
)

import (
    "github.com/coopersong/mini-spider/loader"
)

func TestNewScheduler(t *testing.T) {
    scheduler := NewScheduler()
    if scheduler == nil {
        t.Errorf("scheduler is nil")
        return
    }
}

func TestScheduler_Start(t *testing.T) {
    mkCmd := exec.Command("/bin/bash", "-c", "mkdir ../test_output")
    rmCmd := exec.Command("/bin/bash", "-c", "rm -rf ../test_output")

    err := mkCmd.Start()
    if err != nil {
        t.Errorf("mkCmd.Start(): %s", err.Error())
        return
    }
    err = mkCmd.Wait()
    if err != nil {
        t.Errorf("mkCmd.Wait(): %s", err.Error())
    }
    defer rmCmd.Start()

    scheduler := NewScheduler()
    cfg := loader.Config{
        loader.Spider{
            UrlListFile:     "../data/url.data",
            OutputDirectory: "../test_output",
            MaxDepth:        1,
            CrawlInterval:   1,
            CrawlTimeout:    1,
            TargetUrl:       ".*.(htm|html)$",
            ThreadCount:     8,
        },
    }
    seeds := []string{"http://www.baidu.com", "http://www.sina.com"}

    scheduler.Init(cfg, seeds)
    scheduler.Start()
}

func TestScheduler_RunTask(t *testing.T) {
    mkCmd := exec.Command("/bin/bash", "-c", "mkdir ../test_output")
    rmCmd := exec.Command("/bin/bash", "-c", "rm -rf ../test_output")
    err := mkCmd.Start()
    if err != nil {
        t.Errorf("mkCmd.Start(): %s", err.Error())
        return
    }
    err = mkCmd.Wait()
    if err != nil {
        t.Errorf("mkCmd.Wait(): %s", err.Error())
        return
    }
    defer rmCmd.Start()

    scheduler := NewScheduler()
    cfg := loader.Config{
        loader.Spider{
            UrlListFile:     "../data/url.data",
            OutputDirectory: "../test_output",
            MaxDepth:        2,
            CrawlInterval:   1,
            CrawlTimeout:    1,
            TargetUrl:       ".*.(htm|html)$",
            ThreadCount:     8,
        },
    }
    seeds := []string{"http://www.baidu.com", "http://www.sina.com"}
    scheduler.Init(cfg, seeds)
    task := &Task{
        Url:"http://www.test.com",
        Depth: 0,
        CommonCfg: scheduler.TaskCommonCfg,
    }
    scheduler.RunTask(task)
    task.Url = "http://www.baidu.com"
    scheduler.RunTask(task)
}
