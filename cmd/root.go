package cmd

import "github.com/spf13/cobra"
func init() {
	rootCmd.PersistentFlags().StringP("namespace","n","default","Namespace to look")
}

var rootCmd = &cobra.Command{
	Use:   "helm-chart-checker",
	Short: "See how out of date our cluster tools are",
	Long:"Retrieve a list of HelmReleases from current context, and retrieve the latest version numbers from the corresponding Helm Repo"
	Run: func(cmd *cobra.Command, args []string) {
	  // Do Stuff Here
	},
  }
  
  func Execute() {
	return rootCmd.Execute()
  }
  