package saver

import (
    "os/exec"
    "testing"
)

func TestSaveData(t *testing.T) {
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

    data := []byte("Hello World!")
    err = SaveData(data, "www.test.com", "../test_output")
    if err != nil {
        t.Errorf("SaveData(): %s", err.Error())
        return
    }
}