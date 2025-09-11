package flags

import (
	"fmt"
	"github.com/urfave/cli/v3"
	"net/url"
	"wangzhiqiang/mpgrm/pkg/credential"
	"wangzhiqiang/mpgrm/pkg/x"
)

const (
	FlagsFormRepo   = "repo"
	FlagsTargetRepo = "target-repo"
	FlagsUseEnv     = "use-env"

	FlagsBranches = "branches"
	FlagsTags     = "tags"
	FlagsFiles    = "files"
)

func FormReleaseUploadFiles() []cli.Flag {
	var flag []cli.Flag
	flag = append(flag, FormFlags()...)
	flag = append(flag, TagsFlags()...)  // tag selection
	flag = append(flag, FilesFlags()...) // files selection
	return flag
}
func FormReleaseDownload() []cli.Flag {
	var flag []cli.Flag
	flag = append(flag, FormFlags()...)
	flag = append(flag, TagsFlags()...) // tag selection
	return flag
}

func FormTargetReleaseSync() []cli.Flag {
	var flag []cli.Flag
	flag = append(flag, FormFlags()...)
	flag = append(flag, TargetFlags()...)
	flag = append(flag, TagsFlags()...)
	return flag
}

// FormTargetRepo combines flags needed for repository operations.
func FormTargetRepo() []cli.Flag {
	var flag []cli.Flag
	flag = append(flag, FormFlags()...)
	flag = append(flag, TargetFlags()...)
	return flag
}

func FormTargetRepoPush() []cli.Flag {
	var flag []cli.Flag
	flag = append(flag, FormFlags()...)
	flag = append(flag, TargetFlags()...)
	flag = append(flag, BranchFlags()...)
	flag = append(flag, TagsFlags()...)
	return flag
}

func CredentialFlags() []cli.Flag {
	var flag []cli.Flag
	flag = append(flag, FormFlags()...)
	flag = append(flag, UseEnvFlags()...)
	return flag
}

func FormFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     FlagsFormRepo,
			Usage:    "Source repository URL, the starting point of your magic",
			Required: true,
		},
	}
}

// GetFormCredential retrieves the repository URL and credential from the command.
// If readEnv is true, it will also try to read credentials from environment variables.
// Returns:
//   - *url.URL: parsed repository URL
//   - *credential.Credential: credential object
//   - error: any parsing or credential error
func GetFormCredential(cmd *cli.Command, readEnv bool) (*url.URL, *credential.Credential, error) {
	repoURL, err := url.Parse(cmd.String(FlagsFormRepo))
	if err != nil {
		return nil, nil, err
	}
	cred, err := credential.GetCredential(repoURL, cmd.String(FlagsFormUsername), cmd.String(FlagsFormToken), readEnv)
	return repoURL, cred, err
}

// BranchFlags returns branch selection flag.
func BranchFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringSliceFlag{
			Name:  FlagsBranches,
			Usage: "Branches at your command, releases under your control",
		},
	}
}

func TargetFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     FlagsTargetRepo,
			Usage:    "Target repository URL, the final destination of your conquest",
			Required: true,
		},
	}
}

// GetTargetCredential retrieves the target repository URL and credential from the command.
// Always reads credentials from environment variables.
// Returns:
//   - *url.URL: parsed target repository URL
//   - *credential.Credential: credential object
//   - error: any parsing or credential error
func GetTargetCredential(cmd *cli.Command) (*url.URL, *credential.Credential, error) {
	repoURL, err := url.Parse(cmd.String(FlagsTargetRepo))
	if err != nil {
		return nil, nil, err
	}
	cred, err := credential.GetCredential(repoURL, cmd.String(FlagsTargetUsername), cmd.String(FlagsTargetToken), true)
	return repoURL, cred, err
}

// GetBranches returns a slice of branch names specified in the command flags.
// Branch names are split by comma.
func GetBranches(cmd *cli.Command) []string {
	return x.StringSplits(cmd.StringSlice(FlagsBranches), ",")
}

// TagsFlags returns tag selection flag.
func TagsFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringSliceFlag{
			Name:  FlagsTags,
			Usage: "Bring your chosen tags to the battlefield, do as you please",
		},
	}
}

// GetTags returns a slice of tag names specified in the command flags.
// Tag names are split by comma.
func GetTags(cmd *cli.Command) []string {
	return x.StringSplits(cmd.StringSlice(FlagsTags), ",")
}

// GetFirstTags returns the first tag from the command flags.
// Returns an error if no tags are provided.
func GetFirstTags(cmd *cli.Command) (string, error) {
	tags := GetTags(cmd)
	if len(tags) <= 0 {
		return "", fmt.Errorf("no tags received; cannot use the first tag by default")
	}
	return tags[0], nil
}

// FilesFlags returns file selection flag.
func FilesFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringSliceFlag{
			Name:  FlagsFiles,
			Usage: "Select your files freely, manage with ease",
		},
	}
}

// GetFiles returns a slice of file paths specified in the command flags.
// Supports comma-separated input and file pattern matching.
func GetFiles(cmd *cli.Command) []string {
	return x.MatchedFiles(x.StringSplits(cmd.StringSlice(FlagsFiles), ","))
}

func UseEnvFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:  FlagsUseEnv,
			Usage: "Typing is boring? Let env vars handle it",
			Value: true,
		},
	}
}

func UseEnv(cmd *cli.Command) bool {
	return cmd.Bool(FlagsUseEnv)
}
