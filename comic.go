// xkcd: downloads json of a comic given its number and stores it under index/
// if it cannot find it under index/
package xkcd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type Comic struct {
	Year, Month, Day string
	Num              uint
	Link             string
	News             string
	SafeTitle        string `json:"safe_title"`
	Transcript       string
	Alt              string
	Img              string
	Title            string
}

const (
	baseUrl = "https://xkcd.com/"
)

func GetComic(comicNumber string) (*Comic, error) {
	comicNumberInt, err := strconv.Atoi(comicNumber)
	if err != nil {
		fmt.Fprintf(os.Stderr, "comic number must be an integer")
		os.Exit(1)
	}

	filePath, err := prepareFilePath(uint(comicNumberInt))
	if err != nil {
		return nil, err
	}

	var comic Comic

	if fileExists(filePath) {
		bytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(bytes, &comic)
		if err != nil {
			return nil, err
		}

	} else {

		comicJson, err := requestComicJson(comicNumber)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(comicJson, &comic); err != nil {
			return nil, err
		}

		filePath, err = prepareFilePath(comic.Num)
		if err != nil {
			return nil, err
		}

		if !fileExists(filePath) {
			// in case of requesting 0th(current comic), the json file may or may not exist on disk
			if err = writeComicJson(filePath, comicJson); err != nil {
				return nil, err
			}
		}
	}

	return &comic, nil
}

func fileExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func prepareUrl(comicNumber string) string {
	if comicNumber == "0" {
		return baseUrl + "/info.0.json"
	}

	return baseUrl + comicNumber + "/info.0.json"
}

func requestComicJson(comicNumber string) ([]byte, error) {
	response, err := http.Get(prepareUrl(comicNumber))
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		response.Body.Close()
		return nil, fmt.Errorf("fetching comic %s failed", comicNumber)
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		response.Body.Close()
		return nil, err
	}

	response.Body.Close()
	return bodyBytes, nil
}

func prepareFilePath(comicNumber uint) (string, error) {
	currentDirectory, err := os.Getwd()
	if err != nil {
		return "", err
	}
	fileName := strconv.Itoa(int(comicNumber)) + ".json"
	return currentDirectory + "/index/" + fileName, nil
}

func writeComicJson(filePath string, comicJson []byte) error {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	_, err = file.Write(comicJson)
	if err != nil {
		file.Close()
		return err
	}
	return nil
}
