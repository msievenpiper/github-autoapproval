package internal

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cli/go-gh/v2"
)

type PullRequestStateReviewAuthor struct {
	Login string `json:"login"`
}

type PullRequestStateReview struct {
	Id                string                       `json:"id"`
	Author            PullRequestStateReviewAuthor `json:"author"`
	AuthorAssociation string                       `json:"authorAssociation"`
	Body              string                       `json:"body"`
	SubmittedAt       string                       `json:"submittedAt"`
	State             string                       `json:"state"`
}

func (s PullRequestStateReview) IsAppoved() bool {
	return strings.Contains(strings.ToLower(s.State), "approve")
}

type PullRequestState struct {
	LatestReviews []PullRequestStateReview `json:"latestReviews"`
	State         string                   `json:"state"`
	Author        string                   `json:"author"`
}

func (s PullRequestState) IsAppoved() bool {
	for _, el := range s.LatestReviews {
		if el.IsAppoved() {
			return true
		}
	}

	return false
}

type PullRequest struct {
	Repo        string
	Number      string
	Title       string
	Branch      string
	Status      string
	CreatedDate string
	State       PullRequestState
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

func (p PullRequest) IsAppoved() bool {
	return p.State.IsAppoved()
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

func GetPullRequests(repo string, branch string) (PullRequestContainer, error) {
	// --head matches on the PR's source branch; --search is free-text and also
	// matches unrelated PRs that merely mention the branch name in their body.
	prs, r, err := gh.Exec("pr", "list", "--repo", repo, "--head", branch)
	if err != nil {
		fmt.Println("Failed to get status for pr")
		fmt.Println("approimate cmd: gh pr list --repo " + repo + " --head " + branch)
		fmt.Println(r.String())
		return PullRequestContainer{}, err
	}

	data := strings.Split(prs.String(), "\n")
	filterTest := func(item string) bool { return item != "" }
	data = Filter(data, filterTest)

	pullRequests := PullRequestContainer{}

	for _, element := range data {
		pr := PullRequest{}
		parts := strings.Split(element, "\t")

		if len(parts) != 5 || parts[2] != branch {
			continue
		}

		pr.Number = parts[0]
		pr.Title = parts[1]
		pr.Branch = parts[2]
		pr.Status = parts[3]
		pr.CreatedDate = parts[4]
		pr.Repo = repo

		state, err := getPullRequestStatus(pr)
		if err != nil {
			fmt.Println("Skipping " + pr.GetUrl() + ": " + err.Error())
			continue
		}
		pr.State = state

		pullRequests.AddItem(pr)
	}

	return pullRequests, nil
}

func getPullRequestStatus(pr PullRequest) (PullRequestState, error) {
	status, r, err := gh.Exec("pr", "view", pr.Number, "--repo", pr.Repo, "--json", "latestReviews,state,author")

	if err != nil {
		fmt.Println("Failed to get status for pr")
		fmt.Println("approimate cmd: gh pr view " + pr.Number + " --repo " + pr.Repo + " --json latestReviews,state,author")
		fmt.Println(r.String())
		return PullRequestState{}, err
	}

	var s PullRequestState

	json.NewDecoder(strings.NewReader(status.String())).Decode(&s)

	return s, nil
}

func ApprovePullRequest(pr PullRequest, probe bool) bool {
	if pr.IsAppoved() {
		return true
	}

	if !probe {
		_, r, err := gh.Exec("pr", "review", pr.Number, "--repo", pr.Repo, "--approve")

		if err != nil {
			fmt.Println("Failed to approve")
			fmt.Println("approimate cmd: gh pr review " + pr.Number + " --repo " + pr.Repo + " --approve")
			fmt.Println(r.String())
			fmt.Println(err)
			return false
		}
	}

	fmt.Println("Pull Request Approved - " + pr.GetUrl())

	return true
}

// updatePullRequestBranch merges the base branch into the PR's branch,
// equivalent to clicking "Update branch" on GitHub.
func updatePullRequestBranch(pr PullRequest) error {
	_, r, err := gh.Exec("pr", "update-branch", pr.Number, "--repo", pr.Repo)
	if err != nil {
		return fmt.Errorf("%v: %s", err, r.String())
	}

	return nil
}

type prMergeState struct {
	Mergeable        string `json:"mergeable"`
	MergeStateStatus string `json:"mergeStateStatus"`
}

func getMergeState(pr PullRequest) (prMergeState, error) {
	var s prMergeState

	status, _, err := gh.Exec("pr", "view", pr.Number, "--repo", pr.Repo, "--json", "mergeable,mergeStateStatus")
	if err != nil {
		return s, err
	}

	err = json.NewDecoder(strings.NewReader(status.String())).Decode(&s)
	return s, err
}

// isBehindBase asks GitHub whether the PR's branch has fallen behind its base
// (e.g. because another PR merged first). Other merge failures such as pending
// required checks are not fixed by updating the branch, and updating anyway
// pushes a new commit that restarts CI.
func isBehindBase(pr PullRequest) bool {
	s, err := getMergeState(pr)
	return err == nil && s.MergeStateStatus == "BEHIND"
}

// waitForMergeable polls the PR until GitHub reports it as mergeable and no
// longer behind/dirty, or until timeout elapses.
func waitForMergeable(pr PullRequest, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		s, err := getMergeState(pr)

		if err == nil {
			if s.Mergeable == "MERGEABLE" && s.MergeStateStatus != "BEHIND" && s.MergeStateStatus != "DIRTY" && s.MergeStateStatus != "UNKNOWN" {
				return true
			}
		}

		time.Sleep(5 * time.Second)
	}

	return false
}

// MergePullRequest expects the caller to have already approved the PR via
// ApprovePullRequest; pr.State is fetched before approval so it can't be used
// to check that here.
func MergePullRequest(pr PullRequest, strategy string) bool {
	strategyFlag := "--squash"
	switch strategy {
	case "merge":
		strategyFlag = "--merge"
	case "rebase":
		strategyFlag = "--rebase"
	}

	const maxAttempts = 3
	var lastErr error
	var lastOutput string

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		_, r, err := gh.Exec("pr", "merge", pr.Number, "--repo", pr.Repo, strategyFlag)

		if err == nil {
			return true
		}

		lastErr = err
		lastOutput = r.String()

		if attempt == maxAttempts || !isBehindBase(pr) {
			break
		}

		fmt.Println("Pull Request is behind the base branch, updating and retrying - " + pr.GetUrl())

		if updateErr := updatePullRequestBranch(pr); updateErr != nil {
			fmt.Println("Failed to update branch for " + pr.GetUrl() + ": " + updateErr.Error())
			break
		}

		if !waitForMergeable(pr, 2*time.Minute) {
			fmt.Println("Timed out waiting for " + pr.GetUrl() + " to become mergeable after update")
			break
		}
	}

	fmt.Println("Failed to merge pr")
	fmt.Println("approimate cmd: gh pr merge " + pr.Number + " --repo " + pr.Repo + " " + strategyFlag)
	fmt.Println(lastOutput)
	fmt.Println(lastErr)

	return false
}
