import React from 'react';
import { InlineField, InlineSwitch, Input, Select, Stack } from '@grafana/ui';
import { SelectableValue } from '@grafana/data';
import { JobRunQueryParams, MyQuery } from '../types';

export function JobRunsEditor({
  resourceParams,
  onChange,
  onRunQuery,
}: {
  resourceParams: JobRunQueryParams;
  onChange: (queryUpdate: Partial<MyQuery>) => void;
  onRunQuery: () => void;
}) {
  const onJobIdChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    onChange({
      resourceParams: {
        ...resourceParams,
        jobId: event.target.value,
      },
    });
  };

  const onActiveOnlyChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    onChange({
      resourceParams: {
        ...resourceParams,
        activeOnly: event.target.checked,
      },
    });
  };

  const onCompletedOnlyChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    onChange({
      resourceParams: {
        ...resourceParams,
        completedOnly: event.target.checked,
      },
    });
  };

  const onRunTypeChange = (value: SelectableValue<string>) => {
    onChange({
      resourceParams: {
        ...resourceParams,
        runType: value.value as 'JOB_RUN' | 'WORKFLOW_RUN' | 'SUBMIT_RUN',
      },
    });
  };

  return (
    <>
      <InlineField label="Job ID" tooltip="Filter runs by job ID" labelWidth={10}>
        <Input
          placeholder="Optional"
          value={resourceParams.jobId || ''}
          onChange={onJobIdChange}
          onBlur={onRunQuery}
          width={24}
        />
      </InlineField>

      <Stack direction="row" gap={1}>
        <InlineSwitch
          label="Active Only"
          showLabel={true}
          value={resourceParams.activeOnly || false}
          onChange={onActiveOnlyChange}
        />

        <InlineSwitch
          label="Completed Only"
          showLabel={true}
          value={resourceParams.completedOnly || false}
          onChange={onCompletedOnlyChange}
        />
      </Stack>

      <InlineField label="Run Type" tooltip="Filter by run type" labelWidth={14}>
        <Select
          options={[
            { label: 'Job Run', value: 'JOB_RUN' },
            { label: 'Workflow Run', value: 'WORKFLOW_RUN' },
            { label: 'Submit Run', value: 'SUBMIT_RUN' },
          ]}
          value={resourceParams.runType}
          onChange={onRunTypeChange}
          width={24}
          isClearable={true}
        />
      </InlineField>
    </>
  );
}
