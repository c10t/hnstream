package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	/*
		news, _ := GetNewStories()
		log.Println(fmt.Sprintf("Result: %v, %v, %v, ...", news[0], news[1], news[2]))

		item, _ := GetItem(news[0])
		log.Println(fmt.Sprintf("Story %v: %v", item.Id, item.Title))

		before, _ := theItemExisted(item.Id)
		log.Println(fmt.Sprintf("JSON existed? -> %v", before))
		writeItemToFile("resources", item)
		after, _ := theItemExisted(item.Id)
		log.Println(fmt.Sprintf("JSON existed? -> %v", after))
	*/

	// ---
	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)

	go func() {
		<-signalChan
		cancel()
		log.Println("main: stopping...")
	}()

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	ids := make(chan int)
	go streamNewStories(ctx, ids)
	writeStories(ids)
}

func writeStories(ids <-chan int) {
	for id := range ids {
		item, err := GetItem(id)
		if err != nil {
			log.Println("failed to get item:", err)
		}
		existed, err := theItemExisted(item.Id)
		if err != nil {
			log.Println("failed to check if item is already exist:", err)
		}
		if !existed && err == nil {
			writeItemToFile("resources", item)
		}
		time.Sleep(10 * time.Second)
	}
	log.Println("[Writer] stopped")
}

func streamNewStories(ctx context.Context, ids chan<- int) {
	defer close(ids)
	for {
		log.Println("start to read new stories...")
		readNewStories(ctx, 1*time.Minute, ids)
		log.Println("--- (waiting) ---")
		select {
		case <-ctx.Done():
			log.Println("stop to read new stories...")
			return
		case <-time.After(5 * time.Second):
		}
		log.Println("repeat to read new stories...")
	}
}

func readNewStories(ctx context.Context, timeout time.Duration, ids chan<- int) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan struct{})
	defer func() { <-done }()

	news, err := GetNewStories()
	if err != nil {
		log.Println("Failed to get new stories:", err)
		return
	}

	go func() {
		defer close(done)
		for _, id := range news[:10] {
			exist, err := theItemExisted(id)
			if err != nil {
				log.Println("failed to check if the item is existed:", err)
			}
			if !exist {
				log.Println("send new id to channel:", id)
				ids <- id
			}
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("received ctx.Done()")
	case <-done:
		log.Println("received done")
	}
}

func theItemExisted(id int) (bool, error) {
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
