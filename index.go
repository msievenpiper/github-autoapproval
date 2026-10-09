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
		reqs, err := internal.GetPullRequests(repo, input.Branch)
		if err != nil {
			log.Println("Skipping repo " + repo + ": " + err.Error())
			continue
		}

		for _, req := range reqs.Requests {
			if !internal.ApprovePullRequest(req, input.Probe) {
				continue
			}

			if !input.Probe && input.Merge {
				internal.MergePullRequest(req, input.MergeStrategy)
			}
		}
	}
}
