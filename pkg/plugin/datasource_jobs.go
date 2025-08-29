package plugin

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/databricks/databricks-sdk-go/service/jobs"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

type jobRunParams struct {
	JobID         string `json:"jobId,omitempty"`
	ActiveOnly    bool   `json:"activeOnly,omitempty"`
	CompletedOnly bool   `json:"completedOnly,omitempty"`
	RunType       string `json:"runType,omitempty"`
}

func parseJobRunParams(_ backend.DataQuery, qm queryModel) (jobRunParams, error) {
	var params jobRunParams
	if qm.ResourceParams != nil {
		if err := json.Unmarshal(qm.ResourceParams, &params); err != nil {
			return params, fmt.Errorf("failed to unmarshal query params: %v", err)
		}
	}

	return params, nil
}

func buildListRunsRequest(params jobRunParams, query backend.DataQuery) (jobs.ListRunsRequest, error) {
	req := jobs.ListRunsRequest{
		Limit:         25, // 25 is the max limit for this API - per page
		ActiveOnly:    params.ActiveOnly,
		CompletedOnly: params.CompletedOnly,
	}

	// apply job id filter, if set
	if params.JobID != "" {
		jobId, err := strconv.ParseInt(params.JobID, 10, 64)
		if err != nil {
			return req, err
		}

		req.JobId = jobId
	}

	// apply time range filter, if set
	if !query.TimeRange.From.IsZero() && !query.TimeRange.To.IsZero() {
		req.StartTimeFrom = query.TimeRange.From.UnixMilli()
		req.StartTimeTo = query.TimeRange.To.UnixMilli()
	}

	// apply run type filter
	if params.RunType != "" {
		req.RunType = jobs.RunType(params.RunType)
	}

	return req, nil
}

func buildJobRunFrame(runs []jobs.BaseRun) *data.Frame {
	frame := data.NewFrame("Databricks Job Runs",
		data.NewField("Start Time", nil, []time.Time{}),
		data.NewField("End Time", nil, []time.Time{}),
		data.NewField("Job ID", nil, []string{}),
		data.NewField("Run ID", nil, []string{}),
		data.NewField("Run Name", nil, []string{}),
		data.NewField("Description", nil, []string{}),
		data.NewField("Attempt Number", nil, []int32{}),
		data.NewField("Status", nil, []string{}),
		data.NewField("Queue Duration (milliseconds)", nil, []int64{}),
		data.NewField("Run Duration (milliseconds)", nil, []int64{}),
		data.NewField("Run URL", nil, []string{}),
	)

	// sort results ascending by StartTime
	slices.SortFunc(runs, func(i, j jobs.BaseRun) int {
		return cmp.Compare(i.StartTime, j.StartTime)
	})

	for _, run := range runs {
		frame.AppendRow(
			time.UnixMilli(run.StartTime),
			time.UnixMilli(run.EndTime),
			fmt.Sprintf("%d", run.JobId),
			fmt.Sprintf("%d", run.RunId),
			run.RunName,
			run.Description,
			int32(run.AttemptNumber),
			string(run.Status.State),
			run.QueueDuration,
			run.RunDuration,
			run.RunPageUrl,
		)
	}

	return frame
}

func (d *Datasource) queryJobRuns(ctx context.Context, pCtx backend.PluginContext, query backend.DataQuery, qm queryModel) backend.DataResponse {
	params, err := parseJobRunParams(query, qm)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("failed to parse query params: %v", err))
	}

	w, err := d.getDatabricksClient(ctx, pCtx)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("failed to get databricks client: %v", err))
	}

	request, err := buildListRunsRequest(params, query)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("failed to build list runs request: %v", err))
	}

	jobsService := &workspaceClientWrapper{client: w}
	jobRuns, err := fetchJobRuns(ctx, jobsService, request, qm.Limit)

	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("failed to fetch job runs: %v", err))
	}

	var response backend.DataResponse
	frame := buildJobRunFrame(jobRuns)
	response.Frames = append(response.Frames, frame)
	return response
}

func fetchJobRuns(ctx context.Context, client DatabricksJobsService, request jobs.ListRunsRequest, maxItems int) ([]jobs.BaseRun, error) {
	iter := client.ListRuns(ctx, request)
	return fetchWithLimit(ctx, iter, maxItems)
}
