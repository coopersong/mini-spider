package crawler

import (
    "testing"
)

func TestCrawl(t *testing.T) {
    // 测试正常的情况
    _, _, err := Crawl("https://www.baidu.com", 3)
    if err != nil {
        t.Errorf("fail to crawl www.baidu.com")
        return
    }

    // 测试url不合法的情况
    _, _, err = Crawl("xxx[]009", 3)
    if err == nil {
        t.Errorf("crawl invalid url should cause error but not")
        return
    }

    // 测试无法访问url的情况
    _, _, err = Crawl("https://www.guangze.com", 3)
    if err == nil {
        t.Errorf("no website named https://www.guangze.com but not cause error")
        return
    }
}
