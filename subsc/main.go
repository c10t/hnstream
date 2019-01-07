package main

import (
	"fmt"
	"log"
)

func main() {
	tops, _ := GetTopStories()
	log.Println(fmt.Sprintf("Result: %v, %v, %v, ...", tops[0], tops[1], tops[2]))

	item, _ := GetItem(tops[0])
	log.Println(fmt.Sprintf("Story %v: %v", item.Id, item.Title))
}
