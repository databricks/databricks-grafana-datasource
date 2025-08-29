import React from 'react';
import { InlineField, Input, Select, Stack } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from '../datasource';
import { JobRunQueryParams, MyDataSourceOptions, MyQuery, PipelineQueryParams } from '../types';
import { JobRunsEditor } from './JobRunsEditor';
import PipelinesEditor from './PipelinesEditor';

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

export function QueryEditor({ query, onChange, onRunQuery }: Props) {
  const resourceParams = query.resourceParams || {};
  const onResourceTypeChange = (value: SelectableValue<string>) => {
    onChange({ ...query, resourceType: value.value! });
    onRunQuery();
  };

  const onLimitChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const limitValue = parseInt(event.target.value, 10);
    if (!isNaN(limitValue)) {
      onChange({ ...query, limit: limitValue });
    }
  };

  const handleQueryChange = (updates: Partial<MyQuery>) => {
    onChange({ ...query, ...updates });
  };

  const resourceEditor = () => {
    switch (query.resourceType) {
      case 'job_runs':
        return (
          <JobRunsEditor
            resourceParams={resourceParams as JobRunQueryParams}
            onChange={handleQueryChange}
            onRunQuery={onRunQuery}
          />
        );

      case 'pipelines':
        return (
          <PipelinesEditor
            resourceParams={resourceParams as PipelineQueryParams}
            onChange={handleQueryChange}
            onRunQuery={onRunQuery}
          />
        );

      default:
        return null;
    }
  };

  return (
    <Stack gap={0}>
      <InlineField label="Resource Type">
        <Select
          options={[
            { label: 'Job Runs', value: 'job_runs' },
            { label: 'Pipelines', value: 'pipelines' },
          ]}
          value={query.resourceType}
          onChange={onResourceTypeChange}
        />
      </InlineField>

      {resourceEditor()}

      <InlineField label="Max Results" labelWidth={14} tooltip="Maximum number of results to fetch">
        <Input
          type="number"
          placeholder="200"
          value={query.limit || 200}
          onChange={onLimitChange}
          onBlur={onRunQuery}
          width={12}
        />
      </InlineField>
    </Stack>
  );
}
