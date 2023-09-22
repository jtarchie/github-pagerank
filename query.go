package main

import "github.com/shurcooL/githubv4"

type UserQuery struct {
	User struct {
		Followers struct {
			Nodes []struct {
				Login githubv4.String
			}
		} `graphql:"followers(first: $count)"`
		Following struct {
			Nodes []struct {
				Login githubv4.String
			}
		} `graphql:"following(first: $count)"`
	} `graphql:"user(login: $username)"`
}