# Mancala

Mancala game implementation in Go as an API.

## Usage

To run the server locally run the following command:

```sh
make setup
make local.run
```

Once you have the server running.

You can access the API docs at [localhost:1323/docs](http://localhost:1323/docs)

To create a new game:

```bash
curl -X POST localhost:1323/v1/games -H "Content-Type: application/json" --data '{"player1":"Rick","player2":"Morty"}'
```

To show the state of a game:

```bash
curl localhost:1323/v1/games/:id
```

To perform the next play:

```bash
curl -X PATCH localhost:1323/v1/games/:id -H "Content-Type: application/json" --data '{"pit_index":0}'
```

### Notes

Games expire after 2 hours.

## License

Copyright 2023 [Pablo Crivella](https://pcriv.com).
Read [LICENSE](LICENSE) for details.
