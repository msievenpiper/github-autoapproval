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

	// These examples assume `gh` is installed and has been authenticated.

	// Shell out to a gh command and read its output.

	// Use an API client to retrieve repository tags.
	// client, err := api.DefaultRESTClient()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// response := []struct {
	// 	Name string
	// }{}
	// err = client.Get("repos/protectednet/dashboard/tags", &response)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(response)
}
