package crawl

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type Cmd struct {
	DBFilename   string        `required:"" default:":memory:" help:"the name of the file to save results to"`
	GithubAPIKey string        `required:"" help:"API Key for the Github API" env:"GITHUB_TOKEN"`
	WaitInterval time.Duration `default:"1s" help:"wait between each request"`
	ResultLimit  int           `required:"" default:"100" help:"the number of results to be returned from GraphQL"`
	MaxFollowing int           `required:"" default:"510" help:"do not include user if they follow more than max users"`
	Username     string        `required:"" help:"Username to start finding associations from"`
}

func (c *Cmd) Run() error {
	ctx := context.Background()

	service, err := NewService(ctx, c.DBFilename)
	if err != nil {
		return fmt.Errorf("could not create service: %w", err)
	}
	defer service.Close()

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.GithubAPIKey},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	err = service.SetUsername(ctx, c.Username)
	if err != nil {
		return fmt.Errorf("could not set username to start from: %w", err)
	}

	for {
		time.Sleep(c.WaitInterval)

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
			if strings.Contains(err.Error(), "Could not resolve to a User with the login") {
				slog.Error("user not a user", slog.String("username", username), slog.String("error", err.Error()))

				err = service.SetProcessed(ctx, username)
				if err != nil {
					return fmt.Errorf("could not set processed: %w", err)
				}
				continue
			}
			return fmt.Errorf("could not query for user %q: %w", username, err)
		}

		for _, follower := range query.User.Following.Nodes {
			// When a user has less users they are following
			if query.User.Following.TotalCount <= githubv4.Int(c.MaxFollowing) {
				err = service.SetFollower(ctx, username, string(follower.Login))
				if err != nil {
					return fmt.Errorf("could not set following %q -> %q: %w", username, follower.Login, err)
				}
			} else {
				err = service.SetUsername(ctx, string(follower.Login))
				if err != nil {
					return fmt.Errorf("could not set user: %w", err)
				}
			}
		}

		err = service.SetProcessed(ctx, username)
		if err != nil {
			return fmt.Errorf("could not set processed: %w", err)
		}
	}
}
