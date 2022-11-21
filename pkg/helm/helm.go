package helm

import (
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/table"
	"golang.org/x/mod/semver"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
)

type Chart struct {
	name          string
	repo          string
	version       string
	git           string
	ref           string
	path          string
	latestVersion string
}

type Charts struct {
	Charts []Chart
}

func GetChartSpec(item map[string]interface{}) *Chart {
	chart := Chart{}
	spec := item["spec"].(map[string]interface{})
	specChart := spec["chart"].(map[string]interface{})
	if specChart["name"] != nil {
		if strings.HasPrefix(specChart["version"].(string), "v") {
			chart.version = specChart["version"].(string)
		} else {
			chart.version = fmt.Sprintf("v%s", specChart["version"].(string))
		}

		chart.name = specChart["name"].(string)
		chart.repo = specChart["repository"].(string)
	} else {
		chart.ref = specChart["ref"].(string)
		chart.git = specChart["git"].(string)
		chart.path = specChart["path"].(string)
	}
	return &chart
}

func (chart *Chart) IsHelmRepo() bool {
	return chart.name != ""
}

func (chart *Chart) GetLatestChartVersion() {
	c, err := repo.NewChartRepository(&repo.Entry{Name: "tmp", URL: chart.repo}, getter.All(&cli.EnvSettings{}))
	if err != nil {
		fmt.Println(err)
	}
	i, err := c.DownloadIndexFile()
	if err != nil {
		fmt.Println(err)
	}
	indexFile, err := repo.LoadIndexFile(i)
	if err != nil {
		fmt.Println(err)
	}
	cv, err := indexFile.Get(chart.name, "")
	if err != nil {
		fmt.Println(err)
	}
	if strings.HasPrefix(cv.Version, "v") {
		chart.latestVersion = cv.Version
	} else {
		chart.latestVersion = fmt.Sprintf("v%s", cv.Version)
	}
}

func (chart *Chart) compareVersion() bool {
	return semver.Compare(chart.version, chart.latestVersion) != 0
}

func (charts *Charts) Append(chart *Chart) {
	charts.Charts = append(charts.Charts, *chart)
}

func (charts *Charts) PrintChartData(ood bool) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "URL", "Version", "Latest"})
	for _, chart := range charts.Charts {
		if chart.IsHelmRepo() {
			if chart.compareVersion() {
				chart.latestVersion = fmt.Sprintf("\033[31m%s\033[0m", chart.latestVersion)
				t.AppendRow([]interface{}{chart.name, chart.repo, chart.version, chart.latestVersion})
			} else {
				if !ood {
					t.AppendRow([]interface{}{chart.name, chart.repo, chart.version, chart.latestVersion})
				}
			}

		}
	}
	t.Render()
}
