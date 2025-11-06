package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/cli/go-gh/v2"
)

type PullRequest struct {
	Repo        string
	Number      string
	Title       string
	Branch      string
	Status      string
	CreatedDate string
	State       PullRequestState
}

type PullRequestState struct {
	LatestReviews []string `json:latestReviews`
	State         string   `json:state`
	author        string   `json:author`
}

func (p PullRequest) GetUrl() string {
	return "https://github.com/" + p.Repo + "/pull/" + p.Number
}

func (p PullRequest) GetBranchUrl() string {
	return "https://github.com/" + p.Repo + "/branch/" + p.Branch
}

func (p PullRequest) GetRepoUrl() string {
	return "https://github.com/" + p.Repo
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

func (c PullRequestContainer) GetItem(id string) (found PullRequest, notFound bool) {
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

func GetPullRequests(repo string, branch string) PullRequestContainer {
	prs, r, err := gh.Exec("pr", "list", "--repo", repo, "--search", branch)
	if err != nil {
		fmt.Println("Failed to get status for pr")
		fmt.Println("approimate cmd: gh pr list --repo " + repo + " --search " + branch)
		fmt.Println(r.String())
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
		pr.Repo = repo
		pr.State = getPullRequestStatus(pr)

		pullRequests.AddItem(pr)
	}

	return pullRequests
}

func getPullRequestStatus(pr PullRequest) PullRequestState {
	status, r, err := gh.Exec("pr", "view", pr.Number, "--repo", pr.Repo, "--json", "latestReviews,state,author")

	if err != nil {
		fmt.Println("Failed to get status for pr")
		fmt.Println("approimate cmd: gh pr view " + pr.Number + " --repo " + pr.Repo + " --json latestReviews,state,author")
		fmt.Println(r.String())
		log.Fatal(err)
	}

	var s PullRequestState

	json.NewDecoder(strings.NewReader(status.String())).Decode(&s)

	return s
}

func ApprovePullRequest(pr PullRequest, probe bool) bool {
	if strings.Contains(strings.Join(pr.State.LatestReviews, ","), "approve") {
		return true
	}

	if !probe {
		_, r, err := gh.Exec("pr", "review", pr.Number, "--repo", pr.Repo, "--approve")

		if err != nil {
			fmt.Println("Failed to approve")
			fmt.Println("approimate cmd: gh pr review " + pr.Number + " --repo " + pr.Repo + " --approve")
			fmt.Println(r.String())
			log.Fatal(err)
		}
	}

	fmt.Println("Pull Request Approved - " + pr.GetUrl())

	return true
}
