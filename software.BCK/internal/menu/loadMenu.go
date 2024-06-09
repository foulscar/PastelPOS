package menu

import (
  "os"
  "io"
  "encoding/json"
)

func LoadMenuFromJSONFile(filePath string) (*Menu, error) {
  jsonFile, err := os.Open(filePath)
	defer jsonFile.Close()
  if err != nil {
    return nil, err
  }

	byteValue, _ := io.ReadAll(jsonFile)

	var menu Menu
	json.Unmarshal([]byte(byteValue), &menu)
	return &menu, nil
}
