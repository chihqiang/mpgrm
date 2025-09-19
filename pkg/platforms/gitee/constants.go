package gitee

import "time"

type CreateReleaseResponse struct {
	ID              int64  `json:"id"`
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
	Prerelease      bool   `json:"prerelease"`
	Name            string `json:"name"`
	Body            string `json:"body"`
	Author          struct {
		Id                int    `json:"id"`
		Login             string `json:"login"`
		Name              string `json:"name"`
		AvatarUrl         string `json:"avatar_url"`
		Url               string `json:"RepoURL"`
		HtmlUrl           string `json:"html_url"`
		Remark            string `json:"remark"`
		FollowersUrl      string `json:"followers_url"`
		FollowingUrl      string `json:"following_url"`
		GistsUrl          string `json:"gists_url"`
		StarredUrl        string `json:"starred_url"`
		SubscriptionsUrl  string `json:"subscriptions_url"`
		OrganizationsUrl  string `json:"organizations_url"`
		ReposUrl          string `json:"repos_url"`
		EventsUrl         string `json:"events_url"`
		ReceivedEventsUrl string `json:"received_events_url"`
		Type              string `json:"type"`
	} `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	Assets    []struct {
		BrowserDownloadUrl string `json:"browser_download_url"`
		Name               string `json:"name"`
	} `json:"assets"`
}

type ListReleaseAssetResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Size     int    `json:"size"`
	Uploader struct {
		Id        int    `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		AvatarUrl string `json:"avatar_url"`
		Url       string `json:"RepoURL"`
		HtmlUrl   string `json:"html_url"`
		Remark    string `json:"remark"`
	} `json:"uploader"`
	BrowserDownloadUrl string `json:"browser_download_url"`
}

type ReleasesTagResponse struct {
	ID              int64  `json:"id"`
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
	Prerelease      bool   `json:"prerelease"`
	Name            string `json:"name"`
	Body            string `json:"body"`
	Author          struct {
		Id                int    `json:"id"`
		Login             string `json:"login"`
		Name              string `json:"name"`
		AvatarUrl         string `json:"avatar_url"`
		Url               string `json:"RepoURL"`
		HtmlUrl           string `json:"html_url"`
		Remark            string `json:"remark"`
		FollowersUrl      string `json:"followers_url"`
		FollowingUrl      string `json:"following_url"`
		GistsUrl          string `json:"gists_url"`
		StarredUrl        string `json:"starred_url"`
		SubscriptionsUrl  string `json:"subscriptions_url"`
		OrganizationsUrl  string `json:"organizations_url"`
		ReposUrl          string `json:"repos_url"`
		EventsUrl         string `json:"events_url"`
		ReceivedEventsUrl string `json:"received_events_url"`
		Type              string `json:"type"`
	} `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	Assets    []struct {
		BrowserDownloadUrl string `json:"browser_download_url"`
		Name               string `json:"name"`
	} `json:"assets"`
}
