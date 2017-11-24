package github

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/aurieh/ddg-ng/utils"
	"net/http"
)

type GithubRepo struct {
	ArchiveURL       string `json:"archive_url"`
	Archived         bool   `json:"archived"`
	AssigneesURL     string `json:"assignees_url"`
	BlobsURL         string `json:"blobs_url"`
	BranchesURL      string `json:"branches_url"`
	CloneURL         string `json:"clone_url"`
	CollaboratorsURL string `json:"collaborators_url"`
	CommentsURL      string `json:"comments_url"`
	CommitsURL       string `json:"commits_url"`
	CompareURL       string `json:"compare_url"`
	ContentsURL      string `json:"contents_url"`
	ContributorsURL  string `json:"contributors_url"`
	CreatedAt        string `json:"created_at"`
	DefaultBranch    string `json:"default_branch"`
	DeploymentsURL   string `json:"deployments_url"`
	Description      string `json:"description"`
	DownloadsURL     string `json:"downloads_url"`
	EventsURL        string `json:"events_url"`
	Fork             bool   `json:"fork"`
	Forks            int    `json:"forks"`
	ForksCount       int    `json:"forks_count"`
	ForksURL         string `json:"forks_url"`
	FullName         string `json:"full_name"`
	GitCommitsURL    string `json:"git_commits_url"`
	GitRefsURL       string `json:"git_refs_url"`
	GitTagsURL       string `json:"git_tags_url"`
	GitURL           string `json:"git_url"`
	HasDownloads     bool   `json:"has_downloads"`
	HasIssues        bool   `json:"has_issues"`
	HasPages         bool   `json:"has_pages"`
	HasProjects      bool   `json:"has_projects"`
	HasWiki          bool   `json:"has_wiki"`
	Homepage         string `json:"homepage"`
	HooksURL         string `json:"hooks_url"`
	HTMLURL          string `json:"html_url"`
	ID               int    `json:"id"`
	IssueCommentURL  string `json:"issue_comment_url"`
	IssueEventsURL   string `json:"issue_events_url"`
	IssuesURL        string `json:"issues_url"`
	KeysURL          string `json:"keys_url"`
	LabelsURL        string `json:"labels_url"`
	Language         string `json:"language"`
	LanguagesURL     string `json:"languages_url"`
	MergesURL        string `json:"merges_url"`
	MilestonesURL    string `json:"milestones_url"`
	MirrorURL        string `json:"mirror_url"`
	Name             string `json:"name"`
	NetworkCount     int    `json:"network_count"`
	NotificationsURL string `json:"notifications_url"`
	OpenIssues       int    `json:"open_issues"`
	OpenIssuesCount  int    `json:"open_issues_count"`
	Owner            `json:"owner"`
	Private          bool   `json:"private"`
	PullsURL         string `json:"pulls_url"`
	PushedAt         string `json:"pushed_at"`
	ReleasesURL      string `json:"releases_url"`
	Size             int    `json:"size"`
	SSHURL           string `json:"ssh_url"`
	StargazersCount  int    `json:"stargazers_count"`
	StargazersURL    string `json:"stargazers_url"`
	StatusesURL      string `json:"statuses_url"`
	SubscribersCount int    `json:"subscribers_count"`
	SubscribersURL   string `json:"subscribers_url"`
	SubscriptionURL  string `json:"subscription_url"`
	SvnURL           string `json:"svn_url"`
	TagsURL          string `json:"tags_url"`
	TeamsURL         string `json:"teams_url"`
	TreesURL         string `json:"trees_url"`
	UpdatedAt        string `json:"updated_at"`
	URL              string `json:"url"`
	Watchers         int    `json:"watchers"`
	WatchersCount    int    `json:"watchers_count"`
}

type Owner struct {
	AvatarURL         string `json:"avatar_url"`
	EventsURL         string `json:"events_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	GravatarID        string `json:"gravatar_id"`
	HTMLURL           string `json:"html_url"`
	ID                int    `json:"id"`
	Login             string `json:"login"`
	OrganizationsURL  string `json:"organizations_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	ReposURL          string `json:"repos_url"`
	SiteAdmin         bool   `json:"site_admin"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	Type              string `json:"type"`
	URL               string `json:"url"`
}

const GITHUB_REPO = "https://api.github.com/repos/%s/%s"

func Repo(client *http.Client, username string, reponame string) (*GithubRepo, error) {
	repo := &GithubRepo{}
	url := fmt.Sprintf(GITHUB_REPO, username, reponame)
	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == 404 {
		return nil, errors.New("repo not found")
	} else if res.StatusCode != 200 {
		log.Errorln("get github repo status: " + res.Status)
		return nil, errors.New("unknown error")
	}
	err = utils.GetJson(res, repo)
	if err != nil {
		return nil, err
	}
	return repo, nil
}