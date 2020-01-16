package external

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/go-github/github"
	"github.com/horalstvo/ghs/util"
	"golang.org/x/oauth2"
)

func GetClient(ctx context.Context, apiToken string) *github.Client {

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: apiToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func GetPullRequests(org string, repo string, ctx context.Context, client *github.Client) []*github.PullRequest {
	var pullRequests []*github.PullRequest
	var page = 1
	lastPage := "?"
	for {
		fmt.Printf("Fetching %s/%s page %d/%s ...\n", org, repo, page, lastPage)
		prs, response, err := client.PullRequests.List(ctx, org, repo, &github.PullRequestListOptions{
			Sort:      "created",
			State:     "all",
			Direction: "desc",
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: 100,
			},
		})

		util.Check(err)
		pullRequests = append(pullRequests, prs...)
		if response.NextPage > 0 {
			page++
			lastPage = fmt.Sprintf("%d", response.LastPage)
			time.Sleep(2 * time.Second)
			// break // TODO Remove
		} else {
			break
		}
	}

	return pullRequests
}

func GetReviews(org string, repo string, number int, ctx context.Context,
	client *github.Client) []*github.PullRequestReview {
	reviews, _, err := client.PullRequests.ListReviews(ctx, org, repo, number, &github.ListOptions{})
	util.Check(err)
	return reviews
}

func GetRepos(org string, repoNames []string, ctx context.Context, client *github.Client) []*github.Repository {
	var waiter sync.WaitGroup
	length := len(repoNames)
	waiter.Add(length)
	repositories := make([]*github.Repository, length, length)
	for index, repoName := range repoNames {
		go func(index int, repoName string) {
			defer waiter.Done()
			repository := GetRepo(org, repoName, ctx, client)
			repositories[index] = repository
		}(index, repoName)
	}
	waiter.Wait()
	return repositories
}

func GetRepo(org string, repoName string, ctx context.Context, client *github.Client) *github.Repository {
	repository, _, err := client.Repositories.Get(ctx, org, repoName)
	util.Check(err)

	return repository
}

func GetTeamRepos(org string, team string, ctx context.Context, client *github.Client) []*github.Repository {
	teamId, getTeamErr := getTeamId(org, team, ctx, client)
	util.Check(getTeamErr)

	repos, _, err := client.Teams.ListTeamRepos(ctx, *teamId, &github.ListOptions{})
	util.Check(err)
	return repos
}

func getTeamId(org string, team string, ctx context.Context, client *github.Client) (*int64, error) {
	teams, _, err := client.Teams.ListTeams(ctx, org, &github.ListOptions{})
	util.Check(err)

	for _, t := range teams {
		if *t.Name == team {
			return t.ID, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Team not found for %s", team))
}
