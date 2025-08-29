package plugin

import (
	"context"

	"github.com/databricks/databricks-sdk-go"
	"github.com/databricks/databricks-sdk-go/listing"
	"github.com/databricks/databricks-sdk-go/service/jobs"
)

type DatabricksJobsService interface {
	ListRuns(ctx context.Context, request jobs.ListRunsRequest) listing.Iterator[jobs.BaseRun]
}

type workspaceClientWrapper struct {
	client *databricks.WorkspaceClient
}

func (w *workspaceClientWrapper) ListRuns(ctx context.Context, request jobs.ListRunsRequest) listing.Iterator[jobs.BaseRun] {
	return w.client.Jobs.ListRuns(ctx, request)
}
