package model

import "os"

func StringArrayAppend(arr []string, str string) []string {
	for _, v := range arr {
		if v == str {
			return arr
		}
	}
	arr = append(arr, str)
	return arr
}

func StringArrayExists(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

func FileWrite(filePath string, content string) error {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	_, err = f.WriteString(content)
	if err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}
