package helper

import (
	"encoding/csv"
	"os"
)

//WriteCSVFile function for writing CSV Files
//It will check if file exists then Write to File
//If not then it will create file and then Write To File
func WriteCSVFile(path string, data [][]string) error {

	fileDescriptor, err := os.Create(path)
	if err != nil {
		return err
	}

	csvWriter := csv.NewWriter(fileDescriptor)
	csvWriter.WriteAll(data)
	csvWriter.Flush()
	return nil
}
