package handlers

import (
	"encoding/json"
	"errors"
	"log"
)

type ImagePatternMap struct {
	Mapping map[string]string
}

// GetNewImagePatternMap constructs an ImagePatternMap instance
func GetNewImagePatternMap() *ImagePatternMap {
	imgPatternMap := new(ImagePatternMap)
	imgPatternMap.Mapping = make(map[string]string)
	return imgPatternMap
}

// AddImagePatternMapFromJson adds pattern-image maps from an json string,
// Ex. {
//		"pattern1": "quay.io/cdis/image1",
//      "pattern2": "quay.io/cdis/image2"
//	}
func (ipm *ImagePatternMap) AddImagePatternMapFromJson(jsonString string) error {
	if ipm.Mapping == nil {
		return errors.New("ImagePatternMap is not initialized yet")
	}
	var mapping map[string]string
	if err := json.Unmarshal([]byte(jsonString), &mapping); err != nil {
		log.Println(err)
	}

	for pattern, quayImage := range mapping {
		ipm.AddImagePatternMap(pattern, quayImage)
	}

	return nil
}

// AddImagePatternMap adds a pattern-image map
func (ipm *ImagePatternMap) AddImagePatternMap(pattern string, jobImage string) error {
	if ipm.Mapping == nil {
		return errors.New("ImagePatternMap is not initialized yet")
	}
	ipm.Mapping[pattern] = jobImage
	return nil
}

// DeleteImagePatternMap deletes a pattern-image map
func (ipm *ImagePatternMap) DeleteImagePatternMap(pattern string) error {
	if ipm.Mapping == nil {
		return errors.New("ImagePatternMap is not initialized yet")
	}
	delete(ipm.Mapping, pattern)
	return nil
}

// ListImagePatternMap lists all pattern-image maps
func (ipm *ImagePatternMap) ListImagePatternMap() (map[string]string, error) {
	if ipm.Mapping == nil {
		return nil, errors.New("ImagePatternMap is not initialized yet")
	}
	return ipm.Mapping, nil
}
