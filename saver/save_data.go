package saver

import (
    "fmt"
    "net/url"
    "os"
    "path/filepath"
)

// Create a file named urlStr in outputDirectory and save data to it.
func SaveData(data []byte, urlStr string, outputDirectory string) error {
    fileName := filepath.Join(outputDirectory, url.QueryEscape(urlStr))

    f, err := os.OpenFile(fileName, os.O_CREATE | os.O_WRONLY, 0666)
    if err != nil {
        return fmt.Errorf("%s: os.OpenFile(): %s", fileName, err.Error())
    }
    defer f.Close()

    _, err = f.Write(data)
    if err != nil {
        return fmt.Errorf("f.Write(): %s", err.Error())
    }

    return nil
}