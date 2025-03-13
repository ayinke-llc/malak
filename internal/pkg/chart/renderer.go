package chart

import (
	"bytes"
	"context"
	"fmt"

	"github.com/adelowo/gulter"
	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ChartDataPoint struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
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
	// Fetch chart data from database
	chartData, err := r.fetchChartData(workspaceID, chartID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch chart data: %w", err)
	}

	var chartHTML string
	switch chartData.ChartType {
	case "bar":
		chartHTML, err = r.renderBarChart(chartData.DataPoints)
	case "pie":
		chartHTML, err = r.renderPieChart(chartData.DataPoints)
	default:
		return "", fmt.Errorf("unsupported chart type: %s", chartData.ChartType)
	}

	if err != nil {
		return "", err
	}

	// Generate a unique key for the chart
	filename := fmt.Sprintf("%s-%s.png", chartData.ChartType, uuid.New().String())

	// Upload the chart HTML to storage
	file, err := r.storage.Upload(context.Background(), bytes.NewReader([]byte(chartHTML)), &gulter.UploadFileOptions{
		FileName: filename,
		Bucket:   r.cfg.Uploader.S3.Bucket,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload chart: %w", err)
	}

	uploadedURL := fmt.Sprintf("%s/%s/%s",
		r.cfg.Uploader.S3.Endpoint,
		file.FolderDestination,
		file.Key)

	return uploadedURL, nil
}

func (r *EChartsRenderer) fetchChartData(workspaceID uuid.UUID, chartID string) (*ChartData, error) {
	chart, err := r.integration.GetChart(context.Background(), malak.FetchChartOptions{
		WorkspaceID: workspaceID,
		Reference:   malak.Reference(chartID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch chart: %w", err)
	}

	dataPoints, err := r.integration.GetDataPoints(context.Background(), chart)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch chart data points: %w", err)
	}

	chartData := make([]ChartDataPoint, len(dataPoints))
	for i, point := range dataPoints {
		chartData[i] = ChartDataPoint{
			Label: point.PointName,
			Value: float64(point.PointValue),
		}
	}

	return &ChartData{
		ChartType:  string(chart.ChartType),
		DataPoints: chartData,
	}, nil
}

func (r *EChartsRenderer) renderBarChart(data []ChartDataPoint) (string, error) {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "600px",
			Height: "400px",
		}),
		charts.WithTitleOpts(opts.Title{
			Title: "Chart",
		}),
	)

	xAxis := make([]string, len(data))
	values := make([]opts.BarData, len(data))

	for i, point := range data {
		xAxis[i] = point.Label
		values[i] = opts.BarData{Value: point.Value}
	}

	bar.SetXAxis(xAxis).AddSeries("Values", values)

	buffer := new(bytes.Buffer)
	if err := bar.Render(buffer); err != nil {
		return "", fmt.Errorf("failed to render bar chart: %w", err)
	}

	return buffer.String(), nil
}

func (r *EChartsRenderer) renderPieChart(data []ChartDataPoint) (string, error) {
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "600px",
			Height: "400px",
		}),
		charts.WithTitleOpts(opts.Title{
			Title: "Chart",
		}),
	)

	items := make([]opts.PieData, len(data))
	for i, point := range data {
		items[i] = opts.PieData{
			Name:  point.Label,
			Value: point.Value,
		}
	}

	pie.AddSeries("Values", items)

	buffer := new(bytes.Buffer)
	if err := pie.Render(buffer); err != nil {
		return "", fmt.Errorf("failed to render pie chart: %w", err)
	}

	return buffer.String(), nil
}
