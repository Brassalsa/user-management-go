package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func UploadFile(r *http.Request) (string, error) {
	// retrieve file from data
	file, _, err := r.FormFile("avatar")

	if err != nil {
		fmt.Println("error retrieving file: ", err)
		return "", err
	}

	// close file when return
	defer file.Close()

	// make temp dir if doesn't exists
	err = os.MkdirAll("./static/img", os.ModePerm)
	if err != nil {
		return "", err
	}

	tempFile, err := os.CreateTemp("./static/img", "upload-*.png")

	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	// read all of contents of our uploaded file into byte arr
	fileBytes, err := io.ReadAll(file)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	tempFile.Write(fileBytes)
	s := ""
	s = strings.ReplaceAll(tempFile.Name(), "\\", "/")

	return s, nil
}

func DeleteFile(s string) error {
	if err := os.Remove(s); err != nil {
		return err
	}
	return nil
}
