# Grafana Databricks Observability Plugin

A Grafana plugin for monitoring and visualizing Databricks Jobs, DLTs and Pipelines directly in your Grafana dashboards.

![screenshot](src/img/screenshot.jpg)

## Features

- Job Runs: Visualize your Databricks job executions, including run status, duration, and other metrics
- Pipelines: Query data about your Databricks Delta Live Tables pipelines
- Filtering Options: Flexible filtering by job ID, run type, and execution status

Please see the [CHANGELOG](./CHANGELOG.md) for the latest updates and changes.

## Installation

Depending on your configuration you can either download the pre-packaged zip file or the source code and build it yourself (see instructions below).

### Using the Release version

Easiest way to use the plugin is to use the [latest release](https://github.com/databricks/databricks-grafana-datasource/releases/latest).

If your Grafana instance is configured to only allow signed plugins to run, you will need to sign the plugin yourself (see instructions below).

#### Signing the release plugin

1. Download the latest release zip file from the [releases page](https://github.com/databricks/databricks-grafana-datasource/releases/latest).
2. Unzip the file into `dist` directory:

   ```bash
   unzip databricks-grafana-datasource-x.y.z.zip -d dist
   ```

3. Create Access Policy Token in your Grafana Cloud instance (Security > Access Policies). This token will be used to sign the plugin.

4. Export the token to your environment:

   ```bash
   export GRAFANA_ACCESS_POLICY_TOKEN=glc...
   ```

5. Sign the plugin, replacing the root URL with your Grafana instance URL:

   ```bash
   npx @grafana/sign-plugin@latest --rootUrls http://localhost:3000/
   ```

6. Re-package the plugin into a zip file (if needed):

   ```bash
   mv dist databricks-grafana-datasource
   zip -r databricks-grafana-datasource.zip databricks-grafana-datasource
   ```

You can now install the signed plugin in your Grafana instance. To quickly validate the plugin, you can spin up a local Grafana instance using Docker (where `$PWD/plugins` is the directory containing the signed plugin `databricks-grafana-datasource`):

```bash
docker run --rm -p 3000:3000 \
  -v "$PWD/plugins:/var/lib/grafana/plugins" \
  -e GF_AUTH_DISABLE_LOGIN_FORM=true \
  -e GF_AUTH_ANONYMOUS_ENABLED=true \
  -e GF_AUTH_ANONYMOUS_ORG_ROLE=Admin \
  grafana/grafana:latest
```

If everything is set up correctly, you should be able to see the "Databricks Community" plugin in the Grafana UI under `Configuration > Data Sources > Add data source` and `msg="Plugin registered" pluginId=databricks-grafana-datasource` message in the Grafana server logs.

### Building from Source

#### Pre-requisites

- Relatively recent version LTS of Node.js
- Relatively recent version of Go installed and configured in your `PATH`

#### Build the plugin

This plugin needs to be built from the source and signed so you could use it in your hosted Grafana instance. Follow the steps below to build and install the plugin:

1. Clone the repository:

   ```bash
   git clone git@github.com:databricks/databricks-grafana-datasource.git
   ```

2. Navigate to the plugin directory:

   ```bash
   cd databricks-grafana-datasource
   ```

3. Build the backend (here we're running custom mage target to skip building for Linux ARM):

   ```bash
   mage buildAllNoLinuxArm
   ```

   Your `dist` directory should now contain native binaries for various platforms.
4. Build the frontend:

   ```bash
   npm install
   npm run build
   ```

   Your `dist` directory should now contain the built frontend files.

#### Signing the plugin

The plugin needs to be signed before it can be used in a Grafana Cloud instance. The signing process ensures that the plugin is verified and trusted by Grafana.

1. Create Access Policy Token in your Grafana Cloud instance (Security > Access Policies). This token will be used to sign the plugin.
2. Export the token to your environment:

   ```bash
   export GRAFANA_ACCESS_POLICY_TOKEN=glc...
   ```

3. Sign the plugin, replacing the root URL with your Grafana instance URL:

   ```bash
   npx @grafana/sign-plugin@latest --rootUrls http://localhost:3000/
   ```

#### Package the plugin

The plugin needs to be packaged into a zip file before it can be installed in Grafana.

1. Create a zip file of the plugin:

   ```bash
   mv dist databricks-grafana-datasource
   zip -r databricks-grafana-datasource.zip databricks-grafana-datasource
   ```

You can now install the plugin in your Grafana instance.

## Pre-requisites (running in Grafana)

- Grafana `10` or later
- Databricks workspace (any cloud provider)
- Databricks Service Principal with the following permissions:
  - Permissions to list either all or specific jobs / pipelines, **or**:
  - Admin permissions in the workspace

## Configuration

1. In Grafana, navigate to `Configuration > Data Sources > Add data source`
2. Search for "Databricks Community" and select this plugin
3. Configure the following settings:
  - `Workspace URL`: Your Databricks workspace URL (e.g., https://adb-xxx.0.azuredatabricks.net)
  - `Client ID`: Service Principal Client ID
  - `Client Secret`: Service Principal Client Secret
4. Click "Save & Test" to verify the connection

## Supported Data Sources

The query editor provides different options based on the selected resource type:

### Job Runs

- `Job ID`: Optional filter to show runs for a specific job
- `Active Only`: Toggle to show only currently running jobs
- `Completed Only`: Toggle to show only completed jobs
- `Run Type`: Filter by job run type (JOB_RUN, WORKFLOW_RUN, or SUBMIT_RUN)
- `Max Results`: Maximum number of results to return (default: 200)

### Pipelines

- `Filter`: Text filter for pipeline queries
- `Max Results`: Maximum number of results to return (default: 200)

## Example Dashboards

Please refer to the [dashboards](./dashboards) directory for example dashboards that demonstrate the capabilities of this plugin.
