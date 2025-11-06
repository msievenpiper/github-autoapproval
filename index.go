package main

import (
	"github-autoapproval/v2/internal"
	"log"
)

func main() {
	// Check auth first
	_, err := internal.GetAuthState()

	if err != nil {
		log.Fatal("There is no auth available")
	}

	input := internal.GetInputs()

	if len(input.Repos) == 0 {
		log.Fatal("There are no repositories available")
	}

	for _, repo := range input.Repos {
		reqs := internal.GetPullRequests(repo, input.Branch)

		for _, req := range reqs.Requests {
			internal.ApprovePullRequest(req, input.Probe)

			if !input.Probe && input.Merge {
				internal.MergePullRequest(req)
			}
		}
	}
}
