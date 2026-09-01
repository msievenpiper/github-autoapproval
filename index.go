package main

import (
	"github-autoapproval/v2/internal"
	"log"
)

func main() {
	input := internal.GetInputs()

	if input.Init {
		internal.RunConfigWizard()
		return
	}

	// Check auth first
	_, err := internal.GetAuthState()

	if err != nil {
		log.Fatal("There is no auth available")
	}

	if len(input.Repos) == 0 {
		log.Fatal("There are no repositories available")
	}

	lock, ok := internal.AcquireRunLock()
	if !ok {
		log.Println("Another instance is already running, skipping this run")
		return
	}
	defer lock.Release()

	for _, repo := range input.Repos {
		reqs := internal.GetPullRequests(repo, input.Branch)

		for _, req := range reqs.Requests {
			internal.ApprovePullRequest(req, input.Probe)

			if !input.Probe && input.Merge {
				internal.MergePullRequest(req, input.MergeStrategy)
			}
		}
	}
}
