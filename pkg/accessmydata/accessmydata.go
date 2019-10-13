package accessmydata

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type AccessMyData interface {
	Exist(id int) (bool, error)
	Write(saveTo string, i Item)
}

type LocalFiles struct{}

var UseLocalFiles LocalFiles

func (LocalFiles) Exist(id int) (bool, error) {
	// todo: consider a better way to check
	matches, err := filepath.Glob(
		filepath.Join("resources", fmt.Sprintf("*-%d.json", id)))

	if err != nil {
		log.Println("failed to glob:", err)
		return false, err
	}

	if len(matches) == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (LocalFiles) Write(saveTo string, i Item) {
	path := filepath.Join(saveTo, fmt.Sprintf("%v-%d.json", i.Type, i.Id))

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
