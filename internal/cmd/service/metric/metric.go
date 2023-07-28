package metric

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/model"
	"time"
)

type Options struct {
	id            string
	environmentID string
	metricType    string
	projectID     string
	hour          uint
}

func NewCmdMetric(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{
		projectID: f.Config.GetContext().GetProject().GetID(),
	}

	cmd := &cobra.Command{
		Use:   "metric <metric-type>",
		Short: "Show metric of a service",
		Long:  `Show metric of a service`,
		Args:  cobra.ExactArgs(1),
		ValidArgs: []string{
			string(model.MetricTypeCPU),
			string(model.MetricTypeMemory),
			//string(model.MetricTypeNetwork), // not supported yet
			//string(model.MetricTypeDisk),	// not supported yet
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.metricType = args[0]
			return runMetric(f, opts)
		},
	}

	ctx := f.Config.GetContext()

	//todo: support identify service by name

	cmd.Flags().StringVar(&opts.id, "id", ctx.GetService().GetID(), "Service ID")
	cmd.Flags().StringVar(&opts.environmentID, "environment-id", ctx.GetEnvironment().GetID(), "Environment ID")
	cmd.Flags().StringVarP(&opts.metricType, "metric-type", "t", "", "Metric type, one of CPU, MEMORY, NETWORK, DISK")
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
	if _, err := f.ParamFiller.ServiceWithEnvironment(&opts.projectID, &opts.id, &opts.environmentID); err != nil {
		return err
	}

	return runMetricNonInteractive(f, opts)
}

func runMetricNonInteractive(f *cmdutil.Factory, opts *Options) error {
	if opts.id == "" {
		return fmt.Errorf("--id is required")
	}

	if opts.environmentID == "" {
		return fmt.Errorf("--environment-id is required")
	}

	if opts.metricType == "" {
		return fmt.Errorf("metric type is required")
	}

	mt := model.MetricType(opts.metricType)

	startTime := time.Now().Add(-time.Duration(opts.hour) * time.Hour)
	endTime := time.Now()

	metrics, err := f.ApiClient.ServiceMetric(context.Background(), opts.id, opts.environmentID, opts.metricType, startTime, endTime)
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

	// todo: support chart? (add a new method in printer.Printer)

	return nil
}
