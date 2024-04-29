package main

import "github.com/gofiber/fiber/v2"

func main() {
    app := fiber.New()

    app.Post("/colors/bet", ColorsBet)
    app.Get("/balance", GetUserBalance)

    app.Listen(":8080")
}
