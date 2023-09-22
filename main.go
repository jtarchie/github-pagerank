package main

import (
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
)

type CLI struct {
	Crawl Crawl `cmd:""`
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, nil)))

	cli := &CLI{}
	ctx := kong.Parse(cli)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
