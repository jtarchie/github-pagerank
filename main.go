package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"golang.org/x/time/rate"
)

type CLI struct {
	DBFilename        string        `required:"" default:":memory:" help:"the name of the file to save results to"`
	GithubAPIKey      string        `required:"" help:"API Key for the Github API" env:"GITHUB_TOKEN"`
	RateLimitInterval time.Duration `default:"1h" help:"rate limit time duration"`
	RateLimitRequest  uint64        `default:"500" help:"number of requests to do per rate limit interval"`
	ResultLimit       int           `required:"" default:"100" help:"the number of results to be returned from GraphQL"`
	Username          string        `required:"" help:"Username to start finding associations from"`
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, nil)))

	cli := &CLI{}
	ctx := kong.Parse(cli)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}

func (c *CLI) Run() error {
	ctx := context.Background()

	service, err := NewService(ctx, c.DBFilename)
	if err != nil {
		return fmt.Errorf("could not create service: %w", err)
	}

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)
	limiter := rate.NewLimiter(
		rate.Every(
			c.RateLimitInterval/time.Duration(c.RateLimitRequest),
		),
		int(c.RateLimitRequest),
	)

	err = service.SetUsername(ctx, c.Username)
	if err != nil {
		return fmt.Errorf("could not set username to start from: %w", err)
	}

	for {
		limiter.Wait(ctx)

		username, err := service.NextUsername(ctx)
		if err != nil {
			return fmt.Errorf("could not find next username: %w", err)
		}

		slog.Info("getting information for user", slog.String("username", username))

		variables := map[string]interface{}{
			"count":    githubv4.Int(c.ResultLimit),
			"username": githubv4.String(username),
		}

		query := UserQuery{}

		err = client.Query(ctx, &query, variables)
		if err != nil {
			return fmt.Errorf("could not query for user %q: %w", username, err)
		}

		for _, follower := range query.User.Followers.Nodes {
			err = service.SetFollower(ctx, username, string(follower.Login))
			if err != nil {
				return fmt.Errorf("could not set follower %q -> %q: %w", follower.Login, username, err)
			}
		}
		for _, follower := range query.User.Followers.Nodes {
			err = service.SetFollower(ctx, string(follower.Login), username)
			if err != nil {
				return fmt.Errorf("could not set following %q -> %q: %w", username, follower.Login, err)
			}
		}

		err = service.SetProcessed(ctx, username)
		if err != nil {
			return fmt.Errorf("could not set processed: %w", err)
		}
	}
}
