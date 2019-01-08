package main

import (
	"fmt"
	"log"
)

func main() {
	news, _ := GetNewStories()
	log.Println(fmt.Sprintf("Result: %v, %v, %v, ...", news[0], news[1], news[2]))

	item, _ := GetItem(news[0])
	log.Println(fmt.Sprintf("Story %v: %v", item.Id, item.Title))

	writeItemToFile(item)
}
