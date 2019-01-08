package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func writeItemToFile(saveDir string, i Item) {
	path := filepath.Join(saveDir, fmt.Sprintf("%v-%d.json", i.Type, i.Id))

	b, err := json.Marshal(i)
	if err != nil {
		fmt.Println("failed to marshal item:", err)
		return
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("failed to open file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	nn, err := writer.Write(b)
	if err != nil {
		fmt.Println("failed to write file:", err)
		return
	}
	writer.Flush()

	log.Println(fmt.Sprintf("%d/%d bytes wrote to %v", len(b), nn, path))
}
