package main

import (
	"encoding/json"
	"math/rand"
	"slices"

	"github.com/gofiber/fiber/v2"
)

var multiplier = map[int]map[int]float32 {
    1: {
        1: 5.88,
        2: 2.94,
        3: 1.96,
        4: 1.47,
        5: 1.17,
    },
    2: {
        1: 35.28,
        2: 8.82,
        3: 3.92,
        4: 2.20,
        5: 1.41,
    },
    3: {
        1: 211.68,
        2: 26.46,
        3: 7.84,
        4: 3.30,
        5: 1.69,
    },
}

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
    Won            bool
    Amount         float32
    WinningColors  []int
    SelectedColors []int
    Error          string
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

    response.Won = true
    response.Amount = body.Amount * multiplier[body.Cubes][len(body.SelectedColors)]

    user.IncreaseBalance(response.Amount)
    wonResponse, _ := json.Marshal(response)
    return context.Send(wonResponse)
}
