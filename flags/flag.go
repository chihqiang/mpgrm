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

	FlagsBranches = "branches"
	FlagsTags     = "tags"
	FlagsFiles    = "files"
)

func FormTargetRelease() []cli.Flag {
	var flag []cli.Flag
	flag = append(flag, FormTargetFlags()...)
	flag = append(flag, BranchFlags()...)
	flag = append(flag, TagsFlags()...)
	flag = append(flag, FilesFlags()...)
	return flag
}
func FormTargetRepo() []cli.Flag {
	var flag []cli.Flag
	flag = append(flag, FormTargetFlags()...)
	flag = append(flag, BranchFlags()...)
	flag = append(flag, TagsFlags()...)
	return flag
}
func FormTargetFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  FlagsFormRepo,
			Usage: "Source repository URL, the starting point of your magic",
		},
		&cli.StringFlag{
			Name:  FlagsTargetRepo,
			Usage: "Target repository URL, the final destination of your conquest",
		},
	}
}
func BranchFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringSliceFlag{
			Name:  FlagsBranches,
			Usage: "Branches at your cmd, releases under your control",
		},
	}
}

func TagsFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringSliceFlag{
			Name:  FlagsTags,
			Usage: "Bring your chosen tags to the battlefield, do as you please",
		},
	}
}

func FilesFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringSliceFlag{
			Name:  FlagsFiles,
			Usage: "Select your files freely, manage with ease",
		},
	}
}

func GetFormCredential(cmd *cli.Command, readEnv bool) (*url.URL, *credential.Credential, error) {
	repoURL, err := url.Parse(cmd.String(FlagsFormRepo))
	if err != nil {
		return nil, nil, err
	}
	cred, err := credential.GetCredential(repoURL, cmd.String(FlagsFormUsername), cmd.String(FlagsFormToken), readEnv)
	return repoURL, cred, err
}
func GetTargetCredential(cmd *cli.Command) (*url.URL, *credential.Credential, error) {
	repoURL, err := url.Parse(cmd.String(FlagsTargetRepo))
	if err != nil {
		return nil, nil, err
	}
	cred, err := credential.GetCredential(repoURL, cmd.String(FlagsTargetUsername), cmd.String(FlagsTargetToken), true)
	return repoURL, cred, err
}

func GetBranches(cmd *cli.Command) []string {
	return x.StringSplits(cmd.StringSlice(FlagsBranches), ",")
}

func GetTags(cmd *cli.Command) []string {
	return x.StringSplits(cmd.StringSlice(FlagsTags), ",")
}
func GetFirstTags(cmd *cli.Command) (string, error) {
	tags := GetTags(cmd)
	if len(tags) <= 0 {
		return "", fmt.Errorf("no tags received; cannot use the first tag by default")
	}
	return tags[0], nil
}

func GetFiles(cmd *cli.Command) []string {
	return x.MatchedFiles(x.StringSplits(cmd.StringSlice(FlagsFiles), ","))
}
