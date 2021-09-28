package scheduler

import (
    "fmt"
    "regexp"
    "strings"
)

import (
    "github.com/coopersong/mini-spider/crawler"
    "github.com/coopersong/mini-spider/parser"
    "github.com/coopersong/mini-spider/saver"
)

type TaskCommonConfig struct {
    // 爬取超时
    CrawlTimeout     int
    // 网页下载目录
    OutputDirectory  string
    // 需要存储的目标网页正则表达式
    TargetUrlPattern *regexp.Regexp
}

type Task struct {
    // 爬取url
    Url       string
    // 爬取深度
    Depth     int
    // 通用配置
    CommonCfg *TaskCommonConfig
}

// Create a new useless task.
func NewUselessTask(maxDepth int) *Task {
    return &Task{
        Url: "",
        Depth: maxDepth,
    }
}

// Run single task.
// A successful call returns sub url list and err == nil.
func (task *Task) Run() ([]string, error) {
    data, contentType, err := crawler.Crawl(task.Url, task.CommonCfg.CrawlTimeout)
    if err != nil {
        return nil, fmt.Errorf("%s: crawler.Crawl(): %s", task.Url, err.Error())
    }

    if !strings.Contains(contentType, "text") {
        return nil, fmt.Errorf("%s: Content-Type: %s", task.Url, contentType)
    }

    data, err = parser.Convert2Utf8(data, contentType)
    if err != nil {
        return nil, fmt.Errorf("%s: parser.Convert2Utf8(): %s", task.Url, err.Error())
    }

    if task.CommonCfg.TargetUrlPattern.MatchString(task.Url) {
        err = task.SaveData(data)
        if err != nil {
            return nil, fmt.Errorf("%s: task.SaveData(): %s", task.Url, err.Error())
        }
    }

    urlList, err := parser.GetUrlList(data, task.Url)
    if err != nil {
        return nil, fmt.Errorf("%s: parser.GetUrlList(): %s", task.Url, err.Error())
    }

    return urlList, nil
}

// Save data to output directory.
func (task *Task) SaveData(data []byte) error {
    err := saver.SaveData(data, task.Url, task.CommonCfg.OutputDirectory)
    if err != nil {
        return fmt.Errorf("saver.SaveData(): %s", err.Error())
    }

    return nil
}
