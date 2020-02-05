# Mancala

Mancala game implementation in Go as an API.

## Usage

If you want to use the game engine:


```go
import (
    "fmt"
    "github.com/pablocrivella/mancala/engine"
)

func main() {
    game := engine.NewGame("Rick", "Morty")

    // Player1 plays
    game.PlayTurn(0)
}
```

If you want to use the API go to https://go-mancala.herokuapp.com/docs

To create a new game:

```bash
curl -X POST https://go-mancala.herokuapp.com/v1/games -H "Content-Type: application/json" --data '{"player1":"Rick","player2":"Morty"}'
```

To show the state of a game:

```bash
curl https://go-mancala.herokuapp.com/v1/games/:id
```

To perform the next play:

```bash
curl -X PATCH https://go-mancala.herokuapp.com/v1/games/:id -H "Content-Type: application/json" --data '{"pit_index":0}'
```

### Notes

Games expire after 2 hours.

## License

Copyright 2020 [Pablo Crivella](https://pablocrivella.me).
Read [LICENSE](LICENSE) for details.