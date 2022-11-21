# helm-chart-checker
A tool to check helm installations on the current kubernetes context with its corresponding repository
to see if the currently installed version is out of date

Currently only supports Flux v1 HelmReleases, with planned support for Argo Applications and Flux v2 HelmReleases

## Usage
```sh
❯ go run main.go --help
Retrieve a list of HelmReleases from current context, and retrieve the latest version numbers from the corresponding Helm Repo

Usage:
  helm-chart-checker [flags]

Flags:
  -h, --help               help for helm-chart-checker
  -n, --namespace string   Namespace to look (default "default")
  -o, --only-out-of-date   Show only out of date charts
```
