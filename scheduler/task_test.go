package scheduler

import (
    "os/exec"
    "regexp"
    "testing"
)

func TestNewUselessTask(t *testing.T) {
    uselessTask := NewUselessTask(3)
    if uselessTask == nil {
        t.Errorf("uselessTask should not be nil but is nil")
        return
    }
}

func TestTask_Run(t *testing.T) {
    mkCmd := exec.Command("/bin/bash", "-c", "mkdir ../test_output")
    rmCmd := exec.Command("/bin/bash", "-c", "rm -rf ../test_output")

    err := mkCmd.Start()
    if err != nil {
        t.Errorf("mkCmd.Start(): %s", err.Error())
        return
    }
    defer rmCmd.Start()

    targetUrlPattern, _ := regexp.Compile(".*.(htm|html)$")
    task := &Task{
        Url: "http://www.test.com",
        Depth: 1,
        CommonCfg: &TaskCommonConfig{
            CrawlTimeout: 2,
            OutputDirectory: "../test_output",
            TargetUrlPattern: targetUrlPattern,
        },
    }

    urlList, err := task.Run()
    if err == nil || len(urlList) > 0 {
        t.Errorf("task.Run() should return error but not")
        return
    }

    task.Url = "http://www.baidu.com"
    _, err = task.Run()
    if err != nil {
        t.Errorf("task.Run(): %s", err.Error())
        return
    }
}

func TestTask_SaveData(t *testing.T) {
    mkCmd := exec.Command("/bin/bash", "-c", "mkdir ../test_output")
    rmCmd := exec.Command("/bin/bash", "-c", "rm -rf ../test_output")

    err := mkCmd.Start()
    if err != nil {
        t.Errorf("mkCmd.Start(): %s", err.Error())
        return
    }
    defer rmCmd.Start()

    task := &Task{
        Url: "http://www.test.com",
        Depth: 1,
        CommonCfg: &TaskCommonConfig{
            OutputDirectory: "../test_output",
        },
    }
    data := []byte("Hello World!")
    err = task.SaveData(data)
    if err != nil {
        t.Errorf("task.SaveData(): %s", err.Error())
        return
    }

    task.CommonCfg.OutputDirectory = "../xxx"
    err = task.SaveData(data)
    if err == nil {
        t.Errorf("task.SaveData() should return error but not")
        return
    }
}
