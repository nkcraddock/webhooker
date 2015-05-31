package testhelp

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

type TestData map[string]map[string]interface{}

func LoadTestData(file string) (TestData, error) {
	data, err := ioutil.ReadFile(path.Join("../", file))
	if err != nil {
		return nil, err
	}

	td := make(TestData)

	if err := json.Unmarshal(data, &td); err != nil {
		return nil, err
	}

	return td, nil
}
