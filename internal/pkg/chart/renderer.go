package chart

import (
	"bytes"
	"context"
	"fmt"

	"github.com/adelowo/gulter"
	"github.com/adelowo/snapshot-chromedp/render"
	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ChartDataPoint struct {
	Label string      `json:"label"`
	Value interface{} `json:"value"`
}

type ChartData struct {
	ChartType  string
	DataPoints []ChartDataPoint
}

type EChartsRenderer struct {
	storage     gulter.Storage
	cfg         config.Config
	db          *bun.DB
	integration malak.IntegrationRepository
}

func NewEChartsRenderer(storage gulter.Storage, cfg config.Config, db *bun.DB,
	integration malak.IntegrationRepository) *EChartsRenderer {
	return &EChartsRenderer{
		storage:     storage,
		cfg:         cfg,
		db:          db,
		integration: integration,
	}
}

func (r *EChartsRenderer) RenderChart(workspaceID uuid.UUID, chartID string) (string, error) {
	chart, err := r.integration.GetChart(context.Background(), malak.FetchChartOptions{
		WorkspaceID: workspaceID,
		Reference:   malak.Reference(chartID),
	})
	if err != nil {
		return "", fmt.Errorf("failed to fetch chart: %w", err)
	}

	// Check for unsupported chart type before fetching data points
	switch chart.ChartType {
	case malak.IntegrationChartTypeBar, malak.IntegrationChartTypePie:
		// These are supported, continue processing
	default:
		return "", fmt.Errorf("unsupported chart type: %s", chart.ChartType)
	}

	dataPoints, err := r.integration.GetDataPoints(context.Background(), chart)
	if err != nil {
		return "", fmt.Errorf("failed to fetch chart data points: %w", err)
	}

	chartData := make([]ChartDataPoint, len(dataPoints))
	for i, point := range dataPoints {
		var value interface{} = point.PointValue

		if chart.DataPointType == malak.IntegrationDataPointTypeCurrency {
			value = float64(point.PointValue) / 100.0
		}

		chartData[i] = ChartDataPoint{
			Label: point.PointName,
			Value: value,
		}
	}

	var chartHTML []byte
	switch chart.ChartType {
	case "bar":
		chartHTML, err = r.renderBarChart(chartData)
	case "pie":
		chartHTML, err = r.renderPieChart(chartData)
	}

	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%s-%s.png", chart.ChartType, uuid.New().String())

	var b = bytes.NewBuffer(nil)

	// Wrap the chart HTML in a container with explicit dimensions
	fullHTML := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <style>
        body { margin: 0; }
        #container { width: 600px; height: 400px; }
    </style>
</head>
<body>
    <div id="container">%s</div>
</body>
</html>`, string(chartHTML))

	if err := render.MakeChartSnapshotWriter([]byte(fullHTML), b); err != nil {
		return "", err
	}

	file, err := r.storage.Upload(context.Background(), b, &gulter.UploadFileOptions{
		FileName: filename,
	})

	uploadedURL, err := r.storage.Path(context.Background(), gulter.PathOptions{
		Key: file.Key,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload chart: %w", err)
	}

	return uploadedURL, nil
}

func (r *EChartsRenderer) renderBarChart(data []ChartDataPoint) ([]byte, error) {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width:      "600px",
			Height:     "400px",
			Theme:      "white",
			AssetsHost: "",
		}),
		charts.WithTitleOpts(opts.Title{
			Title: "Chart",
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: opts.Bool(true)}),
		charts.WithLegendOpts(opts.Legend{Show: opts.Bool(true)}),
		charts.WithYAxisOpts(opts.YAxis{
			Show:      opts.Bool(true),
			Type:      "value",
			AxisLabel: &opts.AxisLabel{Show: opts.Bool(true)},
		}),
		charts.WithXAxisOpts(opts.XAxis{
			AxisLabel: &opts.AxisLabel{Show: opts.Bool(true)},
		}),
		charts.WithToolboxOpts(opts.Toolbox{
			Show: opts.Bool(true),
		}),
	)

	xAxis := make([]string, len(data))
	values := make([]opts.BarData, len(data))

	for i, point := range data {
		xAxis[i] = point.Label
		values[i] = opts.BarData{Value: point.Value}
	}

	bar.SetXAxis(xAxis).AddSeries("Values", values).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:     opts.Bool(true),
				Position: "top",
			}),
			charts.WithAnimationOpts(opts.Animation{
				Animation: opts.Bool(false),
			}),
		)

	return bar.RenderContent(), nil
}

func (r *EChartsRenderer) renderPieChart(data []ChartDataPoint) ([]byte, error) {
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width:      "600px",
			Height:     "400px",
			Theme:      "white",
			AssetsHost: "",
		}),
		charts.WithTitleOpts(opts.Title{
			Title: "Chart",
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: opts.Bool(true)}),
		charts.WithLegendOpts(opts.Legend{Show: opts.Bool(true)}),
		charts.WithToolboxOpts(opts.Toolbox{
			Show: opts.Bool(true),
		}),
	)

	items := make([]opts.PieData, len(data))
	for i, point := range data {
		items[i] = opts.PieData{
			Name:  point.Label,
			Value: point.Value,
		}
	}

	pie.AddSeries("Values", items).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:     opts.Bool(true),
				Position: "outside",
			}),
			charts.WithAnimationOpts(opts.Animation{
				Animation: opts.Bool(false),
			}),
		)

	return pie.RenderContent(), nil
}
