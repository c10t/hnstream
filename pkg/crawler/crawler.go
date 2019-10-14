package crawler

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	a "github.com/c10t/gohn/pkg/accessmydata"
)

func Crawl() {
	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)

	go func() {
		<-signalChan
		cancel()
		log.Println("main: stopping...")
	}()

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	var myData a.AccessMyData
	myData = a.UseLocalFiles

	ids := make(chan int)
	go streamNewStories(ctx, ids, myData)

	writeStories(ids, myData)
}

func writeStories(ids <-chan int, myData a.AccessMyData) {
	for id := range ids {
		item, err := GetItem(id)
		if err != nil {
			log.Println("failed to get item:", err)
		}
		existed, err := myData.Exist(item.Id)
		if err != nil {
			log.Println("failed to check if item is already exist:", err)
		}
		if !existed && err == nil {
			myData.Write("resources", item)
			time.Sleep(3 * time.Second)
		}
	}
	log.Println("[Writer] stopped")
}

func streamNewStories(ctx context.Context, ids chan<- int, myData a.AccessMyData) {
	defer close(ids)
	for {
		log.Println("start to read new stories...")
		readNewStories(ctx, 1*time.Minute, ids, myData)
		log.Println("--- (waiting) ---")
		select {
		case <-ctx.Done():
			log.Println("stop reading new stories...")
			return
		case <-time.After(5 * time.Second):
		}
		log.Println("repeat to read new stories...")
	}
}

func readNewStories(ctx context.Context, timeout time.Duration, ids chan<- int, myData a.AccessMyData) {
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
		for _, id := range news {
			exist, err := myData.Exist(id)
			if err != nil {
				log.Println("failed to check if the item is existed:", err)
			}
			if !exist {
				log.Println("send new id to channel:", id)
				ids <- id
			}
			select {
			case <-ctx.Done():
				log.Println("received ctx.Done() so stop iteration")
				return
			default:
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
