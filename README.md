# Mancala

Mancala game implementation in Go as an API.

## Usage

[API Documentation](https://<URL>/docs)

To create a new game:

```bash
curl -X POST https://<URL>/v1/games -H "Content-Type: application/json" --data '{"player1":"Rick","player2":"Morty"}'
```

To show the state of a game:

```bash
curl https://<URL>/v1/games/:id
```

To perform the next play:

```bash
curl -X PATCH https://<URL>/v1/games/:id -H "Content-Type: application/json" --data '{"pit_index":0}'
```

### Notes

Games expire after 2 hours.

## License

Copyright 2023 [Pablo Crivella](https://pcriv.com).
Read [LICENSE](LICENSE) for details.
