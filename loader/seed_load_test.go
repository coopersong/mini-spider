package loader

import (
    "testing"
)

func TestSeedLoad(t *testing.T) {
    // 测试文件路径不存在的情况
    _, err := SeedLoad("./not_exist.data")
    if err == nil {
        t.Errorf("./not_exist.data is not exist but there is no error")
        return
    }
    // 测试json解析有问题的情况
    _, err = SeedLoad("./test_url1.data")
    if err == nil {
        t.Errorf("./test_url1.data's json format is invalid but there is no error")
        return
    }
    // 测试种子文件没有种子的情况
    _, err = SeedLoad("./test_url2.data")
    if err == nil {
        t.Errorf("there is no seed in ./test_url2.data but there is no error")
        return
    }
    // 测试正常情况
    _, err = SeedLoad("../data/url.data")
    if err != nil {
        t.Errorf("../data/url.data is valid but there is an error")
        return
    }
}
