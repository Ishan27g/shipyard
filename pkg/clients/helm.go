package clients

import (
	"fmt"
	"sync"

	"github.com/hashicorp/go-hclog"
	"golang.org/x/xerrors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/kube"
)

var helmLock sync.Mutex

func init() {
	// create a global lock as it seems map write in Helm is not thread safe
	helmLock = sync.Mutex{}
}

// Helm defines an interface for a client which can manage Helm charts
type Helm interface {
	Create(kubeConfig, name, namespace string, createNamespace bool, chartPath, valuesPath string, valuesString map[string]string) error
	Destroy(kubeConfig, name, namespace string) error
}

type HelmImpl struct {
	log hclog.Logger
}

func NewHelm(l hclog.Logger) Helm {
	return &HelmImpl{l}
}

// Create a new install of the chart
func (h *HelmImpl) Create(kubeConfig, name, namespace string, createNamespace bool, chartPath, valuesPath string, valuesString map[string]string) error {
	// set the kubeclient for Helm
	s := kube.GetConfig(kubeConfig, "default", namespace)
	cfg := &action.Configuration{}
	err := cfg.Init(s, namespace, "", func(format string, v ...interface{}) {
		h.log.Debug("Helm debug message", "message", fmt.Sprintf(format, v...))
	})

	if err != nil {
		return xerrors.Errorf("unalbe to iniailize Helm: %w", err)
	}

	client := action.NewInstall(cfg)
	client.ReleaseName = name
	client.Namespace = namespace
	client.CreateNamespace = createNamespace

	settings := cli.EnvSettings{}
	p := getter.All(&settings)
	vo := values.Options{}
	vo.StringValues = []string{}

	// add the string values to the collection
	for k, v := range valuesString {
		vo.StringValues = append(vo.StringValues, fmt.Sprintf("%s=%s", k, v))
	}

	// if we have an overriden values file set it
	if valuesPath != "" {
		vo.ValueFiles = []string{valuesPath}
	}

	h.log.Debug("Creating chart from config", "ref", name, "path", chartPath)
	cp, err := client.ChartPathOptions.LocateChart(chartPath, &settings)
	if err != nil {
		return xerrors.Errorf("Error locating chart: %w", err)
	}

	h.log.Debug("Loading chart", "ref", name, "path", cp)
	chartRequested, err := loader.Load(cp)
	if err != nil {
		return xerrors.Errorf("Error loading chart: %w", err)
	}

	vals, err := vo.MergeValues(p)
	if err != nil {
		return xerrors.Errorf("Error merging Helm values: %w", err)
	}

	h.log.Debug("Validate chart", "ref", name)
	err = chartRequested.Validate()
	if err != nil {
		return xerrors.Errorf("Error validating chart: %w", err)
	}

	h.log.Debug("Run chart", "ref", name)
	_, err = client.Run(chartRequested, vals)
	if err != nil {
		return xerrors.Errorf("Error running chart: %w", err)
	}

	return nil
}

// Destroy removes an installed Helm chart from the system
func (h *HelmImpl) Destroy(kubeConfig, name, namespace string) error {
	s := kube.GetConfig(kubeConfig, "default", namespace)
	cfg := &action.Configuration{}
	err := cfg.Init(s, namespace, "", func(format string, v ...interface{}) {
		h.log.Debug("Helm debug message", "message", fmt.Sprintf(format, v...))
	})

	//settings := cli.EnvSettings{}
	//p := getter.All(&settings)
	//vo := values.Options{}
	client := action.NewUninstall(cfg)
	_, err = client.Run(name)
	if err != nil {
		h.log.Debug("Unable to remove chart, exit silently", "err", err)
		return err
	}

	return nil
}
