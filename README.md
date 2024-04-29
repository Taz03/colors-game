Backend implementation of [colors gambling game](https://www.hacksawgaming.com/games/colors)

### Tech stack: Go, Fiber, Redis

## API Docs
1. `POST /colors/bet`
   - Body
     ```json
     {
       "username": "",
       "cubes": 0,
       "selected_colors": [],
       "amount": 0.0
     }
     ```

2. `GET /balance`
    - Query Parameters
      - `username`: Username of the player to fetch balance for
