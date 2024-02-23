package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"sync"
	"time"
)

var ctx = context.Background()

func main() {
	var host, pass, key string
	var port, db, producers, consumers int
	var wg sync.WaitGroup

	flag.StringVar(&host, "host", "localhost", "")
	flag.IntVar(&port, "port", 6379, "")
	flag.StringVar(&pass, "password", "", "")
	flag.IntVar(&db, "db", 0, "")
	flag.IntVar(&producers, "producers", 1, "")
	flag.IntVar(&consumers, "consumers", 1, "")
	flag.StringVar(&key, "key", "streamKey", "")

	for i := 0; i < consumers; i++ {
		wg.Add(1)

		go func(consumerName string) {
			defer wg.Done()

			rdb := redis.NewClient(&redis.Options{
				Addr:     fmt.Sprintf("%s:%d", host, port),
				Password: pass,
				DB:       db,
			})
			fmt.Printf("%s start\n", consumerName)
			initialMessage := rdb.XRead(ctx, &redis.XReadArgs{Streams: []string{key, "$"}, Count: 1, Block: 0})

			for {
				cursorId := initialMessage.Val()[0].Messages[0].ID
				streams := rdb.XRead(ctx, &redis.XReadArgs{Streams: []string{key, cursorId}, Count: 1, Block: 0})
				for _, message := range streams.Val()[0].Messages {
					for key := range message.Values {
						recordTime, _ := strconv.ParseInt(key, 10, 64)
						//fmt.Printf("%s\n", strconv.FormatInt(recordTime, 10))
						//fmt.Printf("%s\n", strconv.FormatInt(time.Now().UnixMilli(), 10))
						fmt.Printf("latency: %s\n", strconv.FormatInt(time.Now().UnixMilli()-recordTime, 10))
						//time.Sleep(5 * time.Second)
					}
				}
			}

		}(fmt.Sprintf("consumer-%d", i))
	}

	for i := 0; i < producers; i++ {
		wg.Add(1)

		go func(producerName string) {
			defer wg.Done()

			rdb := redis.NewClient(&redis.Options{
				Addr:     fmt.Sprintf("%s:%d", host, port),
				Password: pass,
				DB:       db,
			})
			fmt.Printf("%s start\n", producerName)

			for {
				rdb.XAdd(ctx, &redis.XAddArgs{Stream: key, Values: []interface{}{time.Now().UnixMilli(), ""}})
			}

		}(fmt.Sprintf("producer-%d", i))
	}

	wg.Wait()
}
