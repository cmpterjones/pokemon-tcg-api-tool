# pokemon-tcg-api-tool

This is just a simple toy application for querying the cards API at `api.pokemontcg.io`.

## Execution

You can run this however you want, perhaps directly?

```shell
go run main.go
```

Or maybe build it first?

```shell
go build
```

## Usage

This app makes use of `flags` to set params and provide some basic help output. Check it out:

```shell
go run main.go --help
```

(Can't really take credit for that, it's just in the stdlib)

It's made for a specific purpose, so the defaults (no arguments) are configured to cover that, but theoretically you can run any query you want, as long as it works with the Cards API. See [Search Cards](https://docs.pokemontcg.io/api-reference/cards/search-cards) for more information.

The output is kinda gnarly but it's legible. You will see it log each step of the process, but also the time elapsed in each step.

There are better ways to do this but this was a timeboxed effort!
