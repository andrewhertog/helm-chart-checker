package helm

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/table"
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
		chart.name = specChart["name"].(string)
		chart.repo = specChart["repository"].(string)
		chart.version = specChart["version"].(string)
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
	chart.latestVersion = cv.Version
}

func (charts *Charts) Append(chart *Chart) {
	charts.Charts = append(charts.Charts, *chart)
}

func (charts *Charts) PrintChartData() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "URL", "Version", "Latest"})
	for _, chart := range charts.Charts {
		if chart.name != "" {
			t.AppendRow([]interface{}{chart.name, chart.repo, chart.version, chart.latestVersion})
		}
	}
	t.Render()
}
