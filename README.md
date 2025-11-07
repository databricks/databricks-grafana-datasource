# Grafana Databricks Observability Plugin

A Grafana plugin for monitoring and visualizing Databricks Jobs, Pipelines, Costs and general observability, directly in your Grafana dashboards.

![screenshot](src/img/screenshot.jpg)

## Features

- Job Runs: Visualize your Databricks job executions, including run status, duration, and other metrics
- Pipelines: Query data about your Databricks Delta Live Tables pipelines
- Filtering Options: Flexible filtering by job ID, run type, and execution status

### Currently Available

- [x] Live and historical Job Status
- [x] Live and historical Pipeline Status

### Coming soon

- [ ] Cost exploration (classic and serverless)

## Getting Started (Quick)

## Installation

### Pre-requisites (building from source)

- Relatively recent version LTS of Node.js
- Relatively recent version of Go installed and configured in your `PATH`

### Build the plugin

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

### Signing the plugin

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

### Package the plugin

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

The following Dashboard are available out of the box:

### Job status monitoring

[Dashboard link](dashboards/jobs_demo.json)

This dashboard aggregates the job information from the workspace on the timeline view and gives an overview
of the currently running jobs (with their status) and historical information.

### Serverless cost monitoring

[Dashboard link](dasbhoards/serverless_costs.json)

This dashboard shows the serverless cost overview for a defined time period, with optional WoW or MoM views.
You can also filter by type of serverless workloads, and attribute costs to specific resource.
