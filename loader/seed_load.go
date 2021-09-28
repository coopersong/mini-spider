package loader

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
)

// Load seeds from path.
func SeedLoad(path string) ([]string, error) {
    var seeds []string
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }
    err = json.Unmarshal(data, &seeds)
    if err != nil {
        return nil, fmt.Errorf("json.Unmarshal(): %s", err.Error())
    }
    if len(seeds) == 0 {
        return nil, fmt.Errorf("no seed in %s", path)
    }
    return seeds, nil
}