package plugin

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/databricks/databricks-grafana-datasource/pkg/models"
	"github.com/databricks/databricks-sdk-go"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
)

// Make sure Datasource implements required interfaces. This is important to do
// since otherwise we will only get a not implemented error response from plugin in
// runtime. In this example datasource instance implements backend.QueryDataHandler,
// backend.CheckHealthHandler interfaces. Plugin should not implement all these
// interfaces - only those which are required for a particular task.
var (
	_ backend.QueryDataHandler      = (*Datasource)(nil)
	_ backend.CheckHealthHandler    = (*Datasource)(nil)
	_ instancemgmt.InstanceDisposer = (*Datasource)(nil)
)

const (
	resourceTypeJobRuns         = "job_runs"
	resourceTypePipelines       = "pipelines"
	resourceTypePipelineUpdates = "pipeline_updates"
)

// NewDatasource creates a new datasource instance.
func NewDatasource(_ context.Context, settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	return &Datasource{
		settings: settings,
	}, nil
}

// Datasource is an example datasource which can respond to data queries, reports
// its health and has streaming skills.
type Datasource struct {
	settings backend.DataSourceInstanceSettings
}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created. As soon as datasource settings change detected by SDK old datasource instance will
// be disposed and a new one will be created using NewSampleDatasource factory function.
func (d *Datasource) Dispose() {
	// Clean up datasource instance resources.
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifier).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	// create response struct
	response := backend.NewQueryDataResponse()

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		res := d.query(ctx, req.PluginContext, q)

		// save the response in a hashmap
		// based on with RefID as identifier
		response.Responses[q.RefID] = res
	}

	return response, nil
}

type queryModel struct {
	RawQuery       string          `json:"queryText,omitempty"`
	ResourceType   string          `json:"resourceType"`
	ResourceParams json.RawMessage `json:"resourceParams,omitempty"`
	Limit          int             `json:"limit"`
}

func (d *Datasource) query(ctx context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	var qm queryModel

	err := json.Unmarshal(query.JSON, &qm)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("json unmarshal: %v", err.Error()))
	}

	// TODO: See if it makes sense to fetch the results via SQL statement from system.lakeflow and other tables instead (for perf reasons, in certain cases)
	switch qm.ResourceType {
	case resourceTypeJobRuns:
		return d.queryJobRuns(ctx, pCtx, query, qm)
	case resourceTypePipelines:
		return d.queryPipelines(ctx, pCtx, query, qm)
	default:
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("unknown resource kind: %s", qm.ResourceType))
	}
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (d *Datasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	res := &backend.CheckHealthResult{}
	config, err := models.LoadPluginSettings(*req.PluginContext.DataSourceInstanceSettings)

	if err != nil {
		res.Status = backend.HealthStatusError
		res.Message = "Unable to load settings"
		return res, nil
	}

	if config.Secrets.ClientId == "" || config.Secrets.ClientSecret == "" {
		res.Status = backend.HealthStatusError
		res.Message = "Authentication is missing"
		return res, nil
	}

	w, _ := d.getDatabricksClient(context.Background(), req.PluginContext)
	_, err = w.CurrentUser.Me(context.Background())

	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: fmt.Sprintf("Error: %v", err),
		}, nil
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Data source is working",
	}, nil
}

func (d *Datasource) getDatabricksClient(ctx context.Context, pCtx backend.PluginContext) (*databricks.WorkspaceClient, error) {
	config, err := models.LoadPluginSettings(d.settings)
	if err != nil {
		return nil, fmt.Errorf("load plugin settings: %v", err)
	}

	dbxConfig := databricks.Config{
		Host:         config.Workspace,
		ClientID:     config.Secrets.ClientId,
		ClientSecret: config.Secrets.ClientSecret,
	}

	return databricks.Must(databricks.NewWorkspaceClient(&dbxConfig)), nil
}
