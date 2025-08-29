import { DataSourceJsonData } from '@grafana/data';
import { DataQuery } from '@grafana/schema';

export interface MyQuery extends DataQuery {
  resourceType: string;
  resourceParams: ResourceParams;
  limit?: number;
}

export const DEFAULT_QUERY: Partial<MyQuery> = {
  resourceType: 'job_runs',
  limit: 200
};

export type ResourceParams = JobRunQueryParams | PipelineQueryParams;

export interface JobRunQueryParams {
  jobId?: string;
  activeOnly?: boolean;
  completedOnly?: boolean;
  runType?: 'JOB_RUN' | 'WORKFLOW_RUN' | 'SUBMIT_RUN';
}

export interface PipelineQueryParams {
  filter?: string;
}

/**
 * These are options configured for each DataSource instance
 */
export interface MyDataSourceOptions extends DataSourceJsonData {
  workspace?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface MySecureJsonData {
  clientId?: string;
  clientSecret?: string;
}
