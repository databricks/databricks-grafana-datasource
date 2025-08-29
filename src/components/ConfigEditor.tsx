import React, { ChangeEvent } from 'react';
import { InlineField, Input, SecretInput } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { MyDataSourceOptions, MySecureJsonData } from '../types';

interface Props extends DataSourcePluginOptionsEditorProps<MyDataSourceOptions, MySecureJsonData> {}

export function ConfigEditor(props: Props) {
  const { onOptionsChange, options } = props;
  const { jsonData, secureJsonFields, secureJsonData } = options;

  const onWorkspaceChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...jsonData,
        workspace: event.target.value,
      },
    });
  };

  const onClientIdChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      secureJsonData: {
        ...options.secureJsonData,
        clientId: event.target.value,
      },
    });
  };

  const onClientSecretChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      secureJsonData: {
        ...options.secureJsonData,
        clientSecret: event.target.value,
      },
    });
  };

  const onResetClientId = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        clientId: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        clientId: '',
      },
    });
  };

  const onResetClientSecret = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        clientSecret: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        clientSecret: '',
      },
    });
  };

  return (
    <>
      <InlineField label="Workspace URL" labelWidth={20} interactive tooltip={'Databricks Workspace URL'}>
        <Input
          id="config-editor-workspace"
          onChange={onWorkspaceChange}
          value={jsonData.workspace}
          placeholder="Enter the workspace URL, e.g. https://adb-xxx.0.azuredatabricks.net"
          width={40}
          autoComplete="off"
        />
      </InlineField>
      <InlineField label="Client ID" labelWidth={20} interactive tooltip={'Service principal Client ID'}>
        <SecretInput
          required
          id="config-editor-client-id"
          isConfigured={secureJsonFields.clientId}
          value={secureJsonData?.clientId}
          placeholder="Enter your Databricks SP Client Id"
          width={40}
          onReset={onResetClientId}
          onChange={onClientIdChange}
          autoComplete="new-password"
        />
      </InlineField>

      <InlineField label="Client Secret" labelWidth={20} interactive tooltip={'Service principal Client Secret'}>
        <SecretInput
          required
          id="config-editor-client-secret"
          isConfigured={secureJsonFields.clientSecret}
          value={secureJsonData?.clientSecret}
          placeholder="Enter your Databricks SP Client Secret"
          width={40}
          onReset={onResetClientSecret}
          onChange={onClientSecretChange}
          autoComplete="new-password"
        />
      </InlineField>
    </>
  );
}
