package main

import (
	"gopkg.in/redis.v5"
	"log"
	"strconv"
)

// This file is using go-redis: https://github.com/go-redis/redis

func WriteCredentials(credentials APICredentials) {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	//Check Redis connection
	pong, err := client.Ping().Result()
	log.Println(pong, err)

	if err != nil {
		log.Println("ERROR: Could not connect to Redis")
	}

	// Store the access token and refresh token in Redis

	err = client.Set("line_access_token", credentials.AccessToken, 0).Err()

	if err != nil {
		panic(err)

	}

	err = client.Set("line_refresh_token", credentials.RefreshToken, 0).Err()

	if err != nil {
		panic(err)
	}

	err = client.Set("line_token_expire", strconv.FormatInt(credentials.Expire, 10), 0).Err()

	if err != nil {
		panic(err)
	}

}

func GetCredentials() APICredentials {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	//Check Redis connection
	pong, err := client.Ping().Result()
	log.Println(pong, err)

	if err != nil {
		log.Println("ERROR: Could not connect to Redis")
	}

	// Store the access token and refresh token in Redis

	accessToken, err := client.Get("line_access_token").Result()
	if err == redis.Nil {
		// set accessToken to empty string if it doesn't exist
		accessToken = ""
	} else if err != nil {
		panic(err)
	}

	refreshToken, err := client.Get("line_refresh_token").Result()
	if err == redis.Nil {
		// set refreshToken to empty string if it doesn't exist
		refreshToken = ""
	} else if err != nil {
		panic(err)
	}

	expire, err := client.Get("line_token_expire").Result()
	if err == redis.Nil {
		// set refreshToken to empty string if it doesn't exist
		expire = "0"
	} else if err != nil {
		panic(err)
	}

	expireInt, _ := strconv.Atoi(expire)

	returnCredentials := APICredentials{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expire:       int64(expireInt),
	}

	return returnCredentials

}

func ClearCredentials() {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	//Check Redis connection
	pong, err := client.Ping().Result()
	log.Println(pong, err)

	if err != nil {
		log.Println("ERROR: Could not connect to Redis")
	}

	// Delete the access token

	err = client.Del("line_access_token").Err()
	if err != nil {
		panic(err)
	}

	// Delete the refresh token

	err = client.Del("line_token_expire").Err()
	if err != nil {
		panic(err)
	}

	err = client.Del("line_refresh_token").Err()
	if err != nil {
		panic(err)
	}

}
