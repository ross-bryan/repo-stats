// CheckArgs parses the cli input before we run our program
package stats

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

/*
Package stats will hold our data structures and
do the work to populate activity for a given repository
*/

// UserStats represents a users participation
type UserStats struct {
	Approvals        int
	Comments         int
	ChangesRequested int
	PullList         []int
	PullsOpened      int
	PullsMerged      int
	Username         string
}

// UserStatsOptions represents configuration options for our queries to generate the stats
type UserStatsOptions struct {
	AfterDate   time.Time
	Owner       string
	Repo        string
	ListOptions github.PullRequestListOptions
}

// StatsManager represents our manager struct containing access to required clients
// and data structures holding our participant data
type StatsManager struct {
	ghcli            *github.Client
	options          *UserStatsOptions
	participantStats map[int64]*UserStats
}

func CheckArgs(args []string) (*UserStatsOptions, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("invalid number of args, expected: 1, got: %v", len(args))
	}
	d, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, fmt.Errorf("couldn't parse parameter to int Days: %v", err)
	}
	var options UserStatsOptions
	options.Owner = args[0]
	options.Repo = args[1]
	options.AfterDate = time.Now().AddDate(0, 0, -d)
	options.ListOptions = github.PullRequestListOptions{
		State: "all",
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 30,
		},
	}
	return &options, nil
}

// Usage gives examples of how to run the program
func Usage() {
	fmt.Println("Example usage:")
	fmt.Println("./repo-stats owner repo days")
	fmt.Println("./repo-stats azure aro-rp 90")
}

func newManager(ctx context.Context, o *UserStatsOptions) (*StatsManager, error) {
	var m = new(StatsManager)

	// If our GH_PAT is not defined, use nil httpclient
	// Otherwise, create a static token resources and oauth client
	// to pass into our github client
	pat := os.Getenv("GH_PAT")
	if pat == "" {
		m.ghcli = github.NewClient(nil)
	} else {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: pat},
		)
		tc := oauth2.NewClient(ctx, ts)
		m.ghcli = github.NewClient(tc)
	}
	m.participantStats = make(map[int64]*UserStats)
	m.options = o

	return m, nil
}
