package main

import (
	"embed"
	"flag"
	"fmt"

	"gopkg.in/yaml.v3"

	"emballm/cli"
)

//go:embed release.yaml
var content embed.FS

type Release struct {
	Version      string   `yaml:"version"`
	Copyright    string   `yaml:"copyright"`
	Contributors []string `yaml:"contributors"`
}

func (r Release) String() (release string) {
	release = fmt.Sprintf("emballm version %s\nCopyright (c) %s", r.Version, r.Copyright)
	release += "\n\nContributors:\n"
	for _, contributor := range r.Contributors {
		release += fmt.Sprintf("- %s\n", contributor)
	}
	return
}

func main() {
	data, _ := content.ReadFile("release.yaml")
	var release Release
	_ = yaml.Unmarshal(data, &release)

	flag.Parse()
	cli.Command(release.String())
}
