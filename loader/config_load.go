package loader

import (
    "fmt"
    "regexp"
)

import (
    "gopkg.in/gcfg.v1"
)

type Spider struct {
    // 种子文件路径
    UrlListFile     string
    // 抓取结果存储目录
    OutputDirectory string
    // 最大抓取深度
    MaxDepth        int
    // 抓取间隔. 单位: 秒
    CrawlInterval   int
    // 抓取超时. 单位: 秒
    CrawlTimeout    int
    // 需要存储的目标网页URL Pattern
    TargetUrl       string
    // 抓取routine数
    ThreadCount     int
}

type Config struct {
    Spider
}

// Load config from confPath.
func ConfigLoad(confPath string) (Config, error) {
    var cfg Config

    if err := gcfg.ReadFileInto(&cfg, confPath); err != nil {
        return cfg, err
    }

    if err := cfg.Check(); err != nil {
        return cfg, err
    }

    return cfg, nil
}

// Check config.
func (c *Config) Check() error {
    if c.UrlListFile == "" {
        return fmt.Errorf("UrlListFile is nil")
    }

    if c.OutputDirectory == "" {
        return fmt.Errorf("OutputDirectory is nil")
    }

    if c.MaxDepth < 1 {
        return fmt.Errorf("MaxDepth is less than 1")
    }

    if c.CrawlInterval < 0 {
        return fmt.Errorf("CrawlInterval is less than 0")
    }

    if c.CrawlTimeout < 1 {
        return fmt.Errorf("CrawlTimeout is less than 1")
    }

    _, err := regexp.Compile(c.TargetUrl)
    if err != nil {
        return fmt.Errorf("%s: regexp.Compile(): %s", c.TargetUrl, err.Error())
    }

    if c.ThreadCount < 1 {
        return fmt.Errorf("ThreadCount is less than 1")
    }

    return nil
}
