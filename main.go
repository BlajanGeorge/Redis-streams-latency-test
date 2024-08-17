package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"strings"
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
			for {
				var cursorId string
				var currentMessage []redis.XStream
				if cursorId == "" {
					currentMessage, _ = rdb.XRead(ctx, &redis.XReadArgs{Streams: []string{key, "$"}, Block: 0, Count: 1}).Result()
				} else {
					currentMessage, _ = rdb.XRead(ctx, &redis.XReadArgs{Streams: []string{key, cursorId}, Block: 0, Count: 1}).Result()
				}
				cursorId = currentMessage[0].Messages[0].ID
				for key := range currentMessage[0].Messages[0].Values {
					recordTime, _ := strconv.ParseInt(key, 10, 64)
					toRedisTime, _ := strconv.ParseInt(strings.Split(cursorId, "-")[0], 10, 64)
					//fmt.Printf("%s\n", strconv.FormatInt(recordTime, 10))
					//fmt.Printf("%s\n", strconv.FormatInt(time.Now().UnixMilli(), 10))
					fmt.Printf("latency to consumer: %s\n", strconv.FormatInt(time.Now().UnixMilli()-recordTime, 10))
					fmt.Printf("latency to redis: %s\n", strconv.FormatInt(time.Now().UnixMilli()-toRedisTime, 10))
				}
			}
		}(fmt.Sprintf("consumer-%d", i))
	}

	//let consumers start first
	time.Sleep(1 * time.Second)

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
				time.Sleep(2 * time.Second)
			}
		}(fmt.Sprintf("producer-%d", i))
	}
	wg.Wait()

}
