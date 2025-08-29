package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/databricks/databricks-sdk-go/service/pipelines"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

type pipelineParams struct {
	Filter string `json:"filter,omitempty"`
}

type pipelineUpdatesParams struct {
	PipelineId string `json:"pipelineId"`
}

func parsePipelineParams(_ backend.DataQuery, qm queryModel) (pipelineParams, error) {
	var params pipelineParams
	if qm.ResourceParams != nil {
		if err := json.Unmarshal(qm.ResourceParams, &params); err != nil {
			return params, fmt.Errorf("failed to unmarshal query params: %v", err)
		}
	}

	return params, nil
}

func parsePipelineUpdateParams(_ backend.DataQuery, qm queryModel) (pipelineUpdatesParams, error) {
	var params pipelineUpdatesParams
	if qm.ResourceParams != nil {
		if err := json.Unmarshal(qm.ResourceParams, &params); err != nil {
			return params, fmt.Errorf("failed to unmarshal query params: %v", err)
		}
	}

	return params, nil
}

func buildPipelineRequest(params pipelineParams, query backend.DataQuery) (pipelines.ListPipelinesRequest, error) {
	req := pipelines.ListPipelinesRequest{
		Filter: params.Filter,
	}

	return req, nil
}

func buildPipelineUpdateRequest(params pipelineUpdatesParams, query backend.DataQuery) (pipelines.ListUpdatesRequest, error) {
	req := pipelines.ListUpdatesRequest{
		PipelineId: params.PipelineId,
		MaxResults: 100,
	}

	return req, nil
}

func buildPipelinesRunFrame(pipelines []pipelines.PipelineStateInfo) *data.Frame {
	frame := data.NewFrame("pipelines")
	frame.Fields = append(frame.Fields,
		data.NewField("Pipeline Id", nil, []string{}),
		data.NewField("Pipeline Name", nil, []string{}),
		data.NewField("State", nil, []string{}),
	)

	for _, pipeline := range pipelines {
		frame.AppendRow(
			pipeline.PipelineId,
			pipeline.Name,
			string(pipeline.State),
		)
	}

	return frame
}

func buildPipelineUpdatesFrame(updates []pipelines.UpdateInfo) *data.Frame {
	frame := data.NewFrame("pipeline updates")
	frame.Fields = append(frame.Fields,
		data.NewField("Creation Time", nil, []time.Time{}),
		data.NewField("Update Id", nil, []string{}),
		data.NewField("Pipeline Id", nil, []string{}),
		data.NewField("Cause", nil, []string{}),
		data.NewField("State", nil, []string{}),
	)

	for _, update := range updates {
		frame.AppendRow(
			time.UnixMilli(update.CreationTime),
			update.UpdateId,
			update.PipelineId,
			string(update.Cause),
			string(update.State),
		)
	}

	return frame
}

func (d *Datasource) queryPipelines(ctx context.Context, pCtx backend.PluginContext, query backend.DataQuery, qm queryModel) backend.DataResponse {
	params, err := parsePipelineParams(query, qm)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("failed to parse query params: %v", err))
	}

	w, err := d.getDatabricksClient(ctx, pCtx)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("failed to get databricks client: %v", err))
	}

	request, err := buildPipelineRequest(params, query)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("failed to build request: %v", err))
	}

	pipelinesIter := w.Pipelines.ListPipelines(context.Background(), request)
	pipelines, err := fetchWithLimit(ctx, pipelinesIter, qm.Limit)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("failed to list pipelines: %v", err))
	}

	frame := buildPipelinesRunFrame(pipelines)
	return backend.DataResponse{
		Frames: []*data.Frame{frame},
	}
}

func (d *Datasource) queryPipelineUpdates(ctx context.Context, pCtx backend.PluginContext, query backend.DataQuery, qm queryModel) backend.DataResponse {
	params, err := parsePipelineUpdateParams(query, qm)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("failed to parse query params: %v", err))
	}

	w, err := d.getDatabricksClient(ctx, pCtx)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("failed to get databricks client: %v", err))
	}

	request, err := buildPipelineUpdateRequest(params, query)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("failed to build request: %v", err))
	}

	updates, err := w.Pipelines.ListUpdates(context.Background(), request)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("failed to list updates: %v", err))
	}

	frame := buildPipelineUpdatesFrame(updates.Updates)
	return backend.DataResponse{
		Frames: []*data.Frame{frame},
	}
}
