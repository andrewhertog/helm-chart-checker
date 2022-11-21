package cmd

import (
	"github.com/andrewhertog/helm-chart-checker/pkg/helm"
	"github.com/andrewhertog/helm-chart-checker/pkg/kubernetes"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
)

func init() {
	rootCmd.PersistentFlags().StringP("namespace", "n", "default", "Namespace to look")
}

var rootCmd = &cobra.Command{
	Use:   "helm-chart-checker",
	Short: "See how out of date our cluster tools are",
	Long:  "Retrieve a list of HelmReleases from current context, and retrieve the latest version numbers from the corresponding Helm Repo",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		namespace := cmd.PersistentFlags().Lookup("namespace").Value
		getCharts(getHelmReleases(namespace.String()))
	},
}

func Execute() {
	rootCmd.Execute()
}
func getHelmReleases(ns string) *kubernetes.Resource {
	r := kubernetes.Resource{}
	r.Resource("helm.fluxcd.io", "v1", "helmreleases")
	r.GetHelmReleaseForNamespace(ns)
	return &r
}

func getCharts(hr *kubernetes.Resource) {
	charts := helm.Charts{}
	for _, item := range hr.Content {
		chart := helm.GetChartSpec(item.Object)
		if !funk.Contains(charts.Charts, *chart) {
			charts.Append(chart)
		}
	}
	for i, chart := range charts.Charts {
		if chart.IsHelmRepo() {
			chart.GetLatestChartVersion()
			charts.Charts[i] = chart
		}
	}
	charts.PrintChartData()
}
