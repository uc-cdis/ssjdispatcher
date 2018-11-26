package ssjdispatcher

import (
	"errors"
)

type ImagePatternMap struct {
	Mapping map[string]string
}

func GetNewImagePatternMap() *ImagePatternMap {
	imgPatternMap := new(ImagePatternMap)
	imgPatternMap.Mapping = make(map[string]string)
	return imgPatternMap
}

func (ipm *ImagePatternMap) AddImagePatternMap(pattern string, jobImage string) error {
	if ipm.Mapping == nil {
		return errors.New("ImagePatternMap is not initialized yet")
	}
	ipm.Mapping[pattern] = jobImage
	return nil
}

func (ipm *ImagePatternMap) DeleteImagePatternMap(pattern string) error {
	if ipm.Mapping == nil {
		return errors.New("ImagePatternMap is not initialized yet")
	}
	delete(ipm.Mapping, pattern)
	return nil
}

func (ipm *ImagePatternMap) ListImagePatternMap() (map[string]string, error) {
	if ipm.Mapping == nil {
		return nil, errors.New("ImagePatternMap is not initialized yet")
	}
	return ipm.Mapping, nil
}
