package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const base = "https://hacker-news.firebaseio.com/v0/"
const item = base + "item/%d.json"
const topstories = base + "topstories.json"
const newstories = base + "newstories.json"

func getHackerNews(url string) ([]byte, error) {
	log.Println("Try to request to " + url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Failed to create request:", err)
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println("Failed to request:", err)
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Failed to read body:", err)
		return nil, err
	}

	return body, nil
}

func GetItem(id int) (Item, error) {
	i := Item{}

	res, err := getHackerNews(fmt.Sprintf(item, id))
	if err != nil {
		log.Println("Failed to get Item:", err)
		return i, err
	}

	err = json.Unmarshal(res, &i)
	if err != nil {
		log.Println("Failed to unmarshall response:", err)
		return i, err
	}

	return i, nil
}

func GetTopStories() ([]int, error) {
	res, err := getHackerNews(topstories)
	if err != nil {
		log.Println("Failed to get Top Stories:", err)
		return nil, err
	}

	var tops []int
	err = json.Unmarshal(res, &tops)
	if err != nil {
		log.Println("Failed to unmarshall response:", err)
		return nil, err
	}

	return tops, nil
}
