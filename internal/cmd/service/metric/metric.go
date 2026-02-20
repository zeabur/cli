package metric

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
	"github.com/zeabur/cli/pkg/fill"
	"github.com/zeabur/cli/pkg/model"
)

type Options struct {
	id            string
	name          string
	environmentID string
	metricType    string
	hour          uint
}

func NewCmdMetric(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "metric [CPU|MEMORY|NETWORK]",
		Short: "Show metric of a service",
		Long:  `Show metric of a service`,
		Args:  cobra.ExactArgs(1),
		ValidArgs: []string{
			string(model.MetricTypeCPU),
			string(model.MetricTypeMemory),
			string(model.MetricTypeNetwork),
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.metricType = args[0]
			return runMetric(f, opts)
		},
	}

	util.AddServiceParam(cmd, &opts.id, &opts.name)
	util.AddEnvOfServiceParam(cmd, &opts.environmentID)
	cmd.Flags().StringVarP(&opts.metricType, "metric-type", "t", "", "Metric type, one of CPU, MEMORY, NETWORK")
	cmd.Flags().UintVarP(&opts.hour, "hour", "H", 2, "Metric history in hour")

	return cmd
}

func runMetric(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		return runMetricInteractive(f, opts)
	} else {
		return runMetricNonInteractive(f, opts)
	}
}

func runMetricInteractive(f *cmdutil.Factory, opts *Options) error {
	zctx := f.Config.GetContext()
	if _, err := f.ParamFiller.ServiceByNameWithEnvironment(fill.ServiceByNameWithEnvironmentOptions{
		ProjectCtx:    zctx,
		ServiceID:     &opts.id,
		ServiceName:   &opts.name,
		EnvironmentID: &opts.environmentID,
		CreateNew:     false,
	}); err != nil {
		return err
	}

	return runMetricNonInteractive(f, opts)
}

func runMetricNonInteractive(f *cmdutil.Factory, opts *Options) error {
	if opts.id == "" && opts.name != "" {
		service, err := util.GetServiceByName(f.Config, f.ApiClient, opts.name)
		if err != nil {
			return err
		}
		opts.id = service.ID
	}

	if opts.id == "" {
		return fmt.Errorf("--id or --name is required")
	}

	// Resolve environment and project from the service
	service, err := f.ApiClient.GetService(context.Background(), opts.id, "", "", "")
	if err != nil {
		return fmt.Errorf("get service failed: %w", err)
	}

	projectID := ""
	if service.Project != nil {
		projectID = service.Project.ID
	}

	if opts.environmentID == "" {
		envID, err := util.ResolveEnvironmentID(f.ApiClient, projectID)
		if err != nil {
			return err
		}
		opts.environmentID = envID
	}

	if opts.metricType == "" {
		return fmt.Errorf("metric type is required")
	}

	upperCaseMetricType := strings.ToUpper(opts.metricType)
	mt := model.MetricType(opts.metricType)

	startTime := time.Now().Add(-time.Duration(opts.hour) * time.Hour)
	endTime := time.Now()

	metrics, err := f.ApiClient.ServiceMetric(context.Background(), opts.id, projectID, opts.environmentID, upperCaseMetricType, startTime, endTime)
	if err != nil {
		return fmt.Errorf("get service metric failed: %w", err)
	}

	if len(metrics.Metrics) == 0 {
		f.Log.Infof("no metric history found")
		return nil
	}

	sum, avg, max, min := 0.0, 0.0, 0.0, 0.0

	for _, metric := range metrics.Metrics {
		sum += metric.Value
		if metric.Value > max {
			max = metric.Value
		}
		if metric.Value < min {
			min = metric.Value
		}

		if f.Debug {
			f.Log.Debugf("Metric: %s, value: %s, timestamp: %s",
				opts.metricType, mt.WithMeasureUnit(metric.Value), metric.Timestamp.Format(time.RFC3339))
		}
	}

	avg = sum / float64(len(metrics.Metrics))

	f.Log.Infof("Metric type: %s, start time: %s, end time: %s\n",
		opts.metricType, startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))

	header := []string{"Sum", "Avg", "Max", "Min"}
	data := [][]string{{mt.WithMeasureUnit(sum), mt.WithMeasureUnit(avg), mt.WithMeasureUnit(max), mt.WithMeasureUnit(min)}}

	f.Printer.Table(header, data)

	return nil
}
