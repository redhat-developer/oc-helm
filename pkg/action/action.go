package action

type Action interface {
	Run() error
	BuildHelmChartClient() error
}
