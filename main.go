package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/thoas/go-funk"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	ctrl "sigs.k8s.io/controller-runtime"
)

type chart struct {
	name          string
	repo          string
	version       string
	git           string
	ref           string
	path          string
	latestVersion string
}

func main() {
	ctx := context.Background()
	config := ctrl.GetConfigOrDie()
	dynamic := dynamic.NewForConfigOrDie(config)
	charts := []chart{}
	namespace := "system-releases"
	items, err := GetResourcesDynamically(dynamic, ctx,
		"helm.fluxcd.io", "v1", "helmreleases", namespace)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, item := range items {
			chart := getChartSpec(item.Object)
			if !funk.Contains(charts, *chart) {
				charts = append(charts, *chart)
			}
		}
		for _, chart := range charts {
			if chart.name == "" {
				continue
			}
			chart.getLatestChartVersion()
			// c, err := repo.FindChartInRepoURL(chart.repo, chart.name, "", "", "", "", getter.All(&cli.EnvSettings{}))
			// if err != nil {
			// 	fmt.Println(err)
			// }
			// fmt.Println(c)

		}
		printChartData(charts)
	}

}

func printChartData(charts []chart) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "URL", "Version", "Latest"})
	for _, chart := range charts {
		if chart.name != "" {
			t.AppendRow([]interface{}{chart.name, chart.repo, chart.version, chart.latestVersion})
		} else {
			t.AppendRow([]interface{}{"", chart.git, "", ""})
		}
	}
	t.Render()
}

func (chart *chart) getLatestChartVersion() {
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
func getChartSpec(item map[string]interface{}) *chart {
	chart := chart{}
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

func GetResourcesDynamically(dynamic dynamic.Interface, ctx context.Context,
	group string, version string, resource string, namespace string) (
	[]unstructured.Unstructured, error) {

	resourceId := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}
	list, err := dynamic.Resource(resourceId).Namespace(namespace).
		List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return list.Items, nil
}
