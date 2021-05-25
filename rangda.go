package rangda

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"
)

// ######################################################################
// Rangda (controller/service)

// Handles pull request review event. Ensures PR is mergeable and performs merge
type Rangda struct {
	client *github.Client
}

// Instances rangda service with client from context
func NewRangda(apiKey string) *Rangda {
	r := Rangda{}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: apiKey},
	)
	tc := oauth2.NewClient(ctx, ts)

	r.client = github.NewClient(tc)

	return &r
}

// Function called to interpret Github Pull Request Event web hook
func (rangda *Rangda) ReviewEventHandler(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err)
	}

	prReview := github.PullRequestReviewEvent{}
	json.Unmarshal(reqBody, &prReview)

	if strings.ToUpper(*prReview.Review.State) != "APPROVED" {
		fmt.Println("Pull Request review is not an approval.")
		w.Write([]byte("Pull Request review is not an approval."))
		return
	}

	if strings.ToUpper(*prReview.PullRequest.MergeableState) != "CLEAN" {
		fmt.Println("Pull Request cannot be automatically merged.")
		rangda.comment(*prReview.Repo, *prReview.PullRequest, "Pull Request cannot be merged. Please fix Pull Request and re approve.")
		w.Write([]byte("Pull Request is not suitable for merging."))
		return
	}

	rangda.merge(*prReview.Repo, *prReview.PullRequest)
	w.Write([]byte("Pull Request merge was attempted."))
}

// Creates a comment on the passed pull request
func (rangda *Rangda) comment(repo github.Repository, pr github.PullRequest, commentText string) {
	c := github.PullRequestComment{
		Body: &commentText,
	}

	_, res, err := rangda.client.PullRequests.CreateComment(
		context.Background(),
		*repo.Owner.Name,
		*repo.Name,
		int(*pr.ID),
		&c,
	)

	if err != nil {
		panic(err)
	}

	if res.Status != http.StatusText(http.StatusCreated) {
		fmt.Println("Failed to make comment.")
	}
}

// Merges the target Pull Request
func (rangda *Rangda) merge(repo github.Repository, pr github.PullRequest) {
	_, res, err := rangda.client.PullRequests.Merge(
		context.Background(),
		*repo.Owner.Name,
		*repo.Name,
		int(*pr.ID),
		"",
		&github.PullRequestOptions{},
	)

	if err != nil {
		panic(err)
	}

	if res.Status != http.StatusText((http.StatusOK)) {
		fmt.Println("Failed to merge pull request.")
	}
}

// ######################################################################
// Secrets (configuration)

// Model capturing secrets configuration for rangda server
type Secrets struct {
	Port   uint   `json:"port"`
	Host   string `json:"host"`
	ApiKey string `json:"api_key"`
}

// Function to retrieve secrets from secrets.json file in project
func GetSecrets(path string) (*Secrets, error) {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var s Secrets
	json.Unmarshal(data, &s)

	return &s, nil
}
