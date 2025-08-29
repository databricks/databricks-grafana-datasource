import { InlineField, Input } from '@grafana/ui';
import React from 'react';
import { MyQuery, PipelineQueryParams } from 'types';

interface PipelinesEditorProps {
  resourceParams: PipelineQueryParams;
  onChange: (queryUpdate: Partial<MyQuery>) => void;
  onRunQuery: () => void;
}

export default function PipelinesEditor({ resourceParams, onChange, onRunQuery }: PipelinesEditorProps) {
  const onFilterChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    onChange({
      resourceParams: {
        ...resourceParams,
        filter: event.target.value,
      },
    });
  };

  return (
    <>
      <InlineField label="Filter" tooltip="Filter pipelines" grow>
        <Input
          placeholder="Filter pipelines by name, notebook, etc"
          value={resourceParams.filter || ''}
          onChange={onFilterChange}
          onBlur={onRunQuery}
        />
      </InlineField>
    </>
  );
}
