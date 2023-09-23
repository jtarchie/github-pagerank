package main

import (
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
	"github.com/jtarchie/github-pagerank/crawl"
	"github.com/jtarchie/github-pagerank/rank"
)

type CLI struct {
	Crawl crawl.Cmd `cmd:""`
	Rank  rank.Cmd  `cmd:""`
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, nil)))

	cli := &CLI{}
	ctx := kong.Parse(cli)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
