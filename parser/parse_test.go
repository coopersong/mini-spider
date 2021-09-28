package parser

import (
    "testing"
)

import (
    "github.com/coopersong/mini-spider/crawler"
)

func TestGetUrlList(t *testing.T) {
    testUrl := "https://www.baidu.com"

    data, _, err := crawler.Crawl(testUrl, 3)
    if err != nil {
        t.Errorf("%s: crawler.Crawl(): %s", testUrl, err.Error())
        return
    }

    urlList, err := GetUrlList(data, testUrl)
    if err != nil {
        t.Errorf("GetUrlList(): %s", err.Error())
        return
    }

    if len(urlList) == 0 {
        t.Errorf("no sublink in %s", testUrl)
        return
    }
}

func TestConvert2Utf8(t *testing.T) {
    testUrl := "http://vip.stock.finance.sina.com.cn/q/go.php/vDYData/kind/znzd/index.phtml"

    data, contentType, err := crawler.Crawl(testUrl, 3)
    if err != nil {
        t.Errorf("%s: crawler.Crawl(): %s", testUrl, err.Error())
        return
    }

    data, err = Convert2Utf8(data, contentType)
    if err != nil {
        t.Errorf("Convert2Utf8(): %s", err.Error())
        return
    }
}

func TestParseHostName(t *testing.T) {
    rawUrl := "http://xxx.baidu.com/v1/configs"
    hostName, err := ParseHostName(rawUrl)
    if err != nil {
        t.Errorf("%s: ParseHostName(): %s", rawUrl, err.Error())
        return
    }
    if hostName != "xxx.baidu.com" {
        t.Errorf("hostName: %s != xxx.baidu.com", hostName)
        return
    }

    rawUrl = "http://xxx.baidu.com:8080/v1/configs"
    hostName, err = ParseHostName(rawUrl)
    if err != nil {
        t.Errorf("%s: ParseHostName(): %s", rawUrl, err.Error())
        return
    }
    if hostName != "xxx.baidu.com" {
        t.Errorf("hostName: %s != xxx.baidu.com", hostName)
        return
    }

    rawUrl = "http:++xxx.baidu.com:8080/v1/configs/page"
    hostName, err = ParseHostName(rawUrl)
    if err == nil {
        t.Errorf("%s: there should be an error but not, hostName: %s", rawUrl, hostName)
        return
    }

    rawUrl = "http://xxx.baidu.com:8080/v1\n/configs/page"
    hostName, err = ParseHostName(rawUrl)
    if err == nil {
        t.Errorf("%s: there should be an error but not, hostName: %s", rawUrl, hostName)
        return
    }
}