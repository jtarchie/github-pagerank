package crawl

import "github.com/shurcooL/githubv4"

type UserQuery struct {
	User struct {
		// Followers struct {
		// 	TotalCount githubv4.Int
		// 	Nodes      []struct {
		// 		Login githubv4.String
		// 	}
		// } `graphql:"followers(first: $count)"`
		Following struct {
			TotalCount githubv4.Int
			Nodes      []struct {
				Login             githubv4.String
				IsFollowingViewer githubv4.Boolean
			}
		} `graphql:"following(first: $count)"`
	} `graphql:"user(login: $username)"`
}
