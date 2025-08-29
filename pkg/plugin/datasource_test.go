package plugin

import (
	"context"
	"testing"
	"time"

	"github.com/databricks/databricks-sdk-go/service/jobs"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func TestQueryData(t *testing.T) {
	ds := Datasource{}

	resp, err := ds.QueryData(
		context.Background(),
		&backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "A"},
			},
		},
	)
	if err != nil {
		t.Error(err)
	}

	if len(resp.Responses) != 1 {
		t.Fatal("QueryData must return a response")
	}
}

func TestBuildListRunsRequest(t *testing.T) {
	t.Parallel()

	t.Run("should return valid request when no params are used", func(t *testing.T) {
		params := jobRunParams{}
		query := backend.DataQuery{}

		_, err := buildListRunsRequest(params, query)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("should pass valid job id", func(t *testing.T) {
		params := jobRunParams{
			JobID: "123",
		}

		query := backend.DataQuery{}

		req, err := buildListRunsRequest(params, query)
		if err != nil {
			t.Error(err)
		}

		if req.JobId != 123 {
			t.Errorf("expected job id to be 123, got %d", req.JobId)
		}
	})

	t.Run("should set active and completed", func(t *testing.T) {
		params := jobRunParams{
			ActiveOnly:    true,
			CompletedOnly: true,
		}

		query := backend.DataQuery{}

		req, err := buildListRunsRequest(params, query)
		if err != nil {
			t.Error(err)
		}

		if !req.ActiveOnly {
			t.Error("expected active only to be true")
		}

		if !req.CompletedOnly {
			t.Error("expected completed only to be true")
		}
	})

	t.Run("should return error when job id is invalid", func(t *testing.T) {
		params := jobRunParams{
			JobID: "invalid",
		}

		query := backend.DataQuery{}

		_, err := buildListRunsRequest(params, query)
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("should set time range", func(t *testing.T) {
		params := jobRunParams{}
		query := backend.DataQuery{
			TimeRange: backend.TimeRange{
				From: time.Now(),
				To:   time.Now(),
			},
		}

		req, err := buildListRunsRequest(params, query)
		if err != nil {
			t.Error(err)
		}

		if req.StartTimeFrom != query.TimeRange.From.UnixMilli() {
			t.Error("expected start time from to be set")
		}

		if req.StartTimeTo != query.TimeRange.To.UnixMilli() {
			t.Error("expected start time to to be set")
		}
	})

	t.Run("should set run type", func(t *testing.T) {
		params := jobRunParams{
			RunType: "WORKFLOW",
		}

		query := backend.DataQuery{}

		req, err := buildListRunsRequest(params, query)
		if err != nil {
			t.Error(err)
		}

		if req.RunType != "WORKFLOW" {
			t.Error("expected run type to be set")
		}
	})
}

func TestBuildJobRunFrame(t *testing.T) {
	t.Parallel()

	t.Run("should return valid frame", func(t *testing.T) {
		runs := []jobs.BaseRun{
			{
				StartTime:     time.Now().UnixMilli(),
				EndTime:       time.Now().UnixMilli(),
				JobId:         123,
				RunId:         456,
				RunName:       "run name",
				Description:   "description",
				AttemptNumber: 1,
				Status: &jobs.RunStatus{
					State: jobs.RunLifecycleStateV2StatePending,
				},
				QueueDuration: 1,
				RunDuration:   1,
				RunPageUrl:    "http://example.com",
			},
		}

		// TODO: assert values, for now we just trigger internal transformation
		frame := buildJobRunFrame(runs)
		if frame == nil {
			t.Error("expected frame to be returned")
		}

		_, err := frame.MarshalJSON()
		if err != nil {
			t.Error(err)
		}
	})
}

// TODO: Add tests for pipelines datasource
