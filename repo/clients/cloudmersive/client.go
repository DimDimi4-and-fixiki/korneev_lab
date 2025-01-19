package cloudmersive

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type client struct {
	httpC *http.Client
}

func NewClient() *client {
	return &client{
		httpC: &http.Client{},
	}
}

func (c *client) CheckFile(filePath string) (string, error) {
	url := "https://api.cloudmersive.com/virus/scan/file/advanced"
	method := "POST"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	part1, err := writer.CreateFormFile("inputFile", filepath.Base("/path/to/file"))
	if err != nil {
		return "", err
	}

	_, err = io.Copy(part1, file)
	if err != nil {
		return "", err
	}

	err = writer.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return "", err
	}

	req.Header.Add("allowExecutables", "<boolean>")
	req.Header.Add("allowInvalidFiles", "<boolean>")
	req.Header.Add("allowScripts", "<boolean>")
	req.Header.Add("allowPasswordProtectedFiles", "<boolean>")
	req.Header.Add("allowMacros", "<boolean>")
	req.Header.Add("restrictFileTypes", "<string>")
	req.Header.Add("Content-Type", "multipart/form-data")
	req.Header.Add("Apikey", "YOUR-API-KEY-HERE")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := c.httpC.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(body), nil
}
