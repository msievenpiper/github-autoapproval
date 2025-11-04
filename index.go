package main

import (
	"github-autoapproval/v2/internal"
)

func main() {
	availableRepos := [...]string{"protectednet/dashboard", "protectednet/www-total", "protectednet/legal"}

	for _, repo := range availableRepos {
		reqs := internal.GetPullRequests(repo)

		for _, req := range reqs.Requests {
			internal.ApprovePullRequest(repo, req)
		}
	}
}
