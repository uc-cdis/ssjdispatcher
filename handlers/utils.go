package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyz"

// GetRandString returns a random string of lenght N
func GetRandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func ReadFile(path string) ([]byte, error) {
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Can not read file %s. Detail %s", path, err)
	}
	return buff, nil
}

func containKey(mapping map[string]interface{}, key string) bool {
	for k := range mapping {
		if k == key {
			return true
		}
	}
	return false
}

func GetValueFromJSON(jsonBytes []byte, keys []string) (interface{}, error) {
	var dataMap interface{}
	err := json.Unmarshal(jsonBytes, &dataMap)
	if err != nil {
		return nil, err
	}
	if dataMap == nil {
		return nil, fmt.Errorf("Can not unmarshal data bytes")
	}
	for _, key := range keys {
		if containKey(dataMap.(map[string]interface{}), key) {
			dataMap = dataMap.(map[string]interface{})[key]
		} else {
			return nil, errors.New(string(jsonBytes) + " does not contain key " + key)
		}
	}
	return dataMap, nil
}

func StringContainsPrefixInSlice(s string, prefixList []string) bool {
	for _, prefix := range prefixList {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}
