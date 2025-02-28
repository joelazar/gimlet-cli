package artifact

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gimlet-io/gimletd/dx"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"time"
)

var artifactCreateCmd = cli.Command{
	Name:  "create",
	Usage: "Creates a release artifact",
	UsageText: `gimlet artifact create \
     --repository=my-app \
     --sha=26fc62ffa5cf63204ccbce6876c6d610 \
     --branch=master \
     --authorName=Laszlo \
     --authorEmail=laszlo@laszlo.laszlo \
     --committerName=Laszlo \
     --committerEmail=laszlo@laszlo.laszlo \
     --message="Bugfix 123" \
     --url="https://github.com/owner/repo/commits/0017d995e32e3d1998395d971b969bcf682d2085" \
     > artifact.json`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "repository",
			Usage:    "The git repository name (mandatory)",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "sha",
			Usage:    "The git sha (mandatory)",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "branch",
			Usage:    "The git branch, or target branch for pull request builds (mandatory)",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "event",
			Usage: "If this is a push/tag/pr build",
		},
		&cli.StringFlag{
			Name:  "sourceBranch",
			Usage: "For pull requests, the feature branch name",
		},
		&cli.StringFlag{
			Name:  "targetBranch",
			Usage: "For pull requests, the merge target",
		},
		&cli.StringFlag{
			Name:  "tag",
			Usage: "For tag builds, the name of the tag",
		},
		&cli.StringFlag{
			Name:     "authorName",
			Usage:    "The person who originally wrote the code (mandatory)",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "authorEmail",
			Usage:    "The person who originally wrote the code (mandatory)",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "committerName",
			Usage:    "The person who originally wrote the code (mandatory)",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "committerEmail",
			Usage:    "The committer is the person who committed the code. Important in case of history rewrite (mandatory)",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "message",
			Usage:    "The git commit message (mandatory)",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "url",
			Usage:    "URL to the git commit (mandatory)",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "output manifest file",
		},
	},
	Action: create,
}

func create(c *cli.Context) error {
	eventStr := c.String("event")
	event := dx.Push
	err := event.UnmarshalJSON([]byte(`"` + eventStr + `"`))
	if err != nil {
		return fmt.Errorf("cannot parse event: %s", err)
	}

	artifact := &dx.Artifact{
		Version: dx.Version{
			RepositoryName: c.String("repository"),
			SHA:            c.String("sha"),
			Branch:         c.String("branch"),
			Event:          event,
			SourceBranch:   c.String("sourceBranch"),
			TargetBranch:   c.String("targetBranch"),
			Tag:            c.String("tag"),
			AuthorName:     c.String("authorName"),
			AuthorEmail:    c.String("authorEmail"),
			CommitterName:  c.String("committerName"),
			CommitterEmail: c.String("committerEmail"),
			Message:        c.String("message"),
			URL:            c.String("url"),
		},
		Context: map[string]string{

		},
		Environments: []*dx.Manifest{

		},
		Items: []map[string]interface{}{

		},
		Created: time.Now().Unix(),
	}

	jsonString := bytes.NewBufferString("")
	e := json.NewEncoder(jsonString)
	e.SetIndent("", "  ")
	e.Encode(artifact)

	outputPath := c.String("output")
	if outputPath != "" {
		err := ioutil.WriteFile(outputPath, jsonString.Bytes(), 0666)
		if err != nil {
			return fmt.Errorf("cannot write artifact json %s", err)
		}
	} else {
		fmt.Println(jsonString.String())
	}

	return nil
}
