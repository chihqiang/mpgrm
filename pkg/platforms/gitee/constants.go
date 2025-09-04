package gitee

import "time"

type CreateRepoResponse struct {
	Assignee []struct {
		AvatarUrl         string `json:"avatar_url"`
		EventsUrl         string `json:"events_url"`
		FollowersUrl      string `json:"followers_url"`
		FollowingUrl      string `json:"following_url"`
		GistsUrl          string `json:"gists_url"`
		HtmlUrl           string `json:"html_url"`
		Id                int    `json:"id"`
		Login             string `json:"login"`
		Name              string `json:"name"`
		OrganizationsUrl  string `json:"organizations_url"`
		ReceivedEventsUrl string `json:"received_events_url"`
		Remark            string `json:"remark"`
		ReposUrl          string `json:"repos_url"`
		StarredUrl        string `json:"starred_url"`
		SubscriptionsUrl  string `json:"subscriptions_url"`
		Type              string `json:"type"`
		Url               string `json:"RepoURL"`
	} `json:"assignee"`
	AssigneesNumber int `json:"assignees_number"`
	Assigner        struct {
		AvatarUrl         string `json:"avatar_url"`
		EventsUrl         string `json:"events_url"`
		FollowersUrl      string `json:"followers_url"`
		FollowingUrl      string `json:"following_url"`
		GistsUrl          string `json:"gists_url"`
		HtmlUrl           string `json:"html_url"`
		Id                int    `json:"id"`
		Login             string `json:"login"`
		Name              string `json:"name"`
		OrganizationsUrl  string `json:"organizations_url"`
		ReceivedEventsUrl string `json:"received_events_url"`
		Remark            string `json:"remark"`
		ReposUrl          string `json:"repos_url"`
		StarredUrl        string `json:"starred_url"`
		SubscriptionsUrl  string `json:"subscriptions_url"`
		Type              string `json:"type"`
		Url               string `json:"RepoURL"`
	} `json:"assigner"`
	BlobsUrl            string      `json:"blobs_url"`
	BranchesUrl         string      `json:"branches_url"`
	CanComment          bool        `json:"can_comment"`
	CollaboratorsUrl    string      `json:"collaborators_url"`
	CommentsUrl         string      `json:"comments_url"`
	CommitsUrl          string      `json:"commits_url"`
	ContributorsUrl     string      `json:"contributors_url"`
	CreatedAt           time.Time   `json:"created_at"`
	DefaultBranch       interface{} `json:"default_branch"`
	Description         string      `json:"description"`
	Enterprise          interface{} `json:"enterprise"`
	Fork                bool        `json:"fork"`
	ForksCount          int         `json:"forks_count"`
	ForksUrl            string      `json:"forks_url"`
	FullName            string      `json:"full_name"`
	Gvp                 bool        `json:"gvp"`
	HasIssues           bool        `json:"has_issues"`
	HasPage             bool        `json:"has_page"`
	HasWiki             bool        `json:"has_wiki"`
	Homepage            interface{} `json:"homepage"`
	HooksUrl            string      `json:"hooks_url"`
	HtmlUrl             string      `json:"html_url"`
	HumanName           string      `json:"human_name"`
	Id                  int         `json:"id"`
	Internal            bool        `json:"internal"`
	IssueComment        interface{} `json:"issue_comment"`
	IssueCommentUrl     string      `json:"issue_comment_url"`
	IssueTemplateSource string      `json:"issue_template_source"`
	IssuesUrl           string      `json:"issues_url"`
	KeysUrl             string      `json:"keys_url"`
	LabelsUrl           string      `json:"labels_url"`
	Language            interface{} `json:"language"`
	License             interface{} `json:"license"`
	Members             []string    `json:"members"`
	MilestonesUrl       string      `json:"milestones_url"`
	Name                string      `json:"name"`
	Namespace           struct {
		HtmlUrl string `json:"html_url"`
		Id      int    `json:"id"`
		Name    string `json:"name"`
		Path    string `json:"path"`
		Type    string `json:"type"`
	} `json:"namespace"`
	NotificationsUrl string `json:"notifications_url"`
	OpenIssuesCount  int    `json:"open_issues_count"`
	Outsourced       bool   `json:"outsourced"`
	Owner            struct {
		AvatarUrl         string `json:"avatar_url"`
		EventsUrl         string `json:"events_url"`
		FollowersUrl      string `json:"followers_url"`
		FollowingUrl      string `json:"following_url"`
		GistsUrl          string `json:"gists_url"`
		HtmlUrl           string `json:"html_url"`
		Id                int    `json:"id"`
		Login             string `json:"login"`
		Name              string `json:"name"`
		OrganizationsUrl  string `json:"organizations_url"`
		ReceivedEventsUrl string `json:"received_events_url"`
		Remark            string `json:"remark"`
		ReposUrl          string `json:"repos_url"`
		StarredUrl        string `json:"starred_url"`
		SubscriptionsUrl  string `json:"subscriptions_url"`
		Type              string `json:"type"`
		Url               string `json:"RepoURL"`
	} `json:"owner"`
	Paas       interface{} `json:"paas"`
	Parent     interface{} `json:"parent"`
	Path       string      `json:"path"`
	Permission struct {
		Admin bool `json:"admin"`
		Pull  bool `json:"pull"`
		Push  bool `json:"push"`
	} `json:"permission"`
	Private             bool          `json:"private"`
	Programs            []interface{} `json:"programs"`
	ProjectCreator      string        `json:"project_creator"`
	ProjectLabels       []interface{} `json:"project_labels"`
	Public              bool          `json:"public"`
	PullRequestsEnabled bool          `json:"pull_requests_enabled"`
	PullsUrl            string        `json:"pulls_url"`
	PushedAt            interface{}   `json:"pushed_at"`
	Recommend           bool          `json:"recommend"`
	Relation            string        `json:"relation"`
	ReleasesUrl         string        `json:"releases_url"`
	SshUrl              string        `json:"ssh_url"`
	Stared              bool          `json:"stared"`
	StargazersCount     int           `json:"stargazers_count"`
	StargazersUrl       string        `json:"stargazers_url"`
	Status              string        `json:"status"`
	TagsUrl             string        `json:"tags_url"`
	Testers             []struct {
		AvatarUrl         string `json:"avatar_url"`
		EventsUrl         string `json:"events_url"`
		FollowersUrl      string `json:"followers_url"`
		FollowingUrl      string `json:"following_url"`
		GistsUrl          string `json:"gists_url"`
		HtmlUrl           string `json:"html_url"`
		Id                int    `json:"id"`
		Login             string `json:"login"`
		Name              string `json:"name"`
		OrganizationsUrl  string `json:"organizations_url"`
		ReceivedEventsUrl string `json:"received_events_url"`
		Remark            string `json:"remark"`
		ReposUrl          string `json:"repos_url"`
		StarredUrl        string `json:"starred_url"`
		SubscriptionsUrl  string `json:"subscriptions_url"`
		Type              string `json:"type"`
		Url               string `json:"RepoURL"`
	} `json:"testers"`
	TestersNumber int       `json:"testers_number"`
	UpdatedAt     time.Time `json:"updated_at"`
	Url           string    `json:"RepoURL"`
	Watched       bool      `json:"watched"`
	WatchersCount int       `json:"watchers_count"`
}
type RepoInfoResponse struct {
	Assignee []struct {
		AvatarUrl         string `json:"avatar_url"`
		EventsUrl         string `json:"events_url"`
		FollowersUrl      string `json:"followers_url"`
		FollowingUrl      string `json:"following_url"`
		GistsUrl          string `json:"gists_url"`
		HtmlUrl           string `json:"html_url"`
		Id                int    `json:"id"`
		Login             string `json:"login"`
		Name              string `json:"name"`
		OrganizationsUrl  string `json:"organizations_url"`
		ReceivedEventsUrl string `json:"received_events_url"`
		Remark            string `json:"remark"`
		ReposUrl          string `json:"repos_url"`
		StarredUrl        string `json:"starred_url"`
		SubscriptionsUrl  string `json:"subscriptions_url"`
		Type              string `json:"type"`
		Url               string `json:"RepoURL"`
	} `json:"assignee"`
	AssigneesNumber int `json:"assignees_number"`
	Assigner        struct {
		AvatarUrl         string `json:"avatar_url"`
		EventsUrl         string `json:"events_url"`
		FollowersUrl      string `json:"followers_url"`
		FollowingUrl      string `json:"following_url"`
		GistsUrl          string `json:"gists_url"`
		HtmlUrl           string `json:"html_url"`
		Id                int    `json:"id"`
		Login             string `json:"login"`
		Name              string `json:"name"`
		OrganizationsUrl  string `json:"organizations_url"`
		ReceivedEventsUrl string `json:"received_events_url"`
		Remark            string `json:"remark"`
		ReposUrl          string `json:"repos_url"`
		StarredUrl        string `json:"starred_url"`
		SubscriptionsUrl  string `json:"subscriptions_url"`
		Type              string `json:"type"`
		Url               string `json:"RepoURL"`
	} `json:"assigner"`
	BlobsUrl            string      `json:"blobs_url"`
	BranchesUrl         string      `json:"branches_url"`
	CanComment          bool        `json:"can_comment"`
	CollaboratorsUrl    string      `json:"collaborators_url"`
	CommentsUrl         string      `json:"comments_url"`
	CommitsUrl          string      `json:"commits_url"`
	ContributorsUrl     string      `json:"contributors_url"`
	CreatedAt           time.Time   `json:"created_at"`
	DefaultBranch       string      `json:"default_branch"`
	Description         string      `json:"description"`
	Enterprise          interface{} `json:"enterprise"`
	Fork                bool        `json:"fork"`
	ForksCount          int         `json:"forks_count"`
	ForksUrl            string      `json:"forks_url"`
	FullName            string      `json:"full_name"`
	Gvp                 bool        `json:"gvp"`
	HasIssues           bool        `json:"has_issues"`
	HasPage             bool        `json:"has_page"`
	HasWiki             bool        `json:"has_wiki"`
	Homepage            string      `json:"homepage"`
	HooksUrl            string      `json:"hooks_url"`
	HtmlUrl             string      `json:"html_url"`
	HumanName           string      `json:"human_name"`
	Id                  int64       `json:"id"`
	Internal            bool        `json:"internal"`
	IssueComment        bool        `json:"issue_comment"`
	IssueCommentUrl     string      `json:"issue_comment_url"`
	IssueTemplateSource string      `json:"issue_template_source"`
	IssuesUrl           string      `json:"issues_url"`
	KeysUrl             string      `json:"keys_url"`
	LabelsUrl           string      `json:"labels_url"`
	Language            string      `json:"language"`
	License             string      `json:"license"`
	Members             []string    `json:"members"`
	MilestonesUrl       string      `json:"milestones_url"`
	Name                string      `json:"name"`
	Namespace           struct {
		HtmlUrl string `json:"html_url"`
		Id      int    `json:"id"`
		Name    string `json:"name"`
		Path    string `json:"path"`
		Type    string `json:"type"`
	} `json:"namespace"`
	NotificationsUrl string `json:"notifications_url"`
	OpenIssuesCount  int    `json:"open_issues_count"`
	Outsourced       bool   `json:"outsourced"`
	Owner            struct {
		AvatarUrl         string `json:"avatar_url"`
		EventsUrl         string `json:"events_url"`
		FollowersUrl      string `json:"followers_url"`
		FollowingUrl      string `json:"following_url"`
		GistsUrl          string `json:"gists_url"`
		HtmlUrl           string `json:"html_url"`
		Id                int    `json:"id"`
		Login             string `json:"login"`
		Name              string `json:"name"`
		OrganizationsUrl  string `json:"organizations_url"`
		ReceivedEventsUrl string `json:"received_events_url"`
		Remark            string `json:"remark"`
		ReposUrl          string `json:"repos_url"`
		StarredUrl        string `json:"starred_url"`
		SubscriptionsUrl  string `json:"subscriptions_url"`
		Type              string `json:"type"`
		Url               string `json:"RepoURL"`
	} `json:"owner"`
	Paas       interface{} `json:"paas"`
	Parent     interface{} `json:"parent"`
	Path       string      `json:"path"`
	Permission struct {
		Admin bool `json:"admin"`
		Pull  bool `json:"pull"`
		Push  bool `json:"push"`
	} `json:"permission"`
	Private             bool          `json:"private"`
	Programs            []interface{} `json:"programs"`
	ProjectCreator      string        `json:"project_creator"`
	ProjectLabels       []interface{} `json:"project_labels"`
	Public              bool          `json:"public"`
	PullRequestsEnabled bool          `json:"pull_requests_enabled"`
	PullsUrl            string        `json:"pulls_url"`
	PushedAt            time.Time     `json:"pushed_at"`
	Recommend           bool          `json:"recommend"`
	Relation            string        `json:"relation"`
	ReleasesUrl         string        `json:"releases_url"`
	SshUrl              string        `json:"ssh_url"`
	Stared              bool          `json:"stared"`
	StargazersCount     int           `json:"stargazers_count"`
	StargazersUrl       string        `json:"stargazers_url"`
	Status              string        `json:"status"`
	TagsUrl             string        `json:"tags_url"`
	Testers             []struct {
		AvatarUrl         string `json:"avatar_url"`
		EventsUrl         string `json:"events_url"`
		FollowersUrl      string `json:"followers_url"`
		FollowingUrl      string `json:"following_url"`
		GistsUrl          string `json:"gists_url"`
		HtmlUrl           string `json:"html_url"`
		Id                int    `json:"id"`
		Login             string `json:"login"`
		Name              string `json:"name"`
		OrganizationsUrl  string `json:"organizations_url"`
		ReceivedEventsUrl string `json:"received_events_url"`
		Remark            string `json:"remark"`
		ReposUrl          string `json:"repos_url"`
		StarredUrl        string `json:"starred_url"`
		SubscriptionsUrl  string `json:"subscriptions_url"`
		Type              string `json:"type"`
		Url               string `json:"RepoURL"`
	} `json:"testers"`
	TestersNumber int       `json:"testers_number"`
	UpdatedAt     time.Time `json:"updated_at"`
	Url           string    `json:"RepoURL"`
	Watched       bool      `json:"watched"`
	WatchersCount int       `json:"watchers_count"`
}

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

type UploadAttachResponse struct {
	ID       int    `json:"id"`
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

type TagResponse struct {
	Name    string `json:"name"`
	Message string `json:"message,omitempty"`
	Commit  struct {
		Sha  string    `json:"sha"`
		Date time.Time `json:"date"`
	} `json:"commit"`
	Tagger struct {
		Name  string    `json:"name"`
		Email string    `json:"email"`
		Date  time.Time `json:"date"`
	} `json:"tagger"`
}

type ReleaseResponse struct {
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
