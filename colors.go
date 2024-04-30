package main

import (
	"encoding/json"
	"math"
	"math/rand"
	"slices"

	"github.com/gofiber/fiber/v2"
)

const TOTAL_COLORS = 6.0

type Body struct {
    // Username of the player
    Username string `json:"username"`

    // Number of cubes
    Cubes int `json:"cubes"`

    // Array of index(s) of selected colors 0-indexed
    SelectedColors []int `json:"selected_colors"`

    // The bet amount
    Amount float32 `json:"amount"`
}

type Response struct {
    // Winning status
    Won bool `json:"won"`

    // Amount lost/won
    Amount float32 `json:"amount"`

    // Amount multiplier in case of win
    Multiplier float32 `json:"multiplier"`

    WinningColors []int `json:"winning_colors"`

    SelectedColors []int `json:"selected_colors"`

    // Error message in case of erro
    Error string `json:"error"`
}

func calculateMultiplier(selectedColors, cubes int) float32 {
    return float32(math.Pow(TOTAL_COLORS / float64(selectedColors), float64(cubes)))
}

func ColorsBet(context *fiber.Ctx) error {
    var body Body
    json.Unmarshal(context.Body(), &body)

    if len(body.SelectedColors) < 1 || len(body.SelectedColors) > 5 {
        errorResponse, _ := json.Marshal(Response{Error: "The selected colors should be between 1-5"})
        return context.Send(errorResponse)
    }

    user := GetUser(body.Username)
    user.DecreaseBalance(body.Amount)

    var randomColors []int
    for i := 0; i < body.Cubes; i++ {
        randomColors = append(randomColors, rand.Intn(6))
    }

    response := Response {
        WinningColors:  randomColors,
        SelectedColors: body.SelectedColors,
    }
    for _, randomColor := range randomColors {
        if !slices.Contains(body.SelectedColors, randomColor) {
            response.Won = false
            response.Amount = body.Amount

            lostResponse, _ := json.Marshal(response)
            return context.Send(lostResponse)
        }
    }

    multiplier := calculateMultiplier(len(body.SelectedColors), body.Cubes)

    response.Won = true
    response.Multiplier = multiplier
    response.Amount = body.Amount * multiplier

    user.IncreaseBalance(response.Amount)
    wonResponse, _ := json.Marshal(response)
    return context.Send(wonResponse)
}
