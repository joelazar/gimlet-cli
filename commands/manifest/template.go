package manifest

import (
	"fmt"
	"github.com/gimlet-io/gimletd/dx"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"

	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
)

var manifestTemplateCmd = cli.Command{
	Name:  "template",
	Usage: "Templates a Gimlet manifest",
	UsageText: `gimlet manifest template \
    -f .gimlet/staging.yaml \
    -o manifests.yaml \
    --vars ci.env`,
	Action: templateCmd,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "file",
			Aliases:  []string{"f"},
			Required: true,
			Usage:    "Gimlet manifest file to template, or \"-\" for stdin",
		},
		&cli.StringFlag{
			Name:    "vars",
			Aliases: []string{"v"},
			Usage:   "an .env file for template variables",
		},
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "output file",
		},
	},
}

func templateCmd(c *cli.Context) error {
	varsPath := c.String("vars")
	vars := map[string]string{}
	if varsPath != "" {
		yamlString, err := ioutil.ReadFile(varsPath)
		if err != nil {
			return fmt.Errorf("cannot read vars file")
		}

		vars, err = godotenv.Parse(strings.NewReader(string(yamlString)))
		if err != nil {
			return fmt.Errorf("cannot parse vars")
		}
	}

	for _, v := range os.Environ() {
		pair := strings.SplitN(v, "=", 2)
		if _, exists := vars[pair[0]]; !exists {
			vars[pair[0]] = pair[1]
		}
	}

	manifestPath := c.String("file")
	manifestString, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("cannot read manifest file")
	}

	var m dx.Manifest
	err = yaml.Unmarshal(manifestString, &m)
	if err != nil {
		return fmt.Errorf("cannot unmarshal manifest")
	}

	err = m.ResolveVars(vars)
	if err != nil {
		return fmt.Errorf("cannot resolve manifest vars %s", err.Error())
	}

	if strings.HasPrefix(m.Chart.Name, "git@") {
		tmpChartDir, err := dx.CloneChartFromRepo(m, "")
		if err != nil {
			return fmt.Errorf("cannot fetch chart from git %s", err.Error())
		}
		m.Chart.Name = tmpChartDir
		defer os.RemoveAll(tmpChartDir)
	}

	templatesManifests, err := dx.HelmTemplate(m)
	if err != nil {
		return fmt.Errorf("cannot template Helm chart %s", err)
	}

	outputPath := c.String("output")
	if outputPath != "" {
		err := ioutil.WriteFile(outputPath, []byte(templatesManifests), 0666)
		if err != nil {
			return fmt.Errorf("cannot write values file %s", err)
		}
	} else {
		fmt.Println(templatesManifests)
	}

	return nil
}
