package main

import (
	"context"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

const DEFAULT_BALANCE = 5000.0

var ctx = context.Background()
var redisClient = redis.NewClient(&redis.Options {
    Addr:     "localhost:6379",
    Password: "",
    DB:       0,
})

type User struct {
    Username string
    Balance  float32
}

func GetUser(username string) User {
    balance, error := redisClient.Get(ctx, username).Float32()
    if error == redis.Nil {
        balance = DEFAULT_BALANCE
    }

    redisClient.Set(ctx, username, balance, 0)

    return User {
        Username: username,
        Balance:  balance, 
    }
}

func (user *User) IncreaseBalance(amount float32) {
    user.Balance += amount
    redisClient.Set(ctx, user.Username, user.Balance, 0)
}

func (user *User) DecreaseBalance(amount float32) {
    user.Balance -= amount
    redisClient.Set(ctx, user.Username, user.Balance, 0)
}

func GetUserBalance(context *fiber.Ctx) error {
    username := context.Query("username")

    balance, error := redisClient.Get(ctx, username).Float32()
    if error == redis.Nil {
        return context.SendStatus(404)
    }

    payload := User {
        Username: username,
        Balance:  balance,
    }
    jsonResponse, _ := json.Marshal(payload)
    return context.Send(jsonResponse)
}
