package Service

import "io/ioutil"

func saveFileToUploads(file []byte, filename string) error {
	filePath := "uploads/" + filename
	err := ioutil.WriteFile(filePath, file, 0644)
	if err != nil {
		return err
	}
	return nil
}
