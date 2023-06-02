package main

import (
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

func main() {

	client := redis.NewClient(&redis.Options{
		Addr:     "10.0.0.31:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pong)

	defer client.Close()

	fmt.Println("Connected to Redis!")

}
