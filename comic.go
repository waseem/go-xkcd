// xkcd: downloads json of a comic given its number and stores it under index/
// if it cannot find it under index/
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	if _, err := strconv.Atoi(comicNumber); err != nil {
		fmt.Fprintf(os.Stderr, "comic number cannot be 0")
		os.Exit(1)
	}
	indexDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fileName := comicNumber + ".json"
	filePath := indexDir + "/index/" + fileName

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
		var err error
		response, err := http.Get(baseUrl + comicNumber + "/info.0.json")
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

		file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0660)
		if err != nil {
			return nil, err
		}
		_, err = file.Write(bodyBytes)
		if err != nil {
			file.Close()
			return nil, err
		}

		if err = json.Unmarshal(bodyBytes, &comic); err != nil {
			return nil, err
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

func main() {
	comic, err := GetComic(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Comic: %d\n", comic.Num)
	fmt.Printf("+----------+\n\n")
	fmt.Printf("Title: %s\n", comic.Title)
	fmt.Printf("Safe Title: %s\n", comic.SafeTitle)
}
