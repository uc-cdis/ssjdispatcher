package ssjdispatcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

func ReadFile(path string) ([]byte, error) {
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return buff, nil
}

func GetKeyValueFromConfigFile(cfgFile string, keys []string) (interface{}, error) {
	buff, err := ReadFile(cfgFile)
	if err != nil {
		return nil, err
	}

	return GetValueFromKeys(buff, keys)

}

func GetValueFromKeys(buff []byte, keys []string) (interface{}, error) {
	if len(keys) == 0 {
		return nil, errors.New("KeyValue")
	}

	var m map[string]interface{}
	json.Unmarshal(buff, &m)

	result := m[keys[0]]
	err := false

	for _, key := range keys[1:] {
		result, err = result.(map[string]interface{})[key]
		if err == false {
			return nil, errors.New("KeyValue")
		}
	}

	return result, nil

}
