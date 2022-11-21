package main

import (
	"github.com/andrewhertog/helm-chart-checker/cmd"
)

func main() {
	cmd.Execute()
	// charts := []helm.Chart{}
	// charts := helm.Charts{}
	// resource := kubernetes.Resource{}
	// resource.Resource("helm.fluxcd.io", "v1", "helmreleases")
	// namespace := "system-releases"
	// resource.GetHelmReleaseForNamespace(namespace)

	// for _, item := range resource.Content {
	// 	chart := helm.GetChartSpec(item.Object)
	// 	if !funk.Contains(charts.Charts, *chart) {
	// 		charts.Append(chart)
	// 	}
	// }
	// for i, chart := range charts.Charts {
	// 	if chart.IsHelmRepo() {
	// 		chart.GetLatestChartVersion()
	// 		charts.Charts[i] = chart
	// 	}

	// }
	// charts.PrintChartData()

}
