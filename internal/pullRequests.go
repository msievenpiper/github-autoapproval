package internal

import (
	"fmt"
	"log"
	"strings"

	"github.com/cli/go-gh/v2"
)

type PullRequest struct {
	Number      string
	Title       string
	Branch      string
	Status      string
	CreatedDate string
}

type PullRequestContainer struct {
	Requests []PullRequest
}

func (c *PullRequestContainer) AddItem(pr PullRequest) {
	c.Requests = append(c.Requests, pr)
}

func (c *PullRequestContainer) RemoveItem(pr PullRequest) *PullRequestContainer {
	f := func(item PullRequest) bool { return pr.Number != item.Number }

	c.Requests = Filter(c.Requests, f)
	return c
}

func (c *PullRequestContainer) GetItem(id string) (found PullRequest, notFound bool) {
	f := func(item PullRequest) bool { return id == item.Number }
	res := Filter(c.Requests, f)

	if len(res) != 1 {
		found = PullRequest{}
		notFound = true
		return
	}

	found = res[0]
	notFound = false
	return
}

func GetPullRequests(repo string) PullRequestContainer {
	prs, _, err := gh.Exec("pr", "list", "--repo", repo, "--search", "translation-github-action")
	if err != nil {
		log.Fatal(err)
	}

	data := strings.Split(prs.String(), "\n")
	filterTest := func(item string) bool { return item != "" }
	data = Filter(data, filterTest)

	pullRequests := PullRequestContainer{}

	for _, element := range data {
		pr := PullRequest{}
		parts := strings.Split(element, "\t")

		if len(parts) != 5 {
			continue
		}

		pr.Number = parts[0]
		pr.Title = parts[1]
		pr.Branch = parts[2]
		pr.Status = parts[3]
		pr.CreatedDate = parts[4]

		pullRequests.AddItem(pr)
	}

	return pullRequests
}

func ApprovePullRequest(repo string, pr PullRequest) bool {
	_, _, err := gh.Exec("pr", "review", pr.Number, "--repo", repo, "--approve")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Pull Request Approved - https://github.com/" + repo + "/pull/" + pr.Number)

	return true
}
